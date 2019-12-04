package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov"
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
	FlagTitle = "title"
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
		GetCmdSubmitProposal(cdc),
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
$ %s tx track create --title "The Show Must Go On" --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagTitle := viper.GetString(FlagTitle) // Get track title

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create track message
			msg := types.NewMsgCreateTrack(flagTitle, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "the track title")

	return cmd
}

// GetCmdSubmitProposal implements the command to submit a track verify proposal
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-track [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a track verify proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a track verify proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx track verify-track <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "The Show Must Go On",
  "description": "Please, verify my track. BTSG Topic: https://btsg.community/......",
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

			proposal, err := ParseTrackVerifyProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewTrackVerifyProposal(proposal.Title, proposal.Description, proposal.TrackID)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

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
