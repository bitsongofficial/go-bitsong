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
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

const (
	flagEntities      = "entity"
	flagProvider      = "provider"
	flagStreamPrice   = "stream-price"
	flagDownloadPrice = "download-price"
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
		GetCmdCreate(cdc),
		GetCmdAddShare(cdc),
		GetCmdRemoveShare(cdc),
	)...)

	return contentTxCmd
}

func GetCmdCreate(cdc *codec.Codec) *cobra.Command {
	// TODO: remove id

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Smart Media Contract on bitsong",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new Smart Media Contract on bitsong.
Example:
$ %s tx track create [trackID] \
[contract.json] \
--entity=100:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a \
--entity=200:bitsong1dykf46zf3ss442j6cydaajk27xalny9y9chwnz \
--from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			trackID := args[0]
			fileInfo, err := ioutil.ReadFile(args[1])
			if err != nil {
				return err
			}

			trackInfoBz := new(bytes.Buffer)
			if err := json.Compact(trackInfoBz, fileInfo); err != nil {
				return err
			}

			entitiesStr, err := cmd.Flags().GetStringArray(flagEntities)
			if err != nil {
				return fmt.Errorf("invalid entities value")
			}

			var entities []types.Entity
			for _, ent := range entitiesStr {
				eargs := strings.Split(ent, ":")
				if len(eargs) != 2 {
					return fmt.Errorf("the entities format must be \"shares:address\" ex: \"1000:bitsong1xe8z84hcvgavtrtqv9al9lk2u3x5gysu44j54a\"")
				}
				eShares, ok := sdk.NewIntFromString(eargs[0])
				if !ok {
					return fmt.Errorf("invalid entities shares: %s", eargs[0])
				}
				eAddr, err := sdk.AccAddressFromBech32(eargs[1])
				if err != nil {
					return fmt.Errorf("entity address is wrong, %s", err.Error())
				}

				entity := types.Entity{
					Shares:  eShares,
					Address: eAddr,
				}
				entities = append(entities, entity)
			}

			providerArg, err := cmd.Flags().GetString(flagProvider)
			if err != nil {
				return err
			}
			provider, err := sdk.AccAddressFromBech32(providerArg)
			if err != nil {
				return err
			}

			streamPriceArg, err := cmd.Flags().GetString(flagStreamPrice)
			if err != nil {
				return err
			}
			streamPrice, err := sdk.ParseCoin(streamPriceArg)
			if err != nil {
				return err
			}

			downloadPriceArg, err := cmd.Flags().GetString(flagDownloadPrice)
			if err != nil {
				return err
			}
			downloadPrice, err := sdk.ParseCoin(downloadPriceArg)
			if err != nil {
				return err
			}

			//msg := types.NewMsgTrackCreate(trackID, trackInfoBz.Bytes(), cliCtx.FromAddress, provider, entities, streamPrice, downloadPrice)
			msg := types.NewMsgTrackCreate(trackID, cliCtx.FromAddress, provider, entities, streamPrice, downloadPrice)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagProvider, "", "Unique Provider (optional)")
	cmd.Flags().String(flagStreamPrice, "0btsg", "Stream Price (optional)")
	cmd.Flags().String(flagDownloadPrice, "0btsg", "Stream Price (optional)")
	cmd.Flags().StringArray(flagEntities, []string{""}, "Track Entities")

	return cmd
}

func GetCmdAddShare(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-share",
		Short: "Add share to Smart Media Contract",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add share to Smart Media Contract.
Example:
$ %s tx track add-share [trackID] [shares] --from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			amt, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTrackAddShare(args[0], amt, cliCtx.FromAddress)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetCmdRemoveShare(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-share",
		Short: "Remove share to Smart Media Contract",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Remove share from Smart Media Contract.
Example:
$ %s tx track remove-share [trackID] [shares] --from <creator>`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			amt, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTrackRemoveShare(args[0], amt, cliCtx.FromAddress)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
