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

	nft, err := k.GetNFTById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	metadata, err := k.GetMetadataById(ctx, nft.CollId, nft.MetadataId)
	if err != nil {
		return nil, err
	}
	return &types.QueryNFTInfoResponse{
		Nft:      nft,
		Metadata: metadata,
	}, nil
}

func (k Keeper) NFTsByOwner(c context.Context, req *types.QueryNFTsByOwnerRequest) (*types.QueryNFTsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}

	nfts := k.GetNFTsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	metadata := []types.Metadata{}
	for _, nft := range nfts {
		meta, err := k.GetMetadataById(ctx, nft.CollId, nft.MetadataId)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}
	return &types.QueryNFTsByOwnerResponse{
		Nfts:     nfts,
		Metadata: metadata,
	}, nil
}

func (k Keeper) Metadata(c context.Context, req *types.QueryMetadataRequest) (*types.QueryMetadataResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	metadata, err := k.GetMetadataById(ctx, req.CollId, req.Id)
	if err != nil {
		return nil, err
	}
	return &types.QueryMetadataResponse{
		Metadata: metadata,
	}, nil
}

func (k Keeper) Collection(c context.Context, req *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	collection, err := k.GetCollectionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	nfts := k.GetCollectionNfts(ctx, req.Id)
	return &types.QueryCollectionResponse{
		Collection: collection,
		Nfts:       nfts,
	}, nil
}
