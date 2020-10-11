package keeper

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryArtist:
			return queryArtist(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown artist query endpoint")
		}
	}
}

func queryArtist(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryArtistParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	artist, found := k.GetArtist(ctx, btsg.ID(params.ID))
	if !found {
		return nil, sdkerrors.Wrap(types.ErrArtistNotFound, fmt.Sprintf("artist %s not found", params.ID))
	}

	bz, err := types.ModuleCdc.MarshalJSON(artist)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
