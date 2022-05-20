package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// NewTxCmd returns the transaction commands for the merkledrop module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "merkledrop transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateMerkledrop(),
		GetCmdClaimMerkledrop(),
	)

	return txCmd
}

func GetCmdCreateMerkledrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a merkledrop from provided params",
		Example: fmt.Sprintf(
			`$ %s tx merkledrop create
				--merkle-root="98ac4ade3eae2e324922ee68c42976eeaecc39d558fcfc2206ec3ab0bad5a36b"
				--total-amount=100000000000ubtsg`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			merkleRoot, err := cmd.Flags().GetString(FlagMerkleRoot)
			if err != nil {
				return err
			}

			totAmtStr, err := cmd.Flags().GetString(FlagTotalAmount)
			if err != nil {
				return err
			}

			totAmt, err := sdk.ParseCoinNormalized(totAmtStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateMerkledrop(clientCtx.GetFromAddress(), merkleRoot, totAmt)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateMerkledrop())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdClaimMerkledrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "claim",
		Long: "Claim a merkledrop from provided params",
		Args: cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`$ %s tx merkledrop claim [id]
				--proofs="20245fe3fcdbf17069bc0de04e319296766a7138be5e5a27c6f5bc05e0c23de9,b8fedba5a18186d4fb92ffcf9924b408d6048aaeb76b10cad97cf6be4071b710"
				--amount=1000ubtsg`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mdId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			proofsStr, err := cmd.Flags().GetString(FlagProofs)
			if err != nil {
				return err
			}
			proofs := []string{}
			if proofsStr != "" {
				proofs = strings.Split(proofsStr, ",")
			}

			amtStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}

			amt, err := sdk.ParseCoinNormalized(amtStr)
			if err != nil {
				return err
			}

			index, err := cmd.Flags().GetUint64(FlagIndex)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimMerkledrop(index, mdId, amt, proofs, clientCtx.GetFromAddress())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagClaimMerkledrop())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
