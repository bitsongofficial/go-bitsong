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
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

const (
	flagIpfsAddr      = "ipfs-addr"
	flagStreamPrice   = "stream-price"
	flagDownloadPrice = "download-price"
	flagRightHolder   = "right-holder"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	contentTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	contentTxCmd.AddCommand(flags.PostCommands(
		GetCmdAdd(cdc),
		GetCmdStream(cdc),
		GetCmdDownload(cdc),
	)...)

	return contentTxCmd
}

func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new content to bitsong",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a new content inside bitsong.
Example:
$ %s tx content add [uri] [name] [meta-file] [content-file] \
--stream-price=[streamPrice] \
--download-price=[downloadPrice] \
--right-holder "80:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a" \
--right-holder "20:bitsong13r9ryyfltaz8rsqqumqxusgtw0ne4udhxm5jm4" \
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			ipfsAddr := viper.GetString(flagIpfsAddr)
			sh := shell.NewShell(ipfsAddr)

			metadata, err := os.Open(args[2])
			if err != nil {
				return err
			}
			metadataCID, err := sh.Add(metadata)
			if err != nil {
				return err
			}

			content, err := os.Open(args[3])
			if err != nil {
				return err
			}
			contentCID, err := sh.Add(content)
			if err != nil {
				return err
			}

			uri, name := args[0], args[1]
			streamPrice := viper.GetString(flagStreamPrice)
			downloadPrice := viper.GetString(flagDownloadPrice)

			rhsStr, err := cmd.Flags().GetStringArray(flagRightHolder)
			if err != nil {
				return fmt.Errorf("invalid rights holders value")
			}

			rhs := types.RightsHolders{}
			for _, rh := range rhsStr {
				rhArgs := strings.Split(rh, ":")
				if len(rhArgs) != 2 {
					return fmt.Errorf("the right holder format must be \"quota:address\" ex: \"100:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a\"")
				}
				rhq, err := sdk.NewDecFromStr(rhArgs[0])
				if err != nil {
					return err
				}
				rhAddr, err := sdk.AccAddressFromBech32(rhArgs[1])
				if err != nil {
					return fmt.Errorf("right holder address is wrong, %s", err.Error())
				}
				rh := types.NewRightHolder(rhq, rhAddr)
				rhs = append(rhs, rh)
			}

			msg := types.NewMsgAddContent(name, uri, "/ipfs/"+metadataCID, "/ipfs/"+contentCID, streamPrice, downloadPrice, rhs)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagStreamPrice, "", "Stream Price of the content")
	cmd.Flags().String(flagDownloadPrice, "", "Download Price of the content")
	cmd.Flags().StringArray(flagRightHolder, []string{}, "Rights Holders of the content")
	cmd.Flags().String(flagIpfsAddr, "http://localhost:5001", "IPFS address node")

	return cmd
}

func GetCmdStream(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Stream a content",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Stream a content.
Example:
$ %s tx content stream [uri]
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			uri := args[0]

			msg := types.NewMsgStream(uri, cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdDownload(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a new content",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Download a new content.
Example:
$ %s tx content download [uri]
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			uri := args[0]

			msg := types.NewMsgDownload(uri, cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
