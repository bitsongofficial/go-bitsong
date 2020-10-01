package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/profile/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProfile:
			return queryProfile(ctx, path[1:], req, k)
		case types.QueryProfileByAddress:
			return queryProfileByAddress(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown profile query endpoint")
		}
	}
}

func queryProfile(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryProfileParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	profile, found := k.GetProfile(ctx, params.Handle)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrProfileNotFound, fmt.Sprintf("profile %s not found", params.Handle))
	}

	bz, err := types.ModuleCdc.MarshalJSON(profile)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryProfileByAddress(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
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
