package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryAlbums:
			return queryAlbums(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown album query endpoint")
		}
	}
}

// nolint: unparam
func queryAlbums(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryAlbumsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	albums := keeper.GetAlbumsFiltered(ctx, params.Owner, params.AlbumStatus, params.Limit)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, albums)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
