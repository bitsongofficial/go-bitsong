package cli

import (
	"bufio"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/content/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	flagName       = "name"
	flagMetaUri    = "meta-uri"
	flagContentUri = "content-uri"
	flagDenom      = "denom"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	contentTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE: client.ValidateCmd,
	}

	contentTxCmd.AddCommand(flags.PostCommands(
		GetCmdAdd(cdc),
	)...)

	return contentTxCmd
}

func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new content",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a new content inside bitsong.
Example:
$ %s tx content add [uri] --name=[name] --meta-uri=[meta-uri] --content-uri=[content-uri] --denom=[denom]
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			name := viper.GetString(flagName)
			uri := args[0]
			metaUri := viper.GetString(flagMetaUri)
			contentUri := viper.GetString(flagContentUri)
			denom := viper.GetString(flagDenom)

			msg := types.NewMsgAddContent(name, uri, metaUri, contentUri, denom, cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagName, "", "Name of the content")
	cmd.Flags().String(flagMetaUri, "", "Meta Uri of the content")
	cmd.Flags().String(flagContentUri, "", "Content Uri of the content")
	cmd.Flags().String(flagDenom, "", "Denom of the content")

	return cmd
}
