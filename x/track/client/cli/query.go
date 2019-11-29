package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

const (
	FlagOwner    = "owner"
	FlagStatus   = "status"
	flagNumLimit = "limit"
)

// GetQueryCmd returns the cli query commands for the track module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group track queries under a subcommand
	trackQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the track module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	trackQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryTracks(queryRoute, cdc),
		// TODO: create GetCmdQueryTrack
	)...)

	return trackQueryCmd
}

// GetCmdQueryTracks implements a query tracks command.
func GetCmdQueryTracks(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Query all tracks with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all tracks. You can filter the returns with the following flags.
Example:
$ %s query track all --owner cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query track all --status (Verified|Rejected|Failed)
$ %s query track all --limit 10
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechOwnerAddr := viper.GetString(FlagOwner)
			strTrackStatus := viper.GetString(FlagStatus)
			numLimit := uint64(viper.GetInt64(flagNumLimit))

			var ownerAddr sdk.AccAddress
			var trackStatus types.TrackStatus

			params := types.NewQueryTracksParams(ownerAddr, trackStatus, numLimit)

			if len(bechOwnerAddr) != 0 {
				ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
				if err != nil {
					return err
				}
				params.Owner = ownerAddr
			}

			if len(strTrackStatus) != 0 {
				trackStatus, err := types.TrackStatusFromString(NormalizeTrackStatus(strTrackStatus))
				if err != nil {
					return err
				}
				params.TrackStatus = trackStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tracks", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingTracks types.Tracks
			err = cdc.UnmarshalJSON(res, &matchingTracks)
			if err != nil {
				return err
			}

			if len(matchingTracks) == 0 {
				return fmt.Errorf("No matching tracks found")
			}

			return cliCtx.PrintOutput(matchingTracks) // nolint:errcheck
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] tracks. Defaults to all tracks")
	cmd.Flags().String(FlagOwner, "", "(optional) filter by tracks owned by address")
	cmd.Flags().String(FlagStatus, "", "(optional) filter tracks by track status, status: verified/failed/rejected")

	return cmd
}

//NormalizeTrackStatus - normalize user specified track status
func NormalizeTrackStatus(status string) string {
	switch status {
	case "Verified", "verified":
		return "Verified"
	case "Rejected", "rejected":
		return "Rejected"
	case "Failed", "failed":
		return "Failed"
	}
	return ""
}
