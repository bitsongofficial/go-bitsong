package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateNFT(goCtx context.Context, msg *types.MsgCreateNFT) (*types.MsgCreateNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// create metadata
	metadataId := m.Keeper.GetLastMetadataId(ctx) + 1
	m.Keeper.SetLastMetadataId(ctx, metadataId)
	msg.Metadata.Id = metadataId
	for index := range msg.Metadata.Data.Creators {
		msg.Metadata.Data.Creators[index].Verified = false
	}
	m.Keeper.SetMetadata(ctx, msg.Metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMetadataCreation{
		Creator:    msg.Sender,
		MetadataId: msg.Metadata.Id,
	})

	// burn fees before minting an nft
	fee := m.GetParamSet(ctx).IssuePrice
	if fee.IsPositive() {
		feeCoins := sdk.Coins{fee}
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return nil, err
		}
		err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, feeCoins)
		if err != nil {
			return nil, err
		}
		err = m.bankKeeper.BurnCoins(ctx, types.ModuleName, feeCoins)
		if err != nil {
			return nil, err
		}
	}

	// create nft
	nftId := m.Keeper.GetLastNftId(ctx) + 1
	m.Keeper.SetLastNftId(ctx, nftId)
	nft := types.NFT{
		Id:         nftId,
		Owner:      msg.Sender,
		MetadataId: metadataId,
	}
	m.Keeper.SetNFT(ctx, nft)
	ctx.EventManager().EmitTypedEvent(&types.EventNFTCreation{
		Creator: msg.Sender,
		NftId:   nftId,
	})

	return &types.MsgCreateNFTResponse{
		Id:         nftId,
		MetadataId: metadataId,
	}, nil
}

func (m msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.TransferNFT(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgTransferNFTResponse{}, nil
}

func (m msgServer) SignMetadata(goCtx context.Context, msg *types.MsgSignMetadata) (*types.MsgSignMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	metadata, err := m.Keeper.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return nil, err
	}

	isCreator := false
	for index, creator := range metadata.Data.Creators {
		if creator.Address == msg.Sender {
			metadata.Data.Creators[index].Verified = true
			isCreator = true
		}
	}

	if isCreator == false {
		return nil, types.ErrNotEnoughPermission
	}

	m.Keeper.SetMetadata(ctx, metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMetadataSign{
		Signer:     msg.Sender,
		MetadataId: msg.MetadataId,
	})

	return &types.MsgSignMetadataResponse{}, nil
}

func (m msgServer) UpdateMetadata(goCtx context.Context, msg *types.MsgUpdateMetadata) (*types.MsgUpdateMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	metadata, err := m.Keeper.GetMetadataById(ctx, msg.MetadataId)
	if err != nil {
		return nil, err
	}

	if !metadata.IsMutable {
		return nil, types.ErrMetadataImmutable
	}

	if metadata.UpdateAuthority != msg.Sender {
		return nil, types.ErrNotEnoughPermission
	}

	metadata.Data = msg.Data
	for index := range metadata.Data.Creators {
		metadata.Data.Creators[index].Verified = false
	}
	m.Keeper.SetMetadata(ctx, metadata)
	ctx.EventManager().EmitTypedEvent(&types.EventMetadataUpdate{
		Updater:    msg.Sender,
		MetadataId: metadata.Id,
	})

	return &types.MsgUpdateMetadataResponse{}, nil
}

func (m msgServer) UpdateMetadataAuthority(goCtx context.Context, msg *types.MsgUpdateMetadataAuthority) (*types.MsgUpdateMetadataAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.UpdateMetadataAuthority(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateMetadataAuthorityResponse{}, nil
}

func (m msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collectionId := m.Keeper.GetLastCollectionId(ctx) + 1
	m.Keeper.SetLastCollectionId(ctx, collectionId)

	collection := types.Collection{
		Id:              collectionId,
		Name:            msg.Name,
		Uri:             msg.Uri,
		UpdateAuthority: msg.UpdateAuthority,
	}
	m.Keeper.SetCollection(ctx, collection)
	ctx.EventManager().EmitTypedEvent(&types.EventCollectionCreation{
		Creator:      msg.Sender,
		CollectionId: collection.Id,
	})

	return &types.MsgCreateCollectionResponse{
		Id: collectionId,
	}, nil
}

func (m msgServer) VerifyCollection(goCtx context.Context, msg *types.MsgVerifyCollection) (*types.MsgVerifyCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := m.Keeper.GetCollectionById(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}
	if collection.UpdateAuthority != msg.Sender {
		return nil, types.ErrNotEnoughPermission
	}
	if m.Keeper.GetLastCollectionId(ctx) < msg.CollectionId {
		return nil, types.ErrCollectionDoesNotExist
	}
	if m.Keeper.GetLastNftId(ctx) < msg.NftId {
		return nil, types.ErrNFTDoesNotExist
	}

	m.Keeper.SetCollectionNftRecord(ctx, msg.CollectionId, msg.NftId)
	ctx.EventManager().EmitTypedEvent(&types.EventCollectionVerification{
		Verifier:     msg.Sender,
		CollectionId: msg.CollectionId,
		NftId:        msg.NftId,
	})

	return &types.MsgVerifyCollectionResponse{}, nil
}

func (m msgServer) UnverifyCollection(goCtx context.Context, msg *types.MsgUnverifyCollection) (*types.MsgUnverifyCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := m.Keeper.GetCollectionById(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}
	if collection.UpdateAuthority != msg.Sender {
		return nil, types.ErrNotEnoughPermission
	}

	if m.Keeper.GetLastCollectionId(ctx) < msg.CollectionId {
		return nil, types.ErrCollectionDoesNotExist
	}
	if m.Keeper.GetLastNftId(ctx) < msg.NftId {
		return nil, types.ErrNFTDoesNotExist
	}

	m.Keeper.DeleteCollectionNftRecord(ctx, msg.CollectionId, msg.NftId)
	ctx.EventManager().EmitTypedEvent(&types.EventCollectionUnverification{
		Verifier:     msg.Sender,
		CollectionId: msg.CollectionId,
		NftId:        msg.NftId,
	})

	return &types.MsgUnverifyCollectionResponse{}, nil
}

func (m msgServer) UpdateCollectionAuthority(goCtx context.Context, msg *types.MsgUpdateCollectionAuthority) (*types.MsgUpdateCollectionAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := m.Keeper.GetCollectionById(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}
	if collection.UpdateAuthority != msg.Sender {
		return nil, types.ErrNotEnoughPermission
	}

	collection.UpdateAuthority = msg.NewAuthority
	m.Keeper.SetCollection(ctx, collection)
	ctx.EventManager().EmitTypedEvent(&types.EventUpdateCollectionAuthority{
		CollectionId: msg.CollectionId,
		NewAuthority: msg.NewAuthority,
	})

	return &types.MsgUpdateCollectionAuthorityResponse{}, nil
}
