package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
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
				--total-amount=100000000000`,
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

			totAmt, err := cmd.Flags().GetUint64(FlagTotalAmount)
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
