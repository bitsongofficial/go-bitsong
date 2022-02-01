package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// GetTokens returns all existing tokens
func (k Keeper) GetFanTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []types.FanTokenI) {
	store := ctx.KVStore(k.storeKey)

	var it sdk.Iterator
	if owner == nil {
		it = sdk.KVStorePrefixIterator(store, types.PrefixFanTokenForDenom)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var token types.FanToken
			k.cdc.MustUnmarshal(it.Value(), &token)

			tokens = append(tokens, &token)
		}
		return
	}

	it = sdk.KVStorePrefixIterator(store, types.KeyFanTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var denom gogotypes.StringValue
		k.cdc.MustUnmarshal(it.Value(), &denom)

		token, err := k.getFanTokenByDenom(ctx, denom.Value)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return
}

// GetToken returns the token of the specified denom
func (k Keeper) GetFanToken(ctx sdk.Context, denom string) (types.FanTokenI, error) {
	// query token by denom
	if token, err := k.getFanTokenByDenom(ctx, denom); err == nil {
		return &token, nil
	}

	return nil, sdkerrors.Wrapf(types.ErrFanTokenNotExists, "fantoken: %s does not exist", denom)
}

// AddToken saves a new token
func (k Keeper) AddFanToken(ctx sdk.Context, token types.FanToken) error {
	if k.HasFanToken(ctx, token.GetDenom()) {
		return sdkerrors.Wrapf(types.ErrDenomAlreadyExists, "denom already exists: %s", token.GetDenom())
	}

	// set token
	k.setFanToken(ctx, token)

	if len(token.Owner) != 0 {
		// set token to be prefixed with owner
		k.setWithOwner(ctx, token.GetOwner(), token.GetDenom())
	}

	k.bankKeeper.SetDenomMetaData(ctx, token.MetaData)

	return nil
}

// HasToken asserts a token exists
func (k Keeper) HasFanToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyDenom(denom))
}

// GetOwner returns the owner of the specified token
func (k Keeper) GetOwner(ctx sdk.Context, denom string) (sdk.AccAddress, error) {
	token, err := k.GetFanToken(ctx, denom)
	if err != nil {
		return nil, err
	}

	return token.GetOwner(), nil
}

// AddBurnCoin saves the total amount of the burned tokens
func (k Keeper) AddBurnCoin(ctx sdk.Context, coin sdk.Coin) {
	var total = coin
	if hasCoin, err := k.GetBurnCoin(ctx, coin.Denom); err == nil {
		total = total.Add(hasCoin)
	}

	bz := k.cdc.MustMarshal(&total)
	key := types.KeyBurnFanTokenAmt(coin.Denom)

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
}

// GetBurnCoin returns the total amount of the burned tokens
func (k Keeper) GetBurnCoin(ctx sdk.Context, denom string) (sdk.Coin, error) {
	key := types.KeyBurnFanTokenAmt(denom)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrNotFoundTokenAmt, "not found denom: %s", denom)
	}

	var coin sdk.Coin
	k.cdc.MustUnmarshal(bz, &coin)

	return coin, nil
}

// GetAllBurnCoin returns the total amount of all the burned tokens
func (k Keeper) GetAllBurnCoin(ctx sdk.Context) []sdk.Coin {
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

// GetParamSet returns token params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet sets token params to the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, denom string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: denom})

	store.Set(types.KeyFanTokens(owner, denom), bz)
}

func (k Keeper) setFanToken(ctx sdk.Context, token types.FanToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&token)

	store.Set(types.KeyDenom(token.GetDenom()), bz)
}

func (k Keeper) getFanTokenByDenom(ctx sdk.Context, denom string) (token types.FanToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyDenom(denom))
	if bz == nil {
		return token, sdkerrors.Wrap(types.ErrFanTokenNotExists, fmt.Sprintf("token denom %s does not exist", denom))
	}

	k.cdc.MustUnmarshal(bz, &token)
	return token, nil
}

// reset all indices by the new owner for token query
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, denom string, srcOwner, dstOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.KeyFanTokens(srcOwner, denom))

	// add the new key
	k.setWithOwner(ctx, dstOwner, denom)
}

// getTokenSupply queries the token supply from the total supply
func (k Keeper) getFanTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}
