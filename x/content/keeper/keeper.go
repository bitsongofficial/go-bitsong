package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/content/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the track store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

// NewKeeper creates a content keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
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

func (k Keeper) SetContent(ctx sdk.Context, content types.Content) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(content)
	store.Set(types.GetContentKey(content.Uri), bz)
}

func (k Keeper) Add(ctx sdk.Context, content types.Content) (string, error) {
	_, uriExists := k.GetContent(ctx, content.Uri)
	if uriExists {
		return "", fmt.Errorf("uri %s is not avalable", content.Uri)
	}

	content.CreatedAt = ctx.BlockHeader().Time
	k.SetContent(ctx, content)

	return content.Uri, nil
}
