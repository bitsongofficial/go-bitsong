package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryRelease:
			return queryRelease(ctx, path[1:], req, k)
		case types.QueryAllReleaseForCreator:
			return queryAllReleaseForCreator(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown profile query endpoint")
		}
	}
}

func queryRelease(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryReleaseParams

	err := k.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	release, found := k.GetRelease(ctx, params.ReleaseID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrReleaseNotFound, fmt.Sprintf("release %s not found", params.ReleaseID))
	}

	bz, err := k.codec.MarshalJSON(release)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAllReleaseForCreator(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryAllReleaseForCreatorParams

	err := k.codec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	releases := k.GetAllReleaseForCreator(ctx, params.Creator)

	bz, err := k.codec.MarshalJSON(releases)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
