package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	db "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// GetFanTokens returns all existing fantokens
func (k Keeper) GetFanTokens(ctx sdk.Context, owner sdk.AccAddress) (fantokens []types.FanToken) {
	store := ctx.KVStore(k.storeKey)

	var it db.Iterator
	if owner == nil {
		it = storetypes.KVStorePrefixIterator(store, types.PrefixFanTokenForDenom)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var fantoken types.FanToken
			k.cdc.MustUnmarshal(it.Value(), &fantoken)

			fantokens = append(fantokens, fantoken)
		}
		return
	}

	it = storetypes.KVStorePrefixIterator(store, types.KeyFanTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var denom gogotypes.StringValue
		k.cdc.MustUnmarshal(it.Value(), &denom)

		fantoken, err := k.getFanTokenByDenom(ctx, denom.Value)
		if err != nil {
			continue
		}
		fantokens = append(fantokens, fantoken)
	}
	return
}

// GetFanToken returns the fantoken of the specified denom
func (k Keeper) GetFanToken(ctx sdk.Context, denom string) (*types.FanToken, error) {
	// query fantoken by denom
	if fantoken, err := k.getFanTokenByDenom(ctx, denom); err == nil {
		return &fantoken, nil
	}

	return nil, sdkerrors.ErrAppConfig.Wrapf(types.ErrFanTokenNotExists.Error(), "denom %s does not exist", denom)
}

// AddFanToken saves a new token
func (k Keeper) AddFanToken(ctx sdk.Context, token *types.FanToken) error {
	if k.HasFanToken(ctx, token.GetDenom()) {
		return sdkerrors.ErrAppConfig.Wrapf(types.ErrDenomAlreadyExists.Error(), "denom already exists: %s", token.GetDenom())
	}

	// set token
	k.setFanToken(ctx, token)

	if len(token.MetaData.Authority) != 0 {
		// set token to be prefixed with metadata authority
		k.setWithMetadataAuthority(ctx, token.GetAuthority(), token.GetDenom())
	}

	return nil
}

// getFanTokenSupply queries the fantoken supply from the total supply
func (k Keeper) getFanTokenSupply(ctx sdk.Context, denom string) math.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}
