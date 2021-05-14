package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	tokentypes "github.com/bitsongofficial/bitsong/x/fantoken/types"
)

// GetTokens returns all existing tokens
func (k Keeper) GetFanTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []tokentypes.FanTokenI) {
	store := ctx.KVStore(k.storeKey)

	var it sdk.Iterator
	if owner == nil {
		it = sdk.KVStorePrefixIterator(store, tokentypes.PrefixFanTokenForSymbol)
		defer it.Close()

		for ; it.Valid(); it.Next() {
			var token tokentypes.FanToken
			k.cdc.MustUnmarshalBinaryBare(it.Value(), &token)

			tokens = append(tokens, &token)
		}
		return
	}

	it = sdk.KVStorePrefixIterator(store, tokentypes.KeyFanTokens(owner, ""))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var symbol gogotypes.StringValue
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &symbol)

		token, err := k.getFanTokenBySymbol(ctx, symbol.Value)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}
	return
}

// GetToken returns the token of the specified symbol or min uint
func (k Keeper) GetFanToken(ctx sdk.Context, denom string) (tokentypes.FanTokenI, error) {
	// query token by symbol
	if token, err := k.getFanTokenBySymbol(ctx, denom); err == nil {
		return &token, nil
	}

	// query token by min unit
	if token, err := k.getTokenByDenom(ctx, denom); err == nil {
		return &token, nil
	}

	return nil, sdkerrors.Wrapf(tokentypes.ErrTokenNotExists, "token: %s does not exist", denom)
}

// AddToken saves a new token
func (k Keeper) AddFanToken(ctx sdk.Context, token tokentypes.FanToken) error {
	if k.HasFanToken(ctx, token.GetSymbol()) {
		return sdkerrors.Wrapf(tokentypes.ErrSymbolAlreadyExists, "symbol already exists: %s", token.GetSymbol())
	}

	if k.HasFanToken(ctx, token.GetDenom()) {
		return sdkerrors.Wrapf(tokentypes.ErrDenomAlreadyExists, "denom already exists: %s", token.GetDenom())
	}

	// set token
	k.setFanToken(ctx, token)

	// set token to be prefixed with denom
	k.setWithDenom(ctx, token.GetDenom(), token.GetSymbol())

	if len(token.Owner) != 0 {
		// set token to be prefixed with owner
		k.setWithOwner(ctx, token.GetOwner(), token.GetSymbol())
	}

	k.bankKeeper.SetDenomMetaData(ctx, token.MetaData)

	return nil
}

// HasSymbol asserts a token exists by symbol
func (k Keeper) HasSymbol(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(tokentypes.KeySymbol(symbol))
}

// HasToken asserts a token exists
func (k Keeper) HasFanToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	if k.HasSymbol(ctx, denom) {
		return true
	}

	return store.Has(tokentypes.KeyDenom(denom))
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

	bz := k.cdc.MustMarshalBinaryBare(&total)
	key := tokentypes.KeyBurnFanTokenAmt(coin.Denom)

	store := ctx.KVStore(k.storeKey)
	store.Set(key, bz)
}

// GetBurnCoin returns the total amount of the burned tokens
func (k Keeper) GetBurnCoin(ctx sdk.Context, denom string) (sdk.Coin, error) {
	key := tokentypes.KeyBurnFanTokenAmt(denom)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	if len(bz) == 0 {
		return sdk.Coin{}, sdkerrors.Wrapf(tokentypes.ErrNotFoundTokenAmt, "not found symbol: %s", denom)
	}

	var coin sdk.Coin
	k.cdc.MustUnmarshalBinaryBare(bz, &coin)

	return coin, nil
}

// GetAllBurnCoin returns the total amount of all the burned tokens
func (k Keeper) GetAllBurnCoin(ctx sdk.Context) []sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	var coins []sdk.Coin
	it := sdk.KVStorePrefixIterator(store, tokentypes.PefixBurnFanTokenAmt)
	for ; it.Valid(); it.Next() {
		var coin sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &coin)
		coins = append(coins, coin)
	}

	return coins
}

// GetParamSet returns token params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) tokentypes.Params {
	var p tokentypes.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParamSet sets token params to the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params tokentypes.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) setWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(tokentypes.KeyFanTokens(owner, symbol), bz)
}

func (k Keeper) setWithDenom(ctx sdk.Context, denom, symbol string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})

	store.Set(tokentypes.KeyDenom(denom), bz)
}

func (k Keeper) setFanToken(ctx sdk.Context, token tokentypes.FanToken) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&token)

	store.Set(tokentypes.KeySymbol(token.GetSymbol()), bz)
}

func (k Keeper) getFanTokenBySymbol(ctx sdk.Context, symbol string) (token tokentypes.FanToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(tokentypes.KeySymbol(symbol))
	if bz == nil {
		return token, sdkerrors.Wrap(tokentypes.ErrTokenNotExists, fmt.Sprintf("token symbol %s does not exist", symbol))
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token, nil
}

func (k Keeper) getTokenByDenom(ctx sdk.Context, denom string) (token tokentypes.FanToken, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(tokentypes.KeyDenom(denom))
	if bz == nil {
		return token, sdkerrors.Wrap(tokentypes.ErrTokenNotExists, fmt.Sprintf("token denom %s does not exist", denom))
	}

	var symbol gogotypes.StringValue
	k.cdc.MustUnmarshalBinaryBare(bz, &symbol)

	token, err = k.getFanTokenBySymbol(ctx, symbol.Value)
	if err != nil {
		return token, err
	}

	return token, nil
}

// reset all indices by the new owner for token query
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, symbol string, srcOwner, dstOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(tokentypes.KeyFanTokens(srcOwner, symbol))

	// add the new key
	k.setWithOwner(ctx, dstOwner, symbol)
}

// getTokenSupply queries the token supply from the total supply
func (k Keeper) getFanTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}
