package track

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the track Querier
const (
	//QueryResolve = "resolve"
	//QueryWhois   = "whois"
	QueryAllSongsByAddress   = "songs"
)

type QuerySongsParams struct {
	Owner sdk.AccAddress
}

func NewQuerySongsParams(addr sdk.AccAddress) QuerySongsParams {
	return QuerySongsParams{
		Owner: addr,
	}
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryAllSongsByAddress:
			return queryOrdersByAddress(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown track query endpoint")
		}
	}
}

func queryOrdersByAddress(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QuerySongsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	songs, err := keeper.GetSongsByAddr(ctx, params.Owner)
	if err != nil {
		return nil, sdk.NewError(DefaultCodespace, CodeSongNotExist, err.Error())
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, songs)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("could not marshal result to JSON: %s", err))
	}

	return bz, nil
}

/*func querySearch(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	value := keeper.SearchTitle(ctx, path[0])

	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not search title")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResSearch{value})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}*/