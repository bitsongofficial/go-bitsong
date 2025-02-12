package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
)

// HasFanToken asserts a fantoken exists
func (k Keeper) HasFanToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenom(denom))
}

func (k Keeper) setWithMetadataAuthority(ctx sdk.Context, owner sdk.AccAddress, denom string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: denom})
	store.Set(types.KeyFanTokens(owner, denom), bz)
}

func (k Keeper) setFanToken(ctx sdk.Context, token *types.FanToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(token)
	store.Set(types.KeyDenom(token.GetDenom()), bz)
}

func (k Keeper) getFanTokenByDenom(ctx sdk.Context, denom string) (fantoken types.FanToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenom(denom))
	if bz == nil {
		return fantoken, errors.Wrap(types.ErrFanTokenNotExists, fmt.Sprintf("fantoken denom %s does not exist", denom))
	}

	k.cdc.MustUnmarshal(bz, &fantoken)
	return fantoken, nil
}

// reset all indices by the new owner for fantoken query
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, denom string, srcOwner, dstOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyFanTokens(srcOwner, denom))

	// add the new key
	k.setWithMetadataAuthority(ctx, dstOwner, denom)
}
