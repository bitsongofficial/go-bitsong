package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router) {
	// issue a token
	r.HandleFunc(fmt.Sprintf("/%s/tokens", types.ModuleName), issueTokenHandlerFn(cliCtx)).Methods("POST")
	// edit a token
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}", types.ModuleName, RestParamSymbol), updateFanTokenHandlerFn(cliCtx)).Methods("PUT")
	// transfer owner
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}/transfer", types.ModuleName, RestParamSymbol), transferOwnerHandlerFn(cliCtx)).Methods("POST")
	// mint token
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}/mint", types.ModuleName, RestParamSymbol), mintFanTokenHandlerFn(cliCtx)).Methods("POST")
	// burn token
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}/burn", types.ModuleName, RestParamSymbol), burnFanTokenHandlerFn(cliCtx)).Methods("POST")
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

		// create the MsgIssueToken message
		msg := &types.MsgIssueFanToken{
			Symbol:    req.Symbol,
			Name:      req.Name,
			MaxSupply: maxSupply,
			Mintable:  req.Mintable,
			Owner:     req.Owner,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}

func updateFanTokenHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[RestParamSymbol]

		var req updateFanTokenMintableReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		mintable := req.Mintable

		// create the MsgEditToken message
		msg := types.NewMsgUpdateFanTokenMintable(symbol, mintable, req.Owner)
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
		symbol := vars[RestParamSymbol]

		var req transferFanTokenOwnerReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the MsgTransferTokenOwner message
		msg := types.NewMsgTransferFanTokenOwner(symbol, req.SrcOwner, req.DstOwner)
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
		symbol := vars[RestParamSymbol]

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
		msg := types.NewMsgMintFanToken(req.Recipient, symbol, req.Owner, amount)
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
		symbol := vars[RestParamSymbol]

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
		msg := types.NewMsgBurnFanToken(symbol, req.Sender, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
