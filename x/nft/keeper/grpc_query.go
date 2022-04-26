package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) NFTInfo(c context.Context, req *types.QueryNFTInfoRequest) (*types.QueryNFTInfoResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	// TODO: implement!
	return &types.QueryNFTInfoResponse{}, nil
}

func (k Keeper) Metadata(c context.Context, req *types.QueryMetadataRequest) (*types.QueryMetadataResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	// TODO: implement!
	return &types.QueryMetadataResponse{}, nil
}

func (k Keeper) Collection(c context.Context, req *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx

	// TODO: implement!
	return &types.QueryCollectionResponse{}, nil
}
