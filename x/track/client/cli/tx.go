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

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

const (
	FlagTitle       = "title"
	FlagDescription = "description"
	FlagAudio       = "audio"
	FlagImage       = "image"
	FlagDuration    = "duration"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	trackTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Track transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	trackTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateTrack(cdc),
		GetCmdPlay(cdc),
		GetCmdDeposit(cdc),
	)...)

	return trackTxCmd
}

// GetCmdCreateTrack implements the create track command handler.
func GetCmdCreateTrack(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new track initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Track initialized with status nil.
Example:
$ %s tx track create --title "The Show Must Go On" --description="The track description" --audio="QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u" --image="QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u" --duration 385 --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagTitle := viper.GetString(FlagTitle) // Get track title
			flagDescription := viper.GetString(FlagDescription)
			flagAudio := viper.GetString(FlagAudio)
			flagImage := viper.GetString(FlagImage)
			flagDuration := viper.GetString(FlagDuration)

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create track message
			msg := types.NewMsgCreateTrack(flagTitle, flagDescription, flagAudio, flagImage, flagDuration, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "the track title")
	cmd.Flags().String(FlagDescription, "", "the track description")
	cmd.Flags().String(FlagAudio, "", "the track audio")
	cmd.Flags().String(FlagImage, "", "the track image")
	cmd.Flags().String(FlagDuration, "", "the track duration")

	return cmd
}

// GetCmdPlay implements creating a new play track command.
func GetCmdPlay(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "play [track-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Add a play on a specific track",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a play on a specific track. You can
find the track-id by running "%s query track all".
Example:
$ %s tx track play 1 --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get acc address
			from := cliCtx.GetFromAddress()

			// validate that the track id is a uint
			trackID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("track-id %s not a valid int, please input a valid track-id", args[0])
			}

			// Build play message and run basic validation
			msg := types.NewMsgPlay(trackID, from)
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
		Use:   "deposit [track-id] [deposit]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens for an unverified track",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a deposit for an unverified track. You can
find the track-id by running "%s query track all".
Example:
$ %s tx track deposit 1 10ubtsg --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the track id is a uint
			trackID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("track-id %s not a valid uint, please input a valid track-id", args[0])
			}

			// Get depositor address
			from := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(from, trackID, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
