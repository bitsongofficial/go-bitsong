package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey sdk.StoreKey
	codec    *codec.Codec

	accountKeeper auth.AccountKeeper
}

func NewKeeper(storeKey sdk.StoreKey, codec *codec.Codec, accountKeeper auth.AccountKeeper) Keeper {
	keeper := Keeper{
		storeKey:      storeKey,
		codec:         codec,
		accountKeeper: accountKeeper,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetChannel(ctx sdk.Context, handle string) (channel types.Channel, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetChannelKey(handle))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &channel)
	return channel, true
}

func (k Keeper) GetChannelByOwner(ctx sdk.Context, addr sdk.AccAddress) (channel types.Channel, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOwnerKey(addr))
	if bz == nil {
		return
	}
	var handle string
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &handle)

	bz = store.Get(types.GetChannelKey(handle))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &channel)
	return channel, true
}

func (k Keeper) SetChannel(ctx sdk.Context, channel types.Channel) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&channel)
	store.Set(types.GetChannelKey(channel.Handle), bz)

	handleBz := k.codec.MustMarshalBinaryLengthPrefixed(&channel.Handle)
	store.Set(types.GetOwnerKey(channel.Owner), handleBz)
}

func (k Keeper) IterateAllChannels(ctx sdk.Context, fn func(channel types.Channel) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ChannelKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var channel types.Channel
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &channel)
		if fn(channel) {
			break
		}
	}
}

func (k Keeper) GetAllChannels(ctx sdk.Context) []types.Channel {
	var channels []types.Channel
	k.IterateAllChannels(ctx, func(channel types.Channel) (stop bool) {
		channels = append(channels, channel)
		return false
	})
	return channels
}

func (k Keeper) IsHandleDuplicated(ctx sdk.Context, handle string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetChannelKey(handle))
}

func (k Keeper) IsChannelPresent(ctx sdk.Context, owner sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOwnerKey(owner))
}

func (k Keeper) CreateChannel(ctx sdk.Context, owner sdk.AccAddress, handle, metadataURI string) (channel types.Channel, err error) {
	if k.IsHandleDuplicated(ctx, handle) {
		return channel, sdkerrors.Wrap(types.ErrChannelCreateError, fmt.Sprintf("handle %s exist", handle))
	}

	if k.IsChannelPresent(ctx, owner) {
		return channel, sdkerrors.Wrap(types.ErrChannelCreateError, fmt.Sprintf("handle exist on account %s", owner.String()))
	}

	channel = types.NewChannel(owner, handle, metadataURI, ctx.BlockHeader().Time)
	k.SetChannel(ctx, channel)

	k.Logger(ctx).Info(fmt.Sprintf("Channel Created %s", channel.String()))

	return channel, nil
}

func (k Keeper) EditChannel(ctx sdk.Context, owner sdk.AccAddress, metadataURI string) (channel types.Channel, err error) {
	channel, found := k.GetChannelByOwner(ctx, owner)
	if !found {
		return channel, sdkerrors.Wrap(types.ErrChannelNotFound, fmt.Sprintf("channel not found"))
	}

	channel.MetadataURI = metadataURI
	k.SetChannel(ctx, channel)

	k.Logger(ctx).Info(fmt.Sprintf("Channel Edited %s", channel.String()))

	return channel, nil
}
