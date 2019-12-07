package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/reward/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group track queries under a subcommand
	trackQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the reward module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	trackQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryAllRewards(queryRoute, cdc),
	)...)

	return trackQueryCmd
}

func GetCmdQueryAllRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Query all rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all rewards.
Example:
$ %s query reward all
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/rewards", queryRoute))
			if err != nil {
				return err
			}

			var rewards types.Rewards
			cdc.MustUnmarshalJSON(res, &rewards)
			return cliCtx.PrintOutput(rewards)
		},
	}
}
