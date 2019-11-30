package client

import (
	"github.com/bitsongofficial/go-bitsong/x/artist/client/cli"
	"github.com/bitsongofficial/go-bitsong/x/artist/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
