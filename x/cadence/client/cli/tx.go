package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
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
		NewRegisterCadenceContract(),
		NewUnregisterCadenceContract(),
		NewUnjailCadenceContract(),
	)
	return txCmd
}

// NewRegisterCadenceContract returns a CLI command handler for registering a
// contract for the cadence module.
func NewRegisterCadenceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [contract_bech32]",
		Short: "Register a cadence contract .",
		Long:  "Register a cadence contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgRegisterCadenceContract{
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

// NewUnregisterCadenceContract returns a CLI command handler for unregistering a
// contract for the cadence module.
func NewUnregisterCadenceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister [contract_bech32]",
		Short: "Unregister a cadence contract .",
		Long:  "Unregister a cadence contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnregisterCadenceContract{
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

// NewUnjailCadenceContract returns a CLI command handler for unjailing a
// contract for the cadence module.
func NewUnjailCadenceContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail [contract_bech32]",
		Short: "Unjail a cadence contract .",
		Long:  "Unjail a cadence contract . Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()
			contractAddress := args[0]

			msg := &types.MsgUnjailCadenceContract{
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
