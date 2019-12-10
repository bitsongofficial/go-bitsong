package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

type (
	DistributorVerifyProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Address     sdk.AccAddress `json:"address" yaml:"address"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

func ProposalRESTHandler(cliCtx context.CLIContext) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "verify_distributor",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DistributorVerifyProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewDistributorVerifyProposal(req.Title, req.Description, req.Address)

		msg := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
