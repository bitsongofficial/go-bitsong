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
		/*case types.QueryProfileByAddress:
		return queryProfileByAddress(ctx, path[1:], req, k)*/
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown profile query endpoint")
		}
	}
}

func queryRelease(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryReleaseParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	release, found := k.GetRelease(ctx, params.ReleaseID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrReleaseNotFound, fmt.Sprintf("release %s not found", params.ReleaseID))
	}

	bz, err := types.ModuleCdc.MarshalJSON(release)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

/*func queryProfileByAddress(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryByAddressParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	profile, found := k.GetProfileByAddress(ctx, params.Address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrProfileNotFound, fmt.Sprintf("profile with address %s not found", params.Address))
	}

	bz, err := types.ModuleCdc.MarshalJSON(profile)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
*/
