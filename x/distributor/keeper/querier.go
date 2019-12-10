package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryDistributors:
			return queryDistributors(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown reward query endpoint")
		}
	}
}

func queryDistributors(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	distributors := keeper.GetAllDistributors(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, distributors)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
