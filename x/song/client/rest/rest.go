package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/BitSongOfficial/go-bitsong/x/song/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"
)

const (
	restTitle = "title"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/titles", storeName), titlesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/titles", storeName), setTitleHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/titles/{%s}", storeName, restTitle), resolveTitleHandler(cliCtx, storeName)).Methods("GET")
}

// --------------------------------------------------------------------------------------
// Tx Handler

type setNameReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Title   string       `json:"title"`
	Owner   string       `json:"owner"`
}

func setTitleHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setTitleReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgSetTitle(req.Title, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

//--------------------------------------------------------------------------------------
// Query Handlers

func resolveTitleHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restTitle]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func titlesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/titles", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}