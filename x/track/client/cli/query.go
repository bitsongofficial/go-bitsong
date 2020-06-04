package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
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
			GetCmqCid(cdc),
		)...,
	)

	return contentQueryCmd
}

func GetCmqCid(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cid [cid]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a track cid",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a cid inside bitsong.
Example:
$ %s query track cid [cid]
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cid := args[0]

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryCid, cid)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not resolve cid %s \n", cid)
				return nil
			}

			var track types.Track
			cdc.MustUnmarshalJSON(res, &track)
			return cliCtx.PrintOutput(track)
		},
	}
}
