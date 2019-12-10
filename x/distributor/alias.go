package distributor

import (
	"github.com/bitsongofficial/go-bitsong/x/distributor/client"
	"github.com/bitsongofficial/go-bitsong/x/distributor/keeper"
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	// Keeper methods
	NewKeeper  = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	// Codec
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec

	// Types
	NewDistributor                      = types.NewDistributor
	NewDistributorVerifyProposal        = types.NewDistributorVerifyProposal
	NewDistributorVerifyProposalHandler = keeper.NewDistributorVerifyProposalHandler
	ProposalHandler                     = client.ProposalHandler

	// Msgs
	NewMsgCreateDistributor = types.NewMsgCreateDistributor
)

type (
	// Keeper
	Keeper = keeper.Keeper

	// Types
	Distributor  = types.Distributor
	Distributors = types.Distributors

	// Msgs
	MsgCreateDistributor = types.MsgCreateDistributor
)
