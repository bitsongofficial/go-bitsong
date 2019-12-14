package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

const (
	FlagOwner    = "owner"
	FlagStatus   = "status"
	flagNumLimit = "limit"
)

// GetQueryCmd returns the cli query commands for the artist module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group artist queries under a subcommand
	artistQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the artist module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	artistQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryArtists(queryRoute, cdc),
		GetCmdQueryDeposits(queryRoute, cdc),
		// TODO: create GetCmdQueryArtist
	)...)

	return artistQueryCmd
}

// GetCmdQueryArtists implements a query artists command.
func GetCmdQueryArtists(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Query all artists with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all artists. You can filter the returns with the following flags.
Example:
$ %s query artist all --owner cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query artist all --status (Verified|Rejected|Failed)
$ %s query artist all --limit 10
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechOwnerAddr := viper.GetString(FlagOwner)
			strArtistStatus := viper.GetString(FlagStatus)
			numLimit := uint64(viper.GetInt64(flagNumLimit))

			var ownerAddr sdk.AccAddress
			var artistStatus types.ArtistStatus

			params := types.NewQueryArtistsParams(ownerAddr, artistStatus, numLimit)

			if len(bechOwnerAddr) != 0 {
				ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
				if err != nil {
					return err
				}
				params.Owner = ownerAddr
			}

			if len(strArtistStatus) != 0 {
				artistStatus, err := types.ArtistStatusFromString(NormalizeArtistStatus(strArtistStatus))
				if err != nil {
					return err
				}
				params.ArtistStatus = artistStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/artists", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingArtists types.Artists
			err = cdc.UnmarshalJSON(res, &matchingArtists)
			if err != nil {
				return err
			}

			if len(matchingArtists) == 0 {
				return fmt.Errorf("No matching artists found")
			}

			return cliCtx.PrintOutput(matchingArtists) // nolint:errcheck
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] artists. Defaults to all artists")
	cmd.Flags().String(FlagOwner, "", "(optional) filter by artists owned by address")
	cmd.Flags().String(FlagStatus, "", "(optional) filter artists by artist status, status: verified/failed/rejected")

	return cmd
}

//NormalizeArtistStatus - normalize user specified artist status
func NormalizeArtistStatus(status string) string {
	switch status {
	case "DepositPeriod", "deposit-period":
		return "DepositPeriod"
	case "Verified", "verified":
		return "Verified"
	case "Rejected", "rejected":
		return "Rejected"
	case "Failed", "failed":
		return "Failed"
	}
	return ""
}

func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposits [artist-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a specific artist",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all deposits on a specific artist.
You can find the artist-id by running "%s query artist all".
Example:
$ %s query artist deposits 1
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the artist id is a uint
			artistID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("artist-id %s not a valid uint, please input a valid artist-id", args[0])
			}

			params := types.NewQueryArtistParams(artistID)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// check to see if the proposal is in the store
			res, err := QueryArtistByID(artistID, cliCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch artist with id %d: %s", artistID, err)
			}

			var artist types.Artist
			cdc.MustUnmarshalJSON(res, &artist)

			artistStatus := artist.Status
			if !(artistStatus == types.StatusDepositPeriod) {
				res, err = QueryDepositsByTxQuery(cliCtx, params)
			} else {
				res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", queryRoute), bz)
			}

			if err != nil {
				return err
			}

			var dep types.Deposits
			cdc.MustUnmarshalJSON(res, &dep)
			return cliCtx.PrintOutput(dep)
		},
	}
}
