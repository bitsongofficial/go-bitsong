package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"

	"github.com/BitSongOfficial/go-bitsong/x/track/types"
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
	FlagRedistributionSplitRate = "redistribution_split_rate"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	songTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Track transaction subcommands",
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
		Short:   "Publish a new track",
		Example: "$ bitsongcli tx track publish --title=SongTitle --content=ipfs_hash --redistribution_split_rate=0.05 --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get flags
			var splitRate sdk.Dec

			from := cliCtx.GetFromAddress()
			title := viper.GetString(FlagTitle)
			content := viper.GetString(FlagContent)
			redistributionSplitRate := viper.GetString(FlagRedistributionSplitRate)
			if redistributionSplitRate != "" {
				rate, err := sdk.NewDecFromStr(redistributionSplitRate)
				if err != nil {
					return fmt.Errorf("invalid new redistribution split rate: %v", err)
				}

				splitRate = rate
			}

			msg := types.NewMsgPublish(title, from, content, splitRate)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "track title, eg. SongTitle")
	cmd.Flags().String(FlagContent, "", "track content, eg. <ipfs_url>")
	cmd.Flags().String(FlagRedistributionSplitRate, "", "percentage redistribution_split_rate, eg. 0.05 for 5%")

	return cmd
}

func GetCmdPlay(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "play",
		Short:   "Play a track",
		Example: "$ bitsongcli tx track play --id=1 --from mykey",
		//Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			//cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Get listener address
			listener := cliCtx.GetFromAddress()
			trackID, err := strconv.ParseUint(viper.GetString(FlagID), 10, 64)
			if err != nil {
				return fmt.Errorf("id %s not a valid uint, please input a valid track id", args[0])
			}

			msg := types.NewMsgPlay(trackID, listener)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagID, "", "track id, eg. 1")

	return cmd
}
