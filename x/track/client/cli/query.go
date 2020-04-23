package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group track queries under a subcommand
	trackQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: client.ValidateCmd,
	}

	trackQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmqQueryTrack(cdc),
		)...,
	)

	return trackQueryCmd
}

func GetCmqQueryTrack(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "track-addr [track-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single track",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single track.
Example:
$ %s query track track-addr B0FA2953B126722264F67828AF7443144C85D867
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the track id is a uint
			trackAddr := args[0]

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTrack, trackAddr)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not find track with addr %s \n", trackAddr)
				return nil
			}

			var track types.Track
			cdc.MustUnmarshalJSON(res, &track)
			return cliCtx.PrintOutput(track)
		},
	}
}
