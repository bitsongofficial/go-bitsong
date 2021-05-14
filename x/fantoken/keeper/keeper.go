package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.Marshaler
	bankKeeper       types.BankKeeper
	paramSpace       paramstypes.Subspace
	blockedAddrs     map[string]bool
	feeCollectorName string
}

func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	bankKeeper types.BankKeeper,
	blockedAddrs map[string]bool,
	feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		bankKeeper:       bankKeeper,
		feeCollectorName: feeCollectorName,
		blockedAddrs:     blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}

// IssueToken issues a new token
func (k Keeper) IssueFanToken(
	ctx sdk.Context,
	symbol string,
	name string,
	maxSupply sdk.Int,
	mintable bool,
	description string,
	owner sdk.AccAddress,
) error {
	denom := fmt.Sprintf("%s%s", "u", symbol)
	denomMetaData := banktypes.Metadata{
		Description: description,
		Base:        denom,
		Display:     symbol,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: denom, Exponent: 0},
			{Denom: symbol, Exponent: types.FanTokenDecimal},
		},
	}
	token := types.NewFanToken(name, maxSupply, mintable, owner, denomMetaData)

	if err := k.AddFanToken(ctx, token); err != nil {
		return err
	}

	return nil
}

// EditToken edits the specified token
func (k Keeper) EditFanToken(
	ctx sdk.Context,
	symbol string,
	mintable bool,
	owner sdk.AccAddress,
) error {
	// get the destination token
	token, err := k.getFanTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if owner.String() != token.Owner {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", owner, symbol)
	}

	token.Mintable = mintable

	k.setFanToken(ctx, token)

	return nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferFanTokenOwner(
	ctx sdk.Context,
	symbol string,
	srcOwner sdk.AccAddress,
	dstOwner sdk.AccAddress,
) error {
	token, err := k.getFanTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if srcOwner.String() != token.Owner {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", srcOwner, symbol)
	}

	token.Owner = dstOwner.String()

	// update token
	k.setFanToken(ctx, token)

	// reset all indices
	k.resetStoreKeyForQueryToken(ctx, token.GetSymbol(), srcOwner, dstOwner)

	return nil
}

// MintToken mints the specified amount of token to the specified recipient
// NOTE: empty owner means that the external caller is responsible to manage the token authority
func (k Keeper) MintFanToken(
	ctx sdk.Context,
	recipient sdk.AccAddress,
	symbol string,
	amount sdk.Int,
	owner sdk.AccAddress,
) error {
	token, err := k.getFanTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if owner.String() != token.Owner {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "the address %s is not the owner of the token %s", owner, symbol)
	}

	if !token.Mintable {
		return sdkerrors.Wrapf(types.ErrNotMintable, "%s", symbol)
	}

	supply := k.getFanTokenSupply(ctx, token.GetDenom())
	precision := sdk.NewIntWithDecimal(1, types.FanTokenDecimal)
	mintableAmt := token.MaxSupply.Mul(precision).Sub(supply)
	mintableMainAmt := mintableAmt.Quo(precision)

	if amount.GT(mintableMainAmt) {
		return sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"the amount exceeds the mintable token amount; expected (0, %d], got %d",
			mintableMainAmt, amount,
		)
	}

	mintCoin := sdk.NewCoin(token.GetDenom(), amount.Mul(precision))
	mintCoins := sdk.NewCoins(mintCoin)

	// mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	if recipient.Empty() {
		recipient = owner
	}

	// sent coins to the recipient account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, mintCoins)
}

// BurnToken burns the specified amount of token
func (k Keeper) BurnFanToken(
	ctx sdk.Context,
	symbol string,
	amount sdk.Int,
	owner sdk.AccAddress,
) error {
	token, err := k.getFanTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	precision := sdk.NewIntWithDecimal(1, types.FanTokenDecimal)
	burnCoin := sdk.NewCoin(token.GetDenom(), amount.Mul(precision))
	burnCoins := sdk.NewCoins(burnCoin)

	// burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, burnCoins); err != nil {
		return err
	}

	k.AddBurnCoin(ctx, burnCoin)

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
}
