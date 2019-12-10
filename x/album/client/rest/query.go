package rest

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/go-bitsong/x/album/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

const (
	RestOwner       = "owner"
	RestAlbumStatus = "status"
	RestNumLimit    = "limit"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/album/all", queryAlbumsWithParameterFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query albums with parameters
func queryAlbumsWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// HTTP Params
		bechOwnerAddr := r.URL.Query().Get(RestOwner)
		strAlbumStatus := r.URL.Query().Get(RestAlbumStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		// Param object
		params := types.QueryAlbumsParams{}

		// Check bech32 Address Address
		if len(bechOwnerAddr) != 0 {
			ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
			params.Owner = ownerAddr
		}

		// Check Album Status
		if len(strAlbumStatus) != 0 {
			albumStatus, err := types.AlbumStatusFromString(cli.NormalizeAlbumStatus(strAlbumStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.AlbumStatus = albumStatus
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
		res, height, err := cliCtx.QueryWithData("custom/album/albums", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
