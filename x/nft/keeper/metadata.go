package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetLastMetadataId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastMetadataId)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLastMetadataId(ctx sdk.Context, id uint64) {
	idBz := sdk.Uint64ToBigEndian(id)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastMetadataId, idBz)
}

func (k Keeper) GetMetadataById(ctx sdk.Context, id uint64) (types.Metadata, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.PrefixMetadata, sdk.Uint64ToBigEndian(id)...))
	if bz == nil {
		return types.Metadata{}, sdkerrors.Wrapf(types.ErrMetadataDoesNotExist, "metadata: %d does not exist", id)
	}
	metadata := types.Metadata{}
	k.cdc.MustUnmarshal(bz, &metadata)
	return metadata, nil
}

func (k Keeper) SetMetadata(ctx sdk.Context, metadata types.Metadata) {
	idBz := sdk.Uint64ToBigEndian(metadata.Id)
	bz := k.cdc.MustMarshal(&metadata)
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.PrefixMetadata, idBz...), bz)
}
