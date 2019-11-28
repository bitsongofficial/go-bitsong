package rest

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Create Artist (POST)
	r.HandleFunc(
		"/artist/create",
		postArtistHandlerFn(cliCtx),
	).Methods("POST")
}

// PostArtistReq defines the properties of an artist request's body.
type PostArtistReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Name    string         `json:"name" yaml:"name"`   // Name of the artist
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"` // Address of the owner
}

func postArtistHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostArtistReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		meta := types.MetaFromArtist(req.Name)

		msg := types.NewMsgCreateArtist(meta, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
