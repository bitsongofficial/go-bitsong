package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/album/types"
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

// GetDeposit gets the deposit of a specific depositor on a specific album
func (keeper Keeper) GetDeposit(ctx sdk.Context, albumID uint64, depositorAddr sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(albumID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) SetDeposit(ctx sdk.Context, albumID uint64, depositorAddr sdk.AccAddress, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(albumID, depositorAddr), bz)
}

func (keeper Keeper) AddDeposit(ctx sdk.Context, albumID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (sdk.Error, bool) {
	// Checks to see if the album exists
	album, ok := keeper.GetAlbum(ctx, albumID)
	if !ok {
		return types.ErrUnknownAlbum(keeper.codespace, fmt.Sprintf("unknown albumID: %d", albumID)), false
	}

	// Check if album is still depositable
	if !(album.Status == types.StatusNil || album.Status == types.StatusDepositPeriod) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("albumID %d already deposited", albumID)), false
	}

	// If status is Nil enable deposit period
	album.Status = types.StatusDepositPeriod

	// Set deposit end time
	blockTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod
	album.DepositEndTime = blockTime.Add(depositPeriod)

	// update the album module's account coins pool
	err := keeper.Sk.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount)
	if err != nil {
		return err, false
	}

	// Increment total deposit
	album.TotalDeposit = album.TotalDeposit.Add(depositAmount)

	// Check if deposit has provided sufficient total funds to transition the album into the verified state
	verified := false
	if album.Status == types.StatusDepositPeriod && album.TotalDeposit.IsAllGTE(keeper.GetDepositParams(ctx).MinDeposit) {
		album.Status = types.StatusVerified
		album.VerifiedTime = blockTime
		verified = true
	}

	// Update album
	keeper.SetAlbum(ctx, album)

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, albumID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
	} else {
		deposit = types.NewDeposit(albumID, depositorAddr, depositAmount)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDepositAlbum,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyAlbumID, fmt.Sprintf("%d", albumID)),
		),
	)

	keeper.SetDeposit(ctx, albumID, depositorAddr, deposit)
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

func (keeper Keeper) IterateDeposits(ctx sdk.Context, albumID uint64, cb func(deposit types.Deposit) (stop bool)) {
	iterator := keeper.GetDepositsIterator(ctx, albumID)

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

// GetDeposits returns all the deposits from a specific album
func (keeper Keeper) GetDeposits(ctx sdk.Context, albumID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, albumID, func(deposit types.Deposit) bool {
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

// RefundDeposits refunds and deletes all the deposits on a specific album
func (keeper Keeper) RefundDeposits(ctx sdk.Context, albumID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, albumID, func(deposit types.Deposit) bool {
		err := keeper.Sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, deposit.Depositor, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(albumID, deposit.Depositor))
		return false
	})
}

// DeleteDeposits deletes all the deposits on a specific album without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, albumID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, albumID, func(deposit types.Deposit) bool {
		err := keeper.Sk.BurnCoins(ctx, types.ModuleName, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(albumID, deposit.Depositor))
		return false
	})
}
