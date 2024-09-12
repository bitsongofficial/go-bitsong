package client

import (
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdUpdateMerkledropFees)
