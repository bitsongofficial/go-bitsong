package cli

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

const (
	FlagTitle                = "title"
	FlagAlbumType            = "album_type"
	FlagReleaseDate          = "release_date"
	FlagReleaseDatePrecision = "release_date_precision"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	albumTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Album transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	albumTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateAlbum(cdc),
	)...)

	return albumTxCmd
}

// GetCmdCreateAlbum implements the create album command handler.
func GetCmdCreateAlbum(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new album initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Album initialized with status nil.
Example:
$ %s tx album create --title "Innuendo" --album_type "Album" --release_date "2018-12-12" --release_date_precision "day" --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagTitle := viper.GetString(FlagTitle)                                                           // Get album title
			flagAlbumType, _ := types.AlbumTypeFromString(NormalizeAlbumType(viper.GetString(FlagAlbumType))) // Get album type
			flagReleaseDate := viper.GetString(FlagReleaseDate)                                               // Get album release date
			flagReleaseDatePrecision := viper.GetString(FlagReleaseDatePrecision)                             // Get album release date precision

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create artist message
			msg := types.NewMsgCreateAlbum(flagAlbumType, flagTitle, flagReleaseDate, flagReleaseDatePrecision, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "the album title")
	cmd.Flags().String(FlagAlbumType, "", "the album type")
	cmd.Flags().String(FlagReleaseDate, "", "the album release date")
	cmd.Flags().String(FlagReleaseDatePrecision, "", "the album release date precision")

	return cmd
}

//NormalizeAlbumType - normalize user specified album type
func NormalizeAlbumType(albumType string) string {
	switch albumType {
	case "Album", "album":
		return "Album"
	case "Single", "single":
		return "Single"
	case "Compilation", "compilation":
		return "Compilation"
	}
	return ""
}
