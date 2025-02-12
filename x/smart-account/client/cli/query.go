package cli

import (
	"github.com/bitsongofficial/go-bitsong/btsgutils/btsgcli"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	// "github.com/bitsongofficial/go-bitsong/btsgutils/btsgcli"
)

func GetQueryCmd() *cobra.Command {
	cmd := btsgcli.QueryIndexCmd(types.ModuleName)
	btsgcli.AddQueryCmd(cmd, types.NewQueryClient, GetCmdAuthenticators)
	btsgcli.AddQueryCmd(cmd, types.NewQueryClient, GetCmdAuthenticator)
	btsgcli.AddQueryCmd(cmd, types.NewQueryClient, GetCmdParams)

	return cmd
}

func GetCmdAuthenticators() (*btsgcli.QueryDescriptor, *types.GetAuthenticatorsRequest) {
	return &btsgcli.QueryDescriptor{
		Use:   "authenticators",
		Short: "Query authenticators by account",
		Long: `{{.Short}}{{.ExampleHeader}}
{{.CommandPrefix}} bitsong12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj`,
	}, &types.GetAuthenticatorsRequest{}
}

func GetCmdAuthenticator() (*btsgcli.QueryDescriptor, *types.GetAuthenticatorRequest) {
	return &btsgcli.QueryDescriptor{
		Use:   "authenticator",
		Short: "Query authenticator by account and id",
		Long: `{{.Short}}{{.ExampleHeader}}
{{.CommandPrefix}} bitsong12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj 17`,
	}, &types.GetAuthenticatorRequest{}
}

func GetCmdParams() (*btsgcli.QueryDescriptor, *types.QueryParamsRequest) {
	return &btsgcli.QueryDescriptor{
		Use:   "params",
		Short: "Query smartaccount params",
		Long: `{{.Short}}{{.ExampleHeader}}
{{.CommandPrefix}} params`,
	}, &types.QueryParamsRequest{}
}
