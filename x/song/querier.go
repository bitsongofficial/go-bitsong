package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	//QueryResolve = "resolve"
	//QueryWhois   = "whois"
	QueryTitles   = "titles"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryTitles:
			return queryTitles(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown song query endpoint")
		}
	}
}

func queryTitles(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var titlesList QueryResTitles

	iterator := keeper.GetTitlesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		titlesList = append(titlesList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, titlesList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}