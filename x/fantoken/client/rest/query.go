package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// Query the token by symbol or denom
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}", types.ModuleName, RestParamDenom), queryFanTokenHandlerFn(cliCtx)).Methods("GET")
	// Query tokens by owner
	r.HandleFunc(fmt.Sprintf("/%s/tokens", types.ModuleName), queryFanTokensHandlerFn(cliCtx)).Methods("GET")
	// Query token params
	r.HandleFunc(fmt.Sprintf("/%s/params", types.ModuleName), queryFanTokenParamsHandlerFn(cliCtx)).Methods("GET")
	// Query the total amount of all burned tokens
	r.HandleFunc(fmt.Sprintf("/%s/total_burn", types.ModuleName), queryTotalBurnHandlerFn(cliCtx)).Methods("GET")
}

// queryTokenHandlerFn is the HTTP request handler to query token
func queryFanTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		params := types.QueryFanTokenParams{
			Denom: vars[RestParamDenom],
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFanToken), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryTokensHandlerFn is the HTTP request handler to query tokens
func queryFanTokensHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := r.FormValue(RestParamOwner)

		var err error
		var owner sdk.AccAddress

		if len(ownerStr) > 0 {
			owner, err = sdk.AccAddressFromBech32(ownerStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		params := types.QueryFanTokensParams{
			Owner: owner,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFanTokens), bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryTokenParamsHandlerFn is the HTTP request handler to query token params
func queryFanTokenParamsHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams), nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// queryTotalBurnHandlerFn is the HTTP request handler to query the total amount of all burn tokens
func queryTotalBurnHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTotalBurn), nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
