package rest

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/release/{releaseID}", queryReleaseHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/release/creator/{creator}", queryAllReleaseForCreatorHandlerFn(cliCtx)).Methods("GET")
}

func queryReleaseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		releaseID := vars["releaseID"]

		params := types.NewQueryReleaseParams(releaseID)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRelease)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryAllReleaseForCreatorHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		creator, err := sdk.AccAddressFromBech32(vars["creator"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		params := types.NewQueryAllReleaseForCreatorParams(creator)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllReleaseForCreator)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
