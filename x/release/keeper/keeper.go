package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey sdk.StoreKey
	codec    *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, codec *codec.Codec) Keeper {
	keeper := Keeper{
		storeKey: storeKey,
		codec:    codec,
	}

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetRelease(ctx sdk.Context, releaseID string) (release types.Release, ok bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetReleaseKey(releaseID))
	if bz == nil {
		return
	}
	k.codec.MustUnmarshalBinaryLengthPrefixed(bz, &release)
	return release, true
}

func (k Keeper) SetRelease(ctx sdk.Context, release types.Release) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&release)
	store.Set(types.GetReleaseKey(release.ReleaseID), bz)
}

func (k Keeper) IterateAllReleases(ctx sdk.Context, fn func(release types.Release) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ReleaseKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var release types.Release
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &release)
		if fn(release) {
			break
		}
	}
}

func (k Keeper) GetAllReleases(ctx sdk.Context) []types.Release {
	var releases []types.Release
	k.IterateAllReleases(ctx, func(release types.Release) (stop bool) {
		releases = append(releases, release)
		return false
	})
	return releases
}

func (k Keeper) SetReleaseForCreator(ctx sdk.Context, release types.Release) {
	store := ctx.KVStore(k.storeKey)
	bz := k.codec.MustMarshalBinaryLengthPrefixed(&release.ReleaseID)
	store.Set(types.ReleaseAddressKey(release.Creator, release.ReleaseID), bz)
}

func (k Keeper) IterateAllReleasesForCreator(ctx sdk.Context, creator sdk.AccAddress, fn func(releaseID string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetReleaseForCreatorKey(creator))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var releaseID string
		k.codec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &releaseID)
		if fn(releaseID) {
			break
		}
	}
}

func (k Keeper) GetAllReleaseForCreator(ctx sdk.Context, creator sdk.AccAddress) []types.Release {
	var releases []types.Release
	k.IterateAllReleasesForCreator(ctx, creator, func(releaseID string) (stop bool) {
		release, _ := k.GetRelease(ctx, releaseID)
		releases = append(releases, release)
		return false
	})
	return releases
}

func (k Keeper) IsReleaseDuplicated(ctx sdk.Context, releaseID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetReleaseKey(releaseID))
}

func (k Keeper) CreateRelease(ctx sdk.Context, address sdk.AccAddress, releaseID, metadataURI string) (release types.Release, err error) {
	if k.IsReleaseDuplicated(ctx, releaseID) {
		return release, sdkerrors.Wrap(types.ErrReleaseCreateError, fmt.Sprintf("releaseID %s exist", releaseID))
	}

	release = types.NewRelease(releaseID, metadataURI, address, ctx.BlockHeader().Time)
	k.SetRelease(ctx, release)
	k.SetReleaseForCreator(ctx, release)

	k.Logger(ctx).Info(fmt.Sprintf("Release Created %s", release.String()))

	return release, nil
}
