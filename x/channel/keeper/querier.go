package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryChannel:
			return queryChannel(ctx, path[1:], req, k)
		case types.QueryChannelByOwner:
			return queryChannelByOwner(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown channel query endpoint")
		}
	}
}

func queryChannel(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryChannelParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	channel, found := k.GetChannel(ctx, params.Handle)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrChannelNotFound, fmt.Sprintf("channel %s not found", params.Handle))
	}

	bz, err := types.ModuleCdc.MarshalJSON(channel)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryChannelByOwner(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryByOwnerParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	channel, found := k.GetChannelByOwner(ctx, params.Owner)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrChannelNotFound, fmt.Sprintf("channel with address %s not found", params.Owner))
	}

	bz, err := types.ModuleCdc.MarshalJSON(channel)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
