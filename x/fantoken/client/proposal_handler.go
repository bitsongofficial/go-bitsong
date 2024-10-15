package client

import (
	"github.com/bitsongofficial/go-bitsong/v018/x/fantoken/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdUpdateFantokenFees)
