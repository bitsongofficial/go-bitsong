package client

import (
	"github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdUpdateFantokenFees, ProposalRESTHandler)

func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{}
}
