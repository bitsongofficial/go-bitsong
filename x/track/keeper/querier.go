package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryTracks:
			return queryTracks(ctx, req, keeper)
		case types.QueryTrack:
			return queryTrack(ctx, req, keeper)
		case types.QueryPlays:
			return queryPlays(ctx, req, keeper)
		case types.QueryShares:
			return queryShares(ctx, keeper)
		case types.QueryDeposits:
			return queryDeposits(ctx, req, keeper)
		default:
			return nil, sdkerr.Wrap(sdkerr.ErrUnknownRequest, "errunknown track query endpoint")
		}
	}
}

// nolint: unparam
func queryTracks(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTracksParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerr.Wrap(err, "incorrectly formatted request data")
	}

	tracks := keeper.GetTracksFiltered(ctx, params.Owner, params.TrackStatus, params.Limit)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tracks)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}

// nolint: unparam
func queryTrack(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTrackParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerr.Wrap(err, "incorrectly formatted request data")
	}

	track, ok := keeper.GetTrack(ctx, params.TrackID)
	if !ok {
		return nil, types.ErrUnknownTrack(types.DefaultCodespace, fmt.Sprintf("unknown track-id %d", params.TrackID))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, track)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}

// nolint: unparam
func queryPlays(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTrackParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)

	if err != nil {
		return nil, sdkerr.Wrap(err, "incorrectly formatted request data")
	}

	plays := keeper.GetPlays(ctx, params.TrackID)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, plays)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}

func queryShares(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	shares := keeper.GetAllShares(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, shares)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}

func queryDeposits(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTrackParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerr.Wrap(err, "incorrectly formatted request data")
	}

	deposits := keeper.GetDeposits(ctx, params.TrackID)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, deposits)
	if err != nil {
		return nil, sdkerr.Wrap(err, "could not marshal result to JSON")
	}
	return bz, nil
}
