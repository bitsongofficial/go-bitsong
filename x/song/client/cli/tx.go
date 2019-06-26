package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/BitSongOfficial/go-bitsong/x/song/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

const (
	FlagID                      = "id"
	FlagTitle                   = "title"
	FlagContent                 = "content"
	FlagRedistributionSplitRate = "redistribuition_split_rate"
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
		GetCmdPlay(cdc),
	)...)

	return songTxCmd
}

// GetCmdPublish is the CLI command for sending a Publish transaction
func GetCmdPublish(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "publish",
		Short:   "Publish a new song",
		Example: "$ bitsongcli song publish --title=SongTitle --content=<ipfs_url> --redistribution_split_rate=5 --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get flags
			from := cliCtx.GetFromAddress()
			title := viper.GetString(FlagTitle)
			content := viper.GetString(FlagContent)
			redistributionSplitRate := viper.GetString(FlagRedistributionSplitRate)

			msg := types.NewMsgPublish(title, from, content, redistributionSplitRate)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// To fix or delete
			//return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "song title, eg. SongTitle")
	cmd.Flags().String(FlagContent, "", "song content, eg. <ipfs_url>")
	cmd.Flags().String(FlagRedistributionSplitRate, "", "song redistribution_split_rate, eg. 5")

	return cmd
}

func GetCmdPlay(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "play",
		Short:   "Play a song",
		Example: "$ bitsongcli song play --id=1 --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get listener address
			listener := cliCtx.GetFromAddress()
			id := viper.GetString(FlagID)

			msg := types.NewMsgPlay(id, listener)
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

	cmd.Flags().String(FlagID, "", "song id, eg. 1")

	return cmd
}
