package rest

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/artist/{id}", queryArtistFn(cliCtx)).Methods("GET")
}

func queryArtistFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		params := types.NewQueryArtistParams(id)
		bz := cliCtx.Codec.MustMarshalJSON(params)

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryArtist)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
