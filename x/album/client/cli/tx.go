package cli

import (
	"fmt"
	"github.com/spf13/viper"
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
	FlagTitle       = "title"
	FlagAlbumType   = "album-type"
	FlagMetadataURI = "metadata-uri"
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
		GetCmdDeposit(cdc),
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
$ %s tx album create --title="Innuendo" --album-type="Single" --metadata-uri="ipfs:QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u" --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagTitle := viper.GetString(FlagTitle)             // Get album title
			flagMetadataURI := viper.GetString(FlagMetadataURI) // Get album metadata uri
			flagAlbumType := viper.GetString(FlagAlbumType)

			albumType, err := types.AlbumTypeFromString(flagAlbumType)
			if err != nil {
				return err
			}

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create artist message
			msg := types.NewMsgCreateAlbum(albumType, flagTitle, flagMetadataURI, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "the album title")
	cmd.Flags().String(FlagMetadataURI, "", "the album metadata uri")
	cmd.Flags().String(FlagAlbumType, "", "the album type")

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

func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [album-id] [deposit]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens for an unverified album",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a deposit for an unverified album. You can
find the album-id by running "%s query album all".
Example:
$ %s tx album deposit 1 10ubtsg --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the album id is a uint
			albumID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("artist-id %s not a valid uint, please input a valid artist-id", args[0])
			}

			// Get depositor address
			from := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(from, albumID, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
