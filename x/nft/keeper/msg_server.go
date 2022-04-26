package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
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
	// TODO: implement!
	return &types.MsgCreateNFTResponse{}, nil
}

func (m msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	// TODO: implement!
	return &types.MsgTransferNFTResponse{}, nil
}

func (m msgServer) SignMetadata(goCtx context.Context, msg *types.MsgSignMetadata) (*types.MsgSignMetadataResponse, error) {
	// TODO: implement!
	return &types.MsgSignMetadataResponse{}, nil
}

func (m msgServer) UpdateMetadata(goCtx context.Context, msg *types.MsgUpdateMetadata) (*types.MsgUpdateMetadataResponse, error) {
	// TODO: implement!
	return &types.MsgUpdateMetadataResponse{}, nil
}

func (m msgServer) UpdateMetadataAuthority(goCtx context.Context, msg *types.MsgUpdateMetadataAuthority) (*types.MsgUpdateMetadataAuthorityResponse, error) {
	// TODO: implement!
	return &types.MsgUpdateMetadataAuthorityResponse{}, nil
}

func (m msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	// TODO: implement!
	return &types.MsgCreateCollectionResponse{}, nil
}

func (m msgServer) VerifyCollection(goCtx context.Context, msg *types.MsgVerifyCollection) (*types.MsgVerifyCollectionResponse, error) {
	// TODO: implement!
	return &types.MsgVerifyCollectionResponse{}, nil
}

func (m msgServer) UnverifyCollection(goCtx context.Context, msg *types.MsgUnverifyCollection) (*types.MsgUnverifyCollectionResponse, error) {
	// TODO: implement!
	return &types.MsgUnverifyCollectionResponse{}, nil
}

func (m msgServer) UpdateCollectionAuthority(goCtx context.Context, msg *types.MsgUpdateCollectionAuthority) (*types.MsgUpdateCollectionAuthorityResponse, error) {
	// TODO: implement!
	return &types.MsgUpdateCollectionAuthorityResponse{}, nil
}
