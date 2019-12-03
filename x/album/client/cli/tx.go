package cli

import (
	"fmt"
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
	return &cobra.Command{
		Use:   "create",
		Short: "create new album initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Album initialized with status nil.
The album details must be supplied via a JSON file.
Example:
$ %s tx album create <path/to/album.json> --from=<key_or_address>
Where album.json contains:
{
  "title": "Innuendo",
  "album_type": "Album",
  "release_date": "2018-12-12",
  "release_date_precision": "day"
}
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			album, err := ParseCreateAlbumJSON(cdc, args[0])
			if err != nil {
				return err
			}

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create artist message
			msg := types.NewMsgCreateAlbum(album.AlbumType, album.Title, album.ReleaseDate, album.ReleaseDatePrecision, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
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
