package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// issue a token
	r.HandleFunc(fmt.Sprintf("/%s/denom", types.ModuleName), issueTokenHandlerFn(cliCtx)).Methods("POST")
	// edit a token
	r.HandleFunc(fmt.Sprintf("/%s/denom/{%s}", types.ModuleName, RestParamDenom), editFanTokenHandlerFn(cliCtx)).Methods("PUT")
	// transfer owner
	r.HandleFunc(fmt.Sprintf("/%s/denom/{%s}/transfer", types.ModuleName, RestParamDenom), transferOwnerHandlerFn(cliCtx)).Methods("POST")
	// mint token
	r.HandleFunc(fmt.Sprintf("/%s/denom/{%s}/mint", types.ModuleName, RestParamDenom), mintFanTokenHandlerFn(cliCtx)).Methods("POST")
	// burn token
	r.HandleFunc(fmt.Sprintf("/%s/denom/{%s}/burn", types.ModuleName, RestParamDenom), burnFanTokenHandlerFn(cliCtx)).Methods("POST")
}

func issueTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueFanTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		maxSupply, ok := sdk.NewIntFromString(req.MaxSupply)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse max supply: %s", req.MaxSupply))
			return
		}

		issueFee, err := sdk.ParseCoinNormalized(req.IssueFee)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to parse issue fee: %s", req.IssueFee))
			return
		}

		// create the MsgIssueToken message
		msg := &types.MsgIssueFanToken{
			Symbol:      req.Symbol,
			Name:        req.Name,
			MaxSupply:   maxSupply,
			Description: req.Description,
			Owner:       req.Owner,
			IssueFee:    issueFee,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func editFanTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestParamDenom]

		var req editFanTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		mintable := req.Mintable

		// create the MsgEditToken message
		msg := types.NewMsgEditFanToken(denom, mintable, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func transferOwnerHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestParamDenom]

		var req transferFanTokenOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferFanTokenOwner(denom, req.SrcOwner, req.DstOwner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func mintFanTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestParamDenom]

		var req mintFanTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid amount %s", amount))
			return
		}

		// create the MsgMintFanToken message
		msg := types.NewMsgMintFanToken(req.Recipient, denom, req.Owner, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func burnFanTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[RestParamDenom]

		var req burnFanTokenReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid amount %s", amount))
			return
		}

		// create the MsgMintToken message
		msg := types.NewMsgBurnFanToken(denom, req.Sender, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
