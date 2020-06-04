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
	"io/ioutil"
	"strings"

	hlsUitls "github.com/bitsongofficial/go-bitsong/x/content/client/utils"
)

const (
	flagDao = "dao"
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
		GetCmdStore(cdc),
		GetCmdAction(cdc),
	)...)

	return contentTxCmd
}

func GetCmdStore(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "Store a new m3u8",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			hls, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			// gzip the hls file
			if hlsUitls.IsHls(hls) {
				hls, err = hlsUitls.GzipIt(hls)
				if err != nil {
					return err
				}
			} else if !hlsUitls.IsGzip(hls) {
				return fmt.Errorf("invalid hls file. please use hls or gzip")
			}

			msg := types.MsgStoreHls{
				From:        cliCtx.GetFromAddress(),
				HLSByteCode: hls,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new content to bitsong",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add a new content inside bitsong.
Example:
$ %s tx content add [uri] [hash] \
--dao "2000:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a" \
--dao "1000:bitsong13r9ryyfltaz8rsqqumqxusgtw0ne4udhxm5jm4" \
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			uri, hash := args[0], args[1]

			daoStr, err := cmd.Flags().GetStringArray(flagDao)
			if err != nil {
				return fmt.Errorf("invalid dao value")
			}

			dao := types.Dao{}
			for _, rh := range daoStr {
				deArgs := strings.Split(rh, ":")
				if len(deArgs) != 2 {
					return fmt.Errorf("the dao format must be \"shares:address\" ex: \"1000:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a\"")
				}
				des, err := sdk.NewDecFromStr(deArgs[0])
				if err != nil {
					return err
				}
				deAddr, err := sdk.AccAddressFromBech32(deArgs[1])
				if err != nil {
					return fmt.Errorf("dao entity address is wrong, %s", err.Error())
				}
				de := types.NewDaoEntity(des, deAddr)
				dao = append(dao, de)
			}

			msg := types.NewMsgAddContent(uri, hash, dao)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringArray(flagDao, []string{}, "DAO of the content")

	return cmd
}

func GetCmdAction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "action",
		Short: "Action to content",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Action to content.
Example:
$ %s tx content action [uri]
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			msg := types.NewMsgContentAction(args[0], cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
