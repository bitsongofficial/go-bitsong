package keeper

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/content/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the track store
type Keeper struct {
	bankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(bk bank.Keeper, cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		bankKeeper: bk,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetContent get Content from store by Uri
func (k Keeper) GetContent(ctx sdk.Context, uri string) (content types.Content, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetContentKey(uri))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &content)
	return content, true
}

func (k Keeper) SetContent(ctx sdk.Context, content types.Content) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(content)
	store.Set(types.GetContentKey(content.Uri), bz)
}

func (k Keeper) GetDenom(ctx sdk.Context, denom string) (_ string, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenomKey(denom))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &denom)
	return denom, true
}

func (k Keeper) SetDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(denom)
	store.Set(types.GetDenomKey(denom), bz)
}

func (k Keeper) IterateContents(ctx sdk.Context, fn func(content types.Content) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ContentKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var content types.Content
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &content)
		if fn(content) {
			break
		}
	}
}

func (k Keeper) GetContents(ctx sdk.Context) []types.Content {
	var contents []types.Content
	k.IterateContents(ctx, func(content types.Content) (stop bool) {
		contents = append(contents, content)
		return false
	})
	return contents
}

func (k Keeper) Add(ctx sdk.Context, content types.Content) (string, error) {
	_, uriExists := k.GetContent(ctx, content.Uri)
	if uriExists {
		return "", fmt.Errorf("uri %s is not avalable", content.Uri)
	}

	// check if denom is duplicated
	_, denomExists := k.GetDenom(ctx, content.Denom)
	if denomExists {
		return "", fmt.Errorf("denom %s is not avalable", content.Denom)
	}
	k.SetDenom(ctx, content.Denom)

	content.CreatedAt = ctx.BlockHeader().Time
	k.SetContent(ctx, content)

	return content.Uri, nil
}

func (k Keeper) Stream(ctx sdk.Context, uri string, from sdk.AccAddress) error {
	// get content
	content, uriExists := k.GetContent(ctx, uri)
	if !uriExists {
		return fmt.Errorf("uri %s is not avalable", uri)
	}

	// subtract stream-price from requester
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(content.StreamPrice))
	if err != nil {
		return err
	}

	// update content with new rewards
	content = allocateFundsRightsHolders(content, content.StreamPrice)

	// mint stream to requester (1 * 10^0)
	unit := sdk.NewInt(1)
	coin := sdk.NewCoin(content.Denom, unit)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	content = increaseTotalSupply(content, coin)
	k.SetContent(ctx, content)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	return nil
}

func increaseTotalSupply(cnt types.Content, coin sdk.Coin) types.Content {
	cnt.TotalSupply = cnt.TotalSupply.Add(coin)
	return cnt
}

func allocateFundsRightsHolders(cnt types.Content, coin sdk.Coin) types.Content {
	// allocate funds to rights holders
	for i, rh := range cnt.RightsHolders {
		price := sdk.NewDecCoinFromCoin(coin)
		allocation := price.Amount.Quo(sdk.NewDec(100).Quo(rh.Quota))
		cnt.RightsHolders[i].Rewards = rh.Rewards.Add(sdk.NewDecCoinFromDec(btsg.BondDenom, allocation))
	}

	// increase volume
	cnt.Volume = cnt.Volume.Add(coin)

	return cnt
}

func (k Keeper) Download(ctx sdk.Context, uri string, from sdk.AccAddress) error {
	// get content
	content, uriExists := k.GetContent(ctx, uri)
	if !uriExists {
		return fmt.Errorf("uri %s is not avalable", uri)
	}

	// subtract download-price from requester
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(content.DownloadPrice))
	if err != nil {
		return err
	}

	// update content with new rewards
	content = allocateFundsRightsHolders(content, content.DownloadPrice)

	// mint download to requester (1 * 10^6)
	unit := sdk.NewInt(1000000)
	coin := sdk.NewCoin(content.Denom, unit)

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	content = increaseTotalSupply(content, coin)
	k.SetContent(ctx, content)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	return nil
}
