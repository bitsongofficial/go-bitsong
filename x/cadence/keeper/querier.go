package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

var _ types.QueryServer = &Querier{}

type Querier struct {
	keeper *Keeper
}

func NewQuerier(k *Keeper) Querier {
	return Querier{
		keeper: k,
	}
}

// ContractModules returns contract addresses which are using the cadence
func (q Querier) CadenceContracts(stdCtx context.Context, req *types.QueryCadenceContracts) (*types.QueryCadenceContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	contracts, err := q.keeper.GetPaginatedContracts(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

// CadenceContract returns the cadence contract information
func (q Querier) CadenceContract(stdCtx context.Context, req *types.QueryCadenceContract) (*types.QueryCadenceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	// Ensure the contract address is valid
	if _, err := sdk.AccAddressFromBech32(req.ContractAddress); err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	contract, err := q.keeper.GetCadenceContract(ctx, req.ContractAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryCadenceContractResponse{
		CadenceContract: *contract,
	}, nil
}

// Params returns the total set of cadence parameters.
func (q Querier) Params(stdCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	p := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &p,
	}, nil
}
