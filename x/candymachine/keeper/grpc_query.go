package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
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

func (k Keeper) CandyMachines(c context.Context, req *types.QueryCandyMachinesRequest) (*types.QueryCandyMachinesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	machines := k.GetAllCandyMachines(ctx)
	return &types.QueryCandyMachinesResponse{
		Machines: machines,
	}, nil
}

func (k Keeper) CandyMachine(c context.Context, req *types.QueryCandyMachineRequest) (*types.QueryCandyMachineResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	machine, err := k.GetCandyMachineByCollId(ctx, req.CollId)
	if err != nil {
		return nil, err
	}
	return &types.QueryCandyMachineResponse{
		Machine: machine,
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
