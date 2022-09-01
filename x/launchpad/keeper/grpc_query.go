package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(ctx),
	}, nil
}

func (k Keeper) LaunchPads(c context.Context, req *types.QueryLaunchPadsRequest) (*types.QueryLaunchPadsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pads := k.GetAllLaunchPads(ctx)
	return &types.QueryLaunchPadsResponse{
		Pads: pads,
	}, nil
}

func (k Keeper) LaunchPad(c context.Context, req *types.QueryLaunchPadRequest) (*types.QueryLaunchPadResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pad, err := k.GetLaunchPadByCollId(ctx, req.CollId)
	if err != nil {
		return nil, err
	}
	return &types.QueryLaunchPadResponse{
		Pad: pad,
	}, nil
}

func (k Keeper) MintableMetadataIds(c context.Context, req *types.QueryMintableMetadataIdsRequest) (*types.QueryMintableMetadataIdsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	metadataIds := k.GetMintableMetadataIds(ctx, req.CollId)

	return &types.QueryMintableMetadataIdsResponse{
		Info: types.MintableMetadataIds{
			CollectionId:        req.CollId,
			MintableMetadataIds: metadataIds,
		},
	}, nil
}
