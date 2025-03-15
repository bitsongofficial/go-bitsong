package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

var _ types.QueryServer = &Querier{}

type Querier struct {
	keeper Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{
		keeper: k,
	}
}

// CadanceContracts returns contract addresses which are using the cadance
func (q Querier) CadanceContracts(stdCtx context.Context, req *types.QueryCadanceContracts) (*types.QueryCadanceContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	contracts, err := q.keeper.GetPaginatedContracts(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

// CadanceContract returns the cadance contract information
func (q Querier) CadanceContract(stdCtx context.Context, req *types.QueryCadanceContract) (*types.QueryCadanceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	// Ensure the contract address is valid
	if _, err := sdk.AccAddressFromBech32(req.ContractAddress); err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	contract, err := q.keeper.GetCadanceContract(ctx, req.ContractAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryCadanceContractResponse{
		CadanceContract: *contract,
	}, nil
}

// Params returns the total set of cadance parameters.
func (q Querier) Params(stdCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(stdCtx)

	p := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &p,
	}, nil
}
