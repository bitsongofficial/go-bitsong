package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/BitSongOfficial/go-bitsong/x/song/types"
)

const (
	FlagTitle = "title"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	songTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Song transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	songTxCmd.AddCommand(client.PostCommands(
		GetCmdPublish(cdc),
	)...)

	return songTxCmd
}

// GetCmdPublish is the CLI command for sending a Publish transaction
func GetCmdPublish(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish a new song",
		Example: "$ bitsongcli song publish --title=SongTitle --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get from address
			from := cliCtx.GetFromAddress()

			// Pull associated account
			/*account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}*/

			title := viper.GetString(FlagTitle)

			msg := types.NewMsgPublish(title, from)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// FIX
			//cliCtx.PrintResponse = true

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "song title, eg. SongTitle")

	return cmd
}