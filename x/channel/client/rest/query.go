package rest

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/channel/{handle}", queryChannelByHandleFn(cliCtx)).Methods("GET")
	r.HandleFunc("/channel/owner/{addr}", queryChannelByOwnerFn(cliCtx)).Methods("GET")
}

func queryChannelByHandleFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		handle := vars["handle"]

		params := types.NewQueryChannelParams(handle)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChannel)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryChannelByOwnerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		owner, _ := sdk.AccAddressFromBech32(vars["addr"])

		params := types.NewQueryByOwnerParams(owner)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryChannelByOwner)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
