package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	flagPage  = "page"
	flagLimit = "limit"
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
			GetCmdAll(cdc),
			GetCmdID(cdc),
			GetCmdCreator(cdc),
		)...,
	)

	return contentQueryCmd
}

func GetCmdAll(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Get all tracks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get all tracks.
Example:
$ %s query track all
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)
			params := types.DefaultQueryTracksParams(page, limit)

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTracks)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var tracks []types.Track
			if err := cdc.UnmarshalJSON(res, &tracks); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tracks)
		},
	}

	cmd.Flags().Int(flagLimit, 100, "pagination limit of tracks to query for")
	cmd.Flags().Int(flagPage, 1, "pagination page of tracks to to query for")

	return cmd
}

func GetCmdID(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "id [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a track by trackID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a track by trackID.
Example:
$ %s query track id [id]
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			id := args[0]

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryID, id)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not find trackID %s \n", id)
				return nil
			}

			var track types.Track
			cdc.MustUnmarshalJSON(res, &track)
			return cliCtx.PrintOutput(track)
		},
	}
}

func GetCmdCreator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "creator [creator-addres]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a track by creator address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a track by creator address.
Example:
$ %s query track creator bitsong1ccy7n32j9vsydn3y7qh2208zz0ap04rfp67ky9
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			creator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryCreatorTracksParams{Creator: creator})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCreatorTracks)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("Could not find creator %s \n", creator)
				return nil
			}

			var tracks []types.Track
			if err := cdc.UnmarshalJSON(res, &tracks); err != nil {
				return err
			}
			return cliCtx.PrintOutput(tracks)
		},
	}
}
