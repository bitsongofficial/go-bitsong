package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

type CreateReleaseReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	ReleaseID   string       `json:"release_id"`
	MetadataURI string       `json:"metadata_uri"`
}

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}
