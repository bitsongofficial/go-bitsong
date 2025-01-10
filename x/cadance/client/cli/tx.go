package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/Clock
// transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Clock subcommands.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterCadanceContract(),
		NewUnregisterCadanceContract(),
		NewUnjailCadanceContract(),
	)
	return txCmd
}

// NewRegisterCadanceContract returns a CLI command handler for registering a
// contract for the cadance module.
func NewRegisterCadanceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [contract_bech32]",
		Short: "Register a cadance contract .",
		Long:  "Register a cadance contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgRegisterCadanceContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewUnregisterCadanceContract returns a CLI command handler for unregistering a
// contract for the cadance module.
func NewUnregisterCadanceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister [contract_bech32]",
		Short: "Unregister a cadance contract .",
		Long:  "Unregister a cadance contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnregisterCadanceContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewUnjailCadanceContract returns a CLI command handler for unjailing a
// contract for the cadance module.
func NewUnjailCadanceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail [contract_bech32]",
		Short: "Unjail a cadance contract .",
		Long:  "Unjail a cadance contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnjailCadanceContract{
				SenderAddress:   senderAddress.String(),
				ContractAddress: contractAddress,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
