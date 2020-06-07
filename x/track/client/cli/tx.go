package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
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
	"strconv"
	"strings"
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
		GetCmdPublish(cdc),
		GetCmdTokenize(cdc),
		GetCmdMint(cdc),
	)...)

	return contentTxCmd
}

func GetCmdPublish(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish a new track to bitsong",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Publish a new track to bitsong.
Example:
$ %s tx track publish [track-info.json] --from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			fileInfo, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			trackInfoBz := new(bytes.Buffer)
			if err := json.Compact(trackInfoBz, fileInfo); err != nil {
				return err
			}

			msg := types.NewMsgTrackPublish(trackInfoBz.Bytes(), cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

			/*daoStr, err := cmd.Flags().GetStringArray(flagDao)
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
			}*/
		},
	}

	// cmd.Flags().StringArray(flagDao, []string{}, "Track DAO")

	return cmd
}
func GetCmdTokenize(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokenize",
		Short: "Tokenize a new track",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Tokenize a new track.
Example:
$ %s tx track tokenize [track-id] [denom] --from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			trackID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			denom := strings.TrimSpace(args[1])

			msg := types.NewMsgTrackTokenize(trackID, denom, cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
func GetCmdMint(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "Mint track token to a recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint track token to a recipient.
Example:
$ %s tx track mint [trackID] [amount] [recipient] --from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			trackID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgTokenMint(trackID, amount, recipient, cliCtx.FromAddress)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// func GetCmdDisableMint
