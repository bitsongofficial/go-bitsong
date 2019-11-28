package rest

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/client"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RestArtistID     = "artist-id"
	RestOwner        = "owner"
	RestArtistStatus = "status"
	RestNumLimit     = "limit"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Get all artists with params
	r.HandleFunc(
		"/artist/artists",
		queryArtistsWithParameterFn(cliCtx),
	).Methods("GET")

	// Get all artists from owner

	// Get all artists by status

	// Get artist by artistID
}

// HTTP request handler to query artists with parameters
func queryArtistsWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bechOwnerAddr := r.URL.Query().Get(RestOwner)
		strArtistStatus := r.URL.Query().Get(RestArtistStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

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
			artistStatus, err := types.ArtistStatusFromString(client.NormalizeArtistStatus(strArtistStatus))
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

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/artist/artists", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
