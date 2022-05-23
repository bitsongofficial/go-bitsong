package keeper

import (
	"context"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Merkledrop(c context.Context, req *types.QueryMerkledropRequest) (*types.QueryMerkledropResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	merkledrop, err := k.GetMerkleDropById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.QueryMerkledropResponse{
		Merkledrop: merkledrop,
	}, nil
}

func (k Keeper) IndexClaimed(c context.Context, req *types.QueryIndexClaimedRequest) (*types.QueryIndexClaimedResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryIndexClaimedResponse{
		IsClaimed: k.IsClaimed(ctx, req.Id, req.Index),
	}, nil
}
