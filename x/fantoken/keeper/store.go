package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"
)

// hasFanToken asserts a fantoken exists
func (k Keeper) hasFanToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenom(denom))
}

func (k Keeper) SetBurnCoin(ctx sdk.Context, total sdk.Coin) {
	bz := k.cdc.MustMarshal(&total)
	key := types.KeyBurnFanTokenAmt(total.Denom)

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
}

// getBurnedCoins returns the total amount of the burned fantoken
func (k Keeper) getBurnedCoins(ctx sdk.Context, denom string) sdk.Coin {
	key := types.KeyBurnFanTokenAmt(denom)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return sdk.Coin{Denom: denom, Amount: sdk.ZeroInt()}
	}

	var coin sdk.Coin
	k.cdc.MustUnmarshal(bz, &coin)

	return coin
}

// GetAllBurnedCoins returns the total amount of all the burned fantokens
func (k Keeper) GetAllBurnedCoins(ctx sdk.Context) []sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	var coins []sdk.Coin
	it := sdk.KVStorePrefixIterator(store, types.PefixBurnFanTokenAmt)
	for ; it.Valid(); it.Next() {
		var coin sdk.Coin
		k.cdc.MustUnmarshal(it.Value(), &coin)
		coins = append(coins, coin)
	}

	return coins
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, denom string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: denom})
	store.Set(types.KeyFanTokens(owner, denom), bz)
}

func (k Keeper) setFanToken(ctx sdk.Context, token *types.FanToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(token)
	store.Set(types.KeyDenom(token.GetDenom()), bz)
}

func (k Keeper) getFanTokenByDenom(ctx sdk.Context, denom string) (fantoken *types.FanToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenom(denom))
	if bz == nil {
		return fantoken, sdkerrors.Wrap(types.ErrFanTokenNotExists, fmt.Sprintf("fantoken denom %s does not exist", denom))
	}

	k.cdc.MustUnmarshal(bz, fantoken)
	return fantoken, nil
}

// reset all indices by the new owner for fantoken query
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, denom string, srcOwner, dstOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyFanTokens(srcOwner, denom))

	// add the new key
	k.setWithOwner(ctx, dstOwner, denom)
}
