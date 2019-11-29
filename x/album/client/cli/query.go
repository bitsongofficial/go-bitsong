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

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

const (
	FlagOwner    = "owner"
	FlagStatus   = "status"
	flagNumLimit = "limit"
)

// GetQueryCmd returns the cli query commands for the album module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group album queries under a subcommand
	albumQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the album module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	albumQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryAlbums(queryRoute, cdc),
		// TODO: create GetCmdQueryAlbum
	)...)

	return albumQueryCmd
}

// GetCmdQueryAlbums implements a query albums command.
func GetCmdQueryAlbums(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Query all albums with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all albums. You can filter the returns with the following flags.
Example:
$ %s query album all --owner cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query album all --status (Verified|Rejected|Failed)
$ %s query album all --limit 10
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechOwnerAddr := viper.GetString(FlagOwner)
			strAlbumStatus := viper.GetString(FlagStatus)
			numLimit := uint64(viper.GetInt64(flagNumLimit))

			var ownerAddr sdk.AccAddress
			var albumStatus types.AlbumStatus

			params := types.NewQueryAlbumsParams(ownerAddr, albumStatus, numLimit)

			if len(bechOwnerAddr) != 0 {
				ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
				if err != nil {
					return err
				}
				params.Owner = ownerAddr
			}

			if len(strAlbumStatus) != 0 {
				albumStatus, err := types.AlbumStatusFromString(NormalizeAlbumStatus(strAlbumStatus))
				if err != nil {
					return err
				}
				params.AlbumStatus = albumStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/albums", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingAlbums types.Albums
			err = cdc.UnmarshalJSON(res, &matchingAlbums)
			if err != nil {
				return err
			}

			if len(matchingAlbums) == 0 {
				return fmt.Errorf("No matching albums found")
			}

			return cliCtx.PrintOutput(matchingAlbums) // nolint:errcheck
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] albums. Defaults to all albums")
	cmd.Flags().String(FlagOwner, "", "(optional) filter by albums owned by address")
	cmd.Flags().String(FlagStatus, "", "(optional) filter albums by album status, status: verified/failed/rejected")

	return cmd
}

//NormalizeAlbumStatus - normalize user specified album status
func NormalizeAlbumStatus(status string) string {
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
