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

func (k Keeper) GetAllMetadata(ctx sdk.Context) []types.Metadata {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.PrefixMetadata)
	defer it.Close()

	allMetadata := []types.Metadata{}
	for ; it.Valid(); it.Next() {
		var metadata types.Metadata
		k.cdc.MustUnmarshal(it.Value(), &metadata)

		allMetadata = append(allMetadata, metadata)
	}

	return allMetadata
}

func (k Keeper) SetPrimarySaleHappened(ctx sdk.Context, metadataId uint64) error {
	metadata, err := k.GetMetadataById(ctx, metadataId)
	if err != nil {
		return err
	}

	if metadata.PrimarySaleHappened == true {
		return types.ErrPrimarySaleAlreadyHappened
	}
	metadata.PrimarySaleHappened = true
	metadata.IsMutable = false
	k.SetMetadata(ctx, metadata)
	return nil
}

func (k Keeper) UpdateMetadataAuthority(ctx sdk.Context, msg *types.MsgUpdateMetadataAuthority) error {
	metadata, err := k.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return err
	}

	if metadata.MetadataAuthority != msg.Sender {
		return types.ErrNotEnoughPermission
	}

	metadata.MetadataAuthority = msg.NewAuthority
	k.SetMetadata(ctx, metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMetadataAuthorityUpdate{
		MetadataId:   msg.Sender,
		NewAuthority: msg.NewAuthority,
	})
	return nil
}

func (k Keeper) UpdateMintAuthority(ctx sdk.Context, msg *types.MsgUpdateMintAuthority) error {
	metadata, err := k.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return err
	}

	if metadata.MintAuthority != msg.Sender {
		return types.ErrNotEnoughPermission
	}

	metadata.MintAuthority = msg.NewAuthority
	k.SetMetadata(ctx, metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMintAuthorityUpdate{
		MetadataId:   msg.Sender,
		NewAuthority: msg.NewAuthority,
	})
	return nil
}
