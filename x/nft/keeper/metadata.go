package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) GetAllLastMetadataIds(ctx sdk.Context) []types.LastMetadataIdInfo {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, types.KeyPrefixLastMetadataId)
	defer it.Close()

	infos := []types.LastMetadataIdInfo{}
	for ; it.Valid(); it.Next() {
		var info types.LastMetadataIdInfo
		k.cdc.MustUnmarshal(it.Value(), &info)

		infos = append(infos, info)
	}
	return infos
}

func (k Keeper) GetLastMetadataId(ctx sdk.Context, collId uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastMetadataId(collId))
	if bz == nil {
		return 0
	}

	var info types.LastMetadataIdInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info.LastMetadataId
}

func (k Keeper) SetLastMetadataId(ctx sdk.Context, collId, id uint64) {
	info := types.LastMetadataIdInfo{
		CollId:         collId,
		LastMetadataId: id,
	}
	bz := k.cdc.MustMarshal(&info)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastMetadataId(collId), bz)
}

func (k Keeper) GetMetadataById(ctx sdk.Context, collId, id uint64) (types.Metadata, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MetadataId(collId, id))
	if bz == nil {
		return types.Metadata{}, sdkerrors.Wrapf(types.ErrMetadataDoesNotExist, "metadata: %d/%d does not exist", collId, id)
	}
	metadata := types.Metadata{}
	k.cdc.MustUnmarshal(bz, &metadata)
	return metadata, nil
}

func (k Keeper) SetMetadata(ctx sdk.Context, metadata types.Metadata) {
	bz := k.cdc.MustMarshal(&metadata)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MetadataId(metadata.CollId, metadata.Id), bz)
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

func (k Keeper) SetPrimarySaleHappened(ctx sdk.Context, collId, metadataId uint64) error {
	metadata, err := k.GetMetadataById(ctx, collId, metadataId)
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
	metadata, err := k.GetMetadataById(ctx, msg.CollId, msg.MetadataId)
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
	metadata, err := k.GetMetadataById(ctx, msg.CollId, msg.MetadataId)
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
