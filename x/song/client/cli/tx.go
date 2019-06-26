package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/BitSongOfficial/go-bitsong/x/song/types"
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
		GetCmdSetName(cdc),
	)...)

	return songTxCmd
}

// GetCmdSetTitle is the CLI command for sending a SetTitle transaction
func GetCmdSetTitle(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-title [title]",
		Short: "set the title of your song",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := types.NewMsgSetTitle(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}