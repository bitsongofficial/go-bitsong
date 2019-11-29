package rest

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/client/cli"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

const (
	RestOwner        = "owner"
	RestArtistStatus = "status"
	RestNumLimit     = "limit"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/artist/all", queryArtistsWithParameterFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query artists with parameters
func queryArtistsWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// HTTP Params
		bechOwnerAddr := r.URL.Query().Get(RestOwner)
		strArtistStatus := r.URL.Query().Get(RestArtistStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		// Param object
		params := types.QueryArtistsParams{}

		// Check bech32 Owner Address
		if len(bechOwnerAddr) != 0 {
			ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
			params.Owner = ownerAddr
		}

		// Check Artist Status
		if len(strArtistStatus) != 0 {
			artistStatus, err := types.ArtistStatusFromString(cli.NormalizeArtistStatus(strArtistStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ArtistStatus = artistStatus
		}

		// Check Limit Param
		if len(strNumLimit) != 0 {
			numLimit, ok := rest.ParseUint64OrReturnBadRequest(w, strNumLimit)
			if !ok {
				return
			}
			params.Limit = numLimit
		}

		// Parse state height
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Run query
		res, height, err := cliCtx.QueryWithData("custom/artist/artists", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
