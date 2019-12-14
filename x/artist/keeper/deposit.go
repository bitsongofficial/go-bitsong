package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDepositParams returns the current DepositParams from the global param store
func (keeper Keeper) GetDepositParams(ctx sdk.Context) types.DepositParams {
	var depositParams types.DepositParams
	keeper.paramSpace.Get(ctx, types.ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

func (keeper Keeper) SetDepositParams(ctx sdk.Context, depositParams types.DepositParams) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeyDepositParams, &depositParams)
}

// GetDeposit gets the deposit of a specific depositor on a specific artist
func (keeper Keeper) GetDeposit(ctx sdk.Context, artistID uint64, depositorAddr sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(artistID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) SetDeposit(ctx sdk.Context, artistID uint64, depositorAddr sdk.AccAddress, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(artistID, depositorAddr), bz)
}

func (keeper Keeper) AddDeposit(ctx sdk.Context, artistID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (sdk.Error, bool) {
	// Checks to see if the artist exists
	artist, ok := keeper.GetArtist(ctx, artistID)
	if !ok {
		return types.ErrUnknownArtist(keeper.codespace, fmt.Sprintf("unknown artistID: %d", artistID)), false
	}

	// Check if artist is still depositable
	if !(artist.Status == types.StatusNil || artist.Status == types.StatusDepositPeriod) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("artistID %d already deposited", artistID)), false
	}

	// If status is Nil enable deposit period
	artist.Status = types.StatusDepositPeriod

	// Set deposit end time
	blockTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod
	artist.DepositEndTime = blockTime.Add(depositPeriod)

	// update the artist module's account coins pool
	err := keeper.Sk.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount)
	if err != nil {
		return err, false
	}

	// Increment total deposit
	artist.TotalDeposit = artist.TotalDeposit.Add(depositAmount)

	// Check if deposit has provided sufficient total funds to transition the artist into the verified state
	verified := false
	if artist.Status == types.StatusDepositPeriod && artist.TotalDeposit.IsAllGTE(keeper.GetDepositParams(ctx).MinDeposit) {
		artist.Status = types.StatusVerified
		artist.VerifiedTime = blockTime
		verified = true
	}

	// Update the artist
	keeper.SetArtist(ctx, artist)

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, artistID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
	} else {
		deposit = types.NewDeposit(artistID, depositorAddr, depositAmount)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositArtist,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyArtistID, fmt.Sprintf("%d", artistID)),
		),
	)

	keeper.SetDeposit(ctx, artistID, depositorAddr, deposit)
	return nil, verified
}

func (keeper Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

func (keeper Keeper) IterateDeposits(ctx sdk.Context, artistID uint64, cb func(deposit types.Deposit) (stop bool)) {
	iterator := keeper.GetDepositsIterator(ctx, artistID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// GetAllDeposits returns all the deposits from the store
func (keeper Keeper) GetAllDeposits(ctx sdk.Context) (deposits types.Deposits) {
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDeposits returns all the deposits from a specific artist
func (keeper Keeper) GetDeposits(ctx sdk.Context, artistID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, artistID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDepositsIterator gets all the deposits on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetDepositsIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))
}

// RefundDeposits refunds and deletes all the deposits on a specific artist
func (keeper Keeper) RefundDeposits(ctx sdk.Context, artistID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, artistID, func(deposit types.Deposit) bool {
		err := keeper.Sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(artistID, deposit.Depositor))
		return false
	})
}

// DeleteDeposits deletes all the deposits on a specific artist without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, artistID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, artistID, func(deposit types.Deposit) bool {
		err := keeper.Sk.BurnCoins(ctx, types.ModuleName, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(artistID, deposit.Depositor))
		return false
	})
}
