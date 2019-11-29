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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

const (
	FlagName   = "name"
	FlagHeight = "imageHeight"
	FlagWidth  = "imageWidth"
	FlagCid    = "cid"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	artistTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Artist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	artistTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateArtist(cdc),
		GetCmdSetArtistImage(cdc),
	)...)

	return artistTxCmd
}

// GetCmdCreateArtist implements the create artist command handler.
func GetCmdCreateArtist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-artist",
		Short: "create new artist initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Artist initialized with status nil.
Example:
$ %s tx artist create-artist --name="Freddy Mercury" --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			name := viper.GetString(FlagName)

			accGetter := authtypes.NewAccountRetriever(cliCtx)
			from := cliCtx.GetFromAddress()
			if err := accGetter.EnsureExists(from); err != nil {
				return err
			}

			msg := types.NewMsgCreateArtist(name, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagName, "", "the artist name")

	return cmd
}

// GetCmdSetArtistImage command to set a new artist image
func GetCmdSetArtistImage(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-image [artist-id] --imageHeight [height] --imageWidth [width] --cid [cid]",
		Args:  cobra.ExactArgs(1),
		Short: "Set a new image to artist",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set a new image to artist. You can
find the artist-id by running "%s query artist artists --owner [your-key]".
Example:
$ %s tx artist set-image 1 --imageHeight 500 --imageWidth 500 --cid QM..... --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagHeight := viper.GetString(FlagHeight)
			flagWidth := viper.GetString(FlagWidth)
			flagCid := viper.GetString(FlagCid)

			// Get owner address
			owner := cliCtx.GetFromAddress()

			// validate that the artist id is a uint
			artistID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("artist-id %s not a valid int, please input a valid artist-id", args[0])
			}

			// validate that the image height is a uint
			height, err := strconv.ParseUint(flagHeight, 10, 64)
			if err != nil {
				return fmt.Errorf("height %s not a valid int, please input a valid height", flagHeight)
			}

			// validate that the image width is a uint
			width, err := strconv.ParseUint(flagWidth, 10, 64)
			if err != nil {
				return fmt.Errorf("width %s not a valid int, please input a valid width", flagWidth)
			}

			// TODO:
			// validate cid ???

			// Build set artist image message and run basic validation
			msg := types.NewMsgSetArtistImage(artistID, height, width, flagCid, owner)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagHeight, "", "the image height")
	cmd.Flags().String(FlagWidth, "", "the image width")
	cmd.Flags().String(FlagCid, "", "the image cid")

	return cmd
}
