package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/bitsongofficial/go-bitsong/x/distributor/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/distributor/client/rest"
)

var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
