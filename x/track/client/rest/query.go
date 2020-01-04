package rest

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/go-bitsong/x/track/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

const (
	RestTrackID     = "track-id"
	RestOwner       = "owner"
	RestTrackStatus = "status"
	RestNumLimit    = "limit"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/track/all", queryTracksWithParameterFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/track/{%s}", RestTrackID), queryTrackHandlerFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query albums with parameters
func queryTracksWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// HTTP Params
		bechOwnerAddr := r.URL.Query().Get(RestOwner)
		strTrackStatus := r.URL.Query().Get(RestTrackStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		// Param object
		params := types.QueryTracksParams{}

		// Check bech32 Address Address
		if len(bechOwnerAddr) != 0 {
			ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			}
			params.Owner = ownerAddr
		}

		// Check Track Status
		if len(strTrackStatus) != 0 {
			trackStatus, err := types.TrackStatusFromString(cli.NormalizeTrackStatus(strTrackStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.TrackStatus = trackStatus
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
		res, height, err := cliCtx.QueryWithData("custom/track/tracks", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Response
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryTrackHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strTrackID := vars[RestTrackID]

		if len(strTrackID) == 0 {
			err := errors.New("trackId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		trackID, ok := rest.ParseUint64OrReturnBadRequest(w, strTrackID)
		if !ok {
			return
		}

		cliCtx, ok = rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryTrackParams(trackID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/track/track", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
