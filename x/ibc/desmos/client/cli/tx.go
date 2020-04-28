package cli

import (
	"bufio"
	"strconv"
	"time"

	"github.com/bitsongofficial/go-bitsong/x/ibc/desmos/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcposts "github.com/desmos-labs/desmos/x/ibc/posts"
	"github.com/spf13/cobra"
)

// GetTransferTxCmd returns the command to create a NewMsgTransfer transaction
func GetIBCDesmosTxCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [src-port] [src-channel] [dest-height] [song-id] [post-owner]",
		Short: "Transfer fungible token through IBC",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			destHeight, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			// Create the post data
			data := types.NewSongCreationData(args[3], time.Now().UTC(), args[4])

			// Create and validate the message
			msg := ibcposts.NewMsgCrossPost(args[0], args[1], uint64(destHeight), data, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
