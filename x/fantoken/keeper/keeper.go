package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	distr "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.Codec
	bankKeeper       types.BankKeeper
	distrKeeper      distr.Keeper
	paramSpace       paramstypes.Subspace
	blockedAddrs     map[string]bool
	feeCollectorName string
}

func NewKeeper(
	cdc codec.Codec,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	bankKeeper types.BankKeeper,
	distrKeeper distr.Keeper,
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
		distrKeeper:      distrKeeper,
		feeCollectorName: feeCollectorName,
		blockedAddrs:     blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}

// Issue issues a new fantoken
func (k Keeper) Issue(ctx sdk.Context, name, symbol, uri string, maxSupply sdk.Int, authority sdk.AccAddress) (denom string, err error) {
	// handle issue fee
	if err := k.deductIssueFee(ctx, authority); err != nil {
		return denom, err
	}

	fantoken := types.NewFanToken(name, symbol, uri, maxSupply, authority, ctx.BlockHeight())
	if err := k.AddFanToken(ctx, fantoken); err != nil {
		return denom, err
	}

	return fantoken.GetDenom(), nil
}

// DisableMint disable the mint of a specific fantoken
func (k Keeper) DisableMint(ctx sdk.Context, denom string, authority sdk.AccAddress) error {
	// get the fantoken
	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if authority.String() != fantoken.Authority {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the address %s is not the authority of the fantoken %s", authority, denom)
	}

	if !fantoken.Mintable {
		return sdkerrors.Wrapf(types.ErrNotMintable, "the fantoken %s is not mintable", denom)
	}

	fantoken.Mintable = false
	fantoken.Authority = ""

	supply := k.getFanTokenSupply(ctx, fantoken.GetDenom())
	fantoken.MaxSupply = supply

	k.setFanToken(ctx, &fantoken)

	return nil
}

// TransferAuthority transfers the owner of the specified fantoken to a new one
func (k Keeper) TransferAuthority(ctx sdk.Context, denom string, srcAuthority, dstAuthority sdk.AccAddress) error {
	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if srcAuthority.String() != fantoken.Authority {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the address %s is not the authority of the fantoken %s", srcAuthority, denom)
	}

	// handle transfer fee
	if err := k.deductTransferFee(ctx, srcAuthority); err != nil {
		return err
	}

	fantoken.Authority = dstAuthority.String()

	// update fantoken
	k.setFanToken(ctx, &fantoken)

	// reset all indices
	k.resetStoreKeyForQueryToken(ctx, fantoken.GetDenom(), srcAuthority, dstAuthority)

	return nil
}

// Mint mints the specified amount of fantoken to the specified recipient
func (k Keeper) Mint(ctx sdk.Context, recipient sdk.AccAddress, denom string, amount sdk.Int, authority sdk.AccAddress) error {
	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if authority.String() != fantoken.Authority {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the address %s is not the authority of the fantoken %s", authority, denom)
	}

	// handle mint fee
	if err := k.deductMintFee(ctx, authority); err != nil {
		return err
	}

	if !fantoken.Mintable {
		return sdkerrors.Wrapf(types.ErrNotMintable, "%s", denom)
	}

	supply := k.getFanTokenSupply(ctx, fantoken.GetDenom())
	burnedCoins := k.getBurnedCoins(ctx, fantoken.GetDenom())
	mintableAmt := fantoken.MaxSupply.Sub(supply).Sub(burnedCoins.Amount)

	if amount.GT(mintableAmt) {
		return sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"the amount exceeds the mintable fantoken amount; expected (0, %d], got %d",
			mintableAmt, amount,
		)
	}

	mintCoin := sdk.NewCoin(fantoken.GetDenom(), amount)
	mintCoins := sdk.NewCoins(mintCoin)

	// mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoins); err != nil {
		return err
	}

	if recipient.Empty() {
		recipient = authority
	}

	// sent coins to the recipient account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, mintCoins)
}

// Burn burns the specified amount of fantoken
func (k Keeper) Burn(ctx sdk.Context, denom string, amount sdk.Int, owner sdk.AccAddress) error {
	found := k.HasFanToken(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrFanTokenNotExists, "fantoken not found: %s", denom)
	}

	// handle burn fee
	if err := k.deductBurnFee(ctx, owner); err != nil {
		return err
	}

	burnCoin := sdk.NewCoin(denom, amount)
	burnCoins := sdk.NewCoins(burnCoin)

	// burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, burnCoins); err != nil {
		return err
	}

	k.AddBurnCoin(ctx, burnCoin)

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
}
