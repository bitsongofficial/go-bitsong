package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"strconv"
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
		GetCmdAddTrack(cdc),
		GetCmdSubmitProposal(cdc),
	)...)

	return albumTxCmd
}

// GetCmdCreateAlbum implements the create album command handler.
func GetCmdCreateAlbum(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Args:  cobra.ExactArgs(1),
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

// GetCmdSubmitProposal implements the command to submit an album verify proposal
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-album [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an album verify proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an album verify proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx album verify-album <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "My new alkbum",
  "description": "Please, verify my new album. BTSG Topic: https://btsg.community/......",
  "id":  1, 
  "deposit": [
    {
      "denom": "ubtsg",
      "amount": "10000"
    }
  ]
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseAlbumVerifyProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewAlbumVerifyProposal(proposal.Title, proposal.Description, proposal.AlbumID)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdAddTrack implements creating a new add track command.
func GetCmdAddTrack(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-track [album-id] [track-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Add a track on a specific album",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a track on a specific album. You can
find the album-id by running "%s query album alls".
Example:
$ %s tx album add-track 1 1 --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get voting address
			from := cliCtx.GetFromAddress()

			// validate that the album id is a uint
			albumID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("album-id %s not a valid int, please input a valid album-id", args[0])
			}

			// validate that the track id is a uint
			trackID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("track-id %s not a valid int, please input a valid track-id", args[1])
			}

			// Find out which vote option user chose
			/*byteVoteOption, err := types.VoteOptionFromString(govutils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}*/

			// Build add-track message and run basic validation
			msg := types.NewMsgAddTrackAlbum(albumID, trackID, from)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
