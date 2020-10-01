package rest

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/profile/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/profile/{handle}", queryProfileHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/profile/addr/{addr}", queryProfileHandlerFn(cliCtx)).Methods("GET")
}

func queryProfileHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr, _ := sdk.AccAddressFromBech32(vars["addr"])

		params := types.NewQueryByAddressParams(addr)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProfileByAddress)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
