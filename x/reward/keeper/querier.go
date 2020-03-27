package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bitsongofficial/go-bitsong/x/reward/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryRewards:
			return queryRewards(ctx, keeper)
		default:
			return nil, sdkerr.Wrap(sdkerr.ErrUnknownRequest, "unknown reward query endpoint")
		}
	}
}

func queryRewards(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	rewards := keeper.GetAllRewards(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rewards)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}
