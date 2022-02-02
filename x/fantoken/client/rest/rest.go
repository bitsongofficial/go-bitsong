package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Rest variable names
// nolint
const (
	RestParamDenom = "denom"
	RestParamOwner = "owner"
)

// RegisterHandlers registers token-related REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type issueFanTokenReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Owner       string       `json:"owner"` // owner of the token
	Symbol      string       `json:"symbol"`
	Name        string       `json:"name"`
	MaxSupply   string       `json:"max_supply"`
	Mintable    bool         `json:"mintable"`
	Description string       `json:"description"`
	IssueFee    string       `json:"issue_fee"`
}

type editFanTokenReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Owner    string       `json:"owner"`    //  owner of the token
	Mintable bool         `json:"mintable"` // mintable of the token
}

type transferFanTokenOwnerReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	SrcOwner string       `json:"src_owner"` // the current owner address of the token
	DstOwner string       `json:"dst_owner"` // the new owner
}

type mintFanTokenReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`     // the current owner address of the token
	Recipient string       `json:"recipient"` // address of minting token to
	Amount    string       `json:"amount"`    // amount of minting token
}

type burnFanTokenReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Sender  string       `json:"owner"`  // the current owner address of the token
	Amount  string       `json:"amount"` // amount of burning token
}
