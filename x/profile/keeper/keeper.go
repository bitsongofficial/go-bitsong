package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/profile/types"
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

func (k Keeper) GetProfile(ctx sdk.Context, handle string) (profile types.Profile, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProfileKey(handle))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &profile)
	return profile, true
}

func (k Keeper) GetProfileByAddress(ctx sdk.Context, addr sdk.AccAddress) (profile types.Profile, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAddressKey(addr))
	if bz == nil {
		return
	}
	var handle string
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &handle)

	bz = store.Get(types.GetProfileKey(handle))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &profile)
	return profile, true
}

func (k Keeper) SetProfile(ctx sdk.Context, acc types.Profile) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&acc)
	store.Set(types.GetProfileKey(acc.Handle), bz)

	handleBz := k.codec.MustMarshalBinaryLengthPrefixed(&acc.Handle)
	store.Set(types.GetAddressKey(acc.Address), handleBz)
}

func (k Keeper) IterateAllProfiles(ctx sdk.Context, fn func(profile types.Profile) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProfileKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var profile types.Profile
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &profile)
		if fn(profile) {
			break
		}
	}
}

func (k Keeper) GetAllProfiles(ctx sdk.Context) []types.Profile {
	var profiles []types.Profile
	k.IterateAllProfiles(ctx, func(profile types.Profile) (stop bool) {
		profiles = append(profiles, profile)
		return false
	})
	return profiles
}

func (k Keeper) IsHandleDuplicated(ctx sdk.Context, handle string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetProfileKey(handle))
}

func (k Keeper) IsAccPresent(ctx sdk.Context, accAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetAddressKey(accAddr))
}

func (k Keeper) CreateProfile(ctx sdk.Context, address sdk.AccAddress, handle, metadataURI string) (profile types.Profile, err error) {
	if k.IsHandleDuplicated(ctx, handle) {
		return profile, sdkerrors.Wrap(types.ErrProfileCreateError, fmt.Sprintf("handle %s exist", handle))
	}

	if k.IsAccPresent(ctx, address) {
		return profile, sdkerrors.Wrap(types.ErrProfileCreateError, fmt.Sprintf("handle exist on account %s", address.String()))
	}

	profile = types.NewProfile(address, handle, metadataURI, ctx.BlockHeader().Time)
	k.SetProfile(ctx, profile)

	k.Logger(ctx).Info(fmt.Sprintf("Profile Created %s", profile.String()))

	return profile, nil
}
