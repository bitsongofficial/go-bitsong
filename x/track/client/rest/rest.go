package rest

import (
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// RegisterRoutes registers track-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type CreateTrackReq struct {
	BaseReq       rest.BaseReq        `json:"base_req"`
	Title         string              `json:"title"`
	Attributes    types.Attributes    `json:"attributes"`
	Media         types.TrackMedia    `json:"media"`
	Rewards       types.TrackRewards  `json:"rewards"`
	RightsHolders types.RightsHolders `json:"rights_holders"`
}
