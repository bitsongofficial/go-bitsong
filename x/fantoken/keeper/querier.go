package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryFanToken:
			return queryFanToken(ctx, req, k, legacyQuerierCdc)
		case types.QueryFanTokens:
			return queryFanTokens(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryTotalBurn:
			return queryTotalBurn(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown token query endpoint")
		}
	}
}

func queryFanToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryFanTokenParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}

	token, err := keeper.GetFanToken(ctx, params.Denom)
	if err != nil {
		return nil, err
	}

	return codec.MarshalJSONIndent(legacyQuerierCdc, token)
}

func queryFanTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryFanTokensParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, err
	}
	tokens := keeper.GetFanTokens(ctx, params.Owner)
	return codec.MarshalJSONIndent(legacyQuerierCdc, tokens)
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := keeper.GetParamSet(ctx)
	return codec.MarshalJSONIndent(legacyQuerierCdc, params)
}

func queryTotalBurn(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	resp, err := keeper.TotalBurn(sdk.WrapSDKContext(ctx), &types.QueryTotalBurnRequest{})
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(legacyQuerierCdc, resp)
}
