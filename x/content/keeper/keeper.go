package keeper

import (
	"fmt"
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
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &content)
	return content, true
}

func (k Keeper) SetContent(ctx sdk.Context, content *types.Content) {
	store := ctx.KVStore(k.storeKey)
	fmt.Println(content.String())
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&content)
	store.Set(types.GetContentKey(content.Uri), bz)
	fmt.Println(content.Uri)
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

func (k Keeper) Add(ctx sdk.Context, content *types.Content) (string, error) {
	_, uriExists := k.GetContent(ctx, content.Uri)
	if uriExists {
		return "", fmt.Errorf("uri %s is not avalable", content.Uri)
	}
	content.CreatedAt = ctx.BlockHeader().Time
	k.SetContent(ctx, content)
	fmt.Println(content.Uri)

	return content.Uri, nil
}

func (k Keeper) Action(ctx sdk.Context, uri string, from sdk.AccAddress) error {
	// get content
	_, uriExists := k.GetContent(ctx, uri)
	if !uriExists {
		return fmt.Errorf("uri %s is not avalable", uri)
	}

	// subtract stream-price from requester
	//err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(content.StreamPrice))
	//if err != nil {
	//	return err
	//}

	// update content with new rewards
	//content = allocateDaoFunds(content, content.StreamPrice)

	//content.TotalStreams = content.TotalStreams + 1
	//k.SetContent(ctx, content)

	return nil
}

/*func allocateDaoFunds(cnt types.Content, coin sdk.Coin) types.Content {
	// allocate dao funds
	for i, rh := range cnt.RightsHolders {
		price := sdk.NewDecCoinFromCoin(coin)
		allocation := price.Amount.Quo(sdk.NewDec(100).Quo(rh.Quota))
		cnt.RightsHolders[i].Rewards = rh.Rewards.Add(sdk.NewDecCoinFromDec(btsg.BondDenom, allocation))
	}

	return cnt
}*/
