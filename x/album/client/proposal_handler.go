package client

import (
	"github.com/bitsongofficial/go-bitsong/x/album/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/album/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
