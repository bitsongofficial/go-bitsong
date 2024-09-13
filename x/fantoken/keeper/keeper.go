package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
	// accountKeeper types.AccountKeeper
	bankKeeper   types.BankKeeper
	distrKeeper  types.DistrKeeper
	paramSpace   types.ParamSubspace
	blockedAddrs map[string]bool
}

func NewKeeper(
	cdc codec.Codec,
	key storetypes.StoreKey,
	paramSpace types.ParamSubspace,
	ak types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	blockedAddrs map[string]bool,
) Keeper {
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the " + types.ModuleName + " module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace,
		bankKeeper:   bankKeeper,
		distrKeeper:  distrKeeper,
		blockedAddrs: blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("go-bitsong/%s", types.ModuleName))
}

// Issue issues a new fantoken
func (k Keeper) Issue(ctx sdk.Context, name, symbol, uri string, maxSupply sdk.Int, minter, authority sdk.AccAddress) (denom string, err error) {
	if k.blockedAddrs[authority.String()] {
		return denom, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", authority.String())
	}

	// at the moment is disabled, will be enabled once some test will be done
	if k.blockedAddrs[minter.String()] {
		return denom, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", minter.String())
	}

	// check minter
	if minter.Empty() {
		return denom, sdkerrors.Wrapf(types.ErrInvalidMinter, "the address %s is not a valid minter address", minter)
	}

	// handle issue fee
	if err := k.deductIssueFee(ctx, minter); err != nil {
		return denom, err
	}

	fantoken := types.NewFanToken(name, symbol, uri, maxSupply, minter, authority, ctx.BlockHeight())
	if err := fantoken.Validate(); err != nil {
		return denom, err
	}

	found := k.HasFanToken(ctx, fantoken.Denom)
	if found {
		return denom, types.ErrDenomAlreadyExists
	}

	if err := k.AddFanToken(ctx, fantoken); err != nil {
		return denom, err
	}

	return fantoken.GetDenom(), nil
}

// Mint mints the specified amount of fantoken to the specified recipient
func (k Keeper) Mint(ctx sdk.Context, minter, recipient sdk.AccAddress, coin sdk.Coin) error {
	if recipient.Empty() {
		return sdkerrors.Wrapf(types.ErrInvalidRecipient, "the address %s is not a valid recipient", recipient.String())
	}

	if minter.Empty() {
		return sdkerrors.Wrapf(types.ErrInvalidMinter, "the address %s is not a valid minter address", minter.String())
	}

	if k.blockedAddrs[minter.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", minter.String())
	}

	if k.blockedAddrs[recipient.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", recipient.String())
	}

	if err := types.ValidateAmount(coin.Amount); err != nil {
		return err
	}

	fantoken, err := k.getFanTokenByDenom(ctx, coin.Denom)
	if err != nil {
		return err
	}

	if minter.String() != fantoken.Minter {
		return sdkerrors.Wrapf(types.ErrInvalidMinter, "the address %s is not the minter of the fantoken %s", minter.String(), coin.Denom)
	}

	// handle Mint fee
	if err := k.deductMintFee(ctx, minter); err != nil {
		return err
	}

	supply := k.getFanTokenSupply(ctx, fantoken.GetDenom())
	mintableAmt := fantoken.MaxSupply.Sub(supply)

	if coin.Amount.GT(mintableAmt) {
		return sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"the amount exceeds the mintable fantoken amount; expected [0, %d], got %d",
			mintableAmt.Int64(), coin.Amount.Int64(),
		)
	}

	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	// send coins to the recipient account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, sdk.NewCoins(coin))
}

// Burn burns the specified amount of fantoken
func (k Keeper) Burn(ctx sdk.Context, coin sdk.Coin, owner sdk.AccAddress) error {
	if k.blockedAddrs[owner.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", owner.String())
	}

	if owner.Empty() {
		return types.ErrInvalidOwner
	}

	// handle Burn fee
	if err := k.deductBurnFee(ctx, owner); err != nil {
		return err
	}

	found := k.HasFanToken(ctx, coin.Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrFanTokenNotExists, "fantoken not found: %s", coin.Denom)
	}

	// Burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
}

// SetAuthority transfers the authority of the specified fantoken to a new one
func (k Keeper) SetAuthority(ctx sdk.Context, denom string, oldAuthority, newAuthority sdk.AccAddress) error {
	if k.blockedAddrs[oldAuthority.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", oldAuthority.String())
	}

	if k.blockedAddrs[newAuthority.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", newAuthority.String())
	}

	if oldAuthority.Empty() {
		return types.ErrInvalidAuthority
	}

	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if oldAuthority.String() != fantoken.MetaData.Authority {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the address %s is not the authority of the fantoken %s", oldAuthority, denom)
	}

	if fantoken.GetAuthority().String() == "" {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the metadata are immutable")
	}

	fantoken.MetaData.Authority = newAuthority.String()

	if err := fantoken.Validate(); err != nil {
		return err
	}

	// update fantoken
	k.setFanToken(ctx, &fantoken)

	// reset all indices
	k.resetStoreKeyForQueryToken(ctx, fantoken.GetDenom(), oldAuthority, newAuthority)

	return nil
}

// SetMinter transfers the minter of the specified fantoken to a new one
func (k Keeper) SetMinter(ctx sdk.Context, denom string, oldMinter, newMinter sdk.AccAddress) error {
	if k.blockedAddrs[oldMinter.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", oldMinter.String())
	}

	if k.blockedAddrs[newMinter.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", newMinter.String())
	}

	if oldMinter.Empty() {
		return types.ErrInvalidMinter
	}

	// get the fantoken
	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if oldMinter.String() != fantoken.Minter {
		return sdkerrors.Wrapf(types.ErrInvalidMinter, "the address %s is not the minter of the fantoken %s", oldMinter, denom)
	}

	if fantoken.Minter == "" {
		return sdkerrors.Wrapf(types.ErrInvalidMinter, "the minting is disabled")
	}

	fantoken.Minter = newMinter.String()

	if newMinter.String() == "" {
		// at this point we can set the official supply
		supply := k.getFanTokenSupply(ctx, fantoken.GetDenom())
		fantoken.MaxSupply = supply
	}

	if err := fantoken.Validate(); err != nil {
		return err
	}

	// update fantoken
	k.setFanToken(ctx, &fantoken)

	return nil
}

func (k Keeper) SetUri(ctx sdk.Context, denom, newUri string, authority sdk.AccAddress) error {
	// get the fantoken
	fantoken, err := k.getFanTokenByDenom(ctx, denom)
	if err != nil {
		return err
	}

	if authority.Empty() {
		return types.ErrInvalidAuthority
	}

	if authority.String() != fantoken.MetaData.Authority {
		return sdkerrors.Wrapf(types.ErrInvalidAuthority, "the address %s is not the authority of the fantoken %s", authority, denom)
	}

	if err := types.ValidateUri(newUri); err != nil {
		return err
	}

	fantoken.MetaData.URI = newUri

	if err := fantoken.Validate(); err != nil {
		return err
	}

	// update fantoken
	k.setFanToken(ctx, &fantoken)

	return nil
}
