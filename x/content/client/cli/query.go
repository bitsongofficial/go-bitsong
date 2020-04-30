package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/content/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"
)

// GetQueryCmd returns the cli query commands
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group content queries under a subcommand
	contentQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	contentQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmqResolve(cdc),
		)...,
	)

	return contentQueryCmd
}

func GetCmqResolve(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [uri]",
		Args:  cobra.ExactArgs(1),
		Short: "Resolve an uri",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Resolve an uri inside bitsong.
Example:
$ %s query content resolve my-best-uri
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			uri := args[0]

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryUri, uri)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve uri %s \n", uri)
				return nil
			}

			var content types.Content
			cdc.MustUnmarshalJSON(res, &content)
			return cliCtx.PrintOutput(content)
		},
	}
}
