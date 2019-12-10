package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group track queries under a subcommand
	trackQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the distributor module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	trackQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryAllDistributors(queryRoute, cdc),
	)...)

	return trackQueryCmd
}

func GetCmdQueryAllDistributors(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Query all distributors",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all distributors.
Example:
$ %s query distributor all
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/distributors", queryRoute))
			if err != nil {
				return err
			}

			var distributors types.Distributors
			cdc.MustUnmarshalJSON(res, &distributors)
			return cliCtx.PrintOutput(distributors)
		},
	}
}
