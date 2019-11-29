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
