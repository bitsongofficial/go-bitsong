package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
)

// NewTxCmd returns the transaction commands for the candymachine module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "candymachine transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateCandyMachine(),
	)

	return txCmd
}

func GetCmdCreateCandyMachine() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-candymachine [flags]",
		Long: "Create candy machine from provided params",
		Example: fmt.Sprintf(
			`$ %s tx candymachine create-candymachine`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCandyMachine(
				clientCtx.GetFromAddress(),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateCandyMachine())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
