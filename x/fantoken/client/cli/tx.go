package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// NewTxCmd returns the transaction commands for the fantoken module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        tokentypes.ModuleName,
		Short:                      "Fan Token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssue(),
		GetCmdDisableMint(),
		GetCmdMint(),
		GetCmdBurn(),
		GetCmdTransferAuthority(),
	)

	return txCmd
}

// GetCmdIssue implements the issue fantoken command
func GetCmdIssue() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue",
		Long: "Issue a new fantoken.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken issue "+
				"--name=\"Kitty Token\" "+
				"--symbol=\"kitty\" "+
				"--max-supply=\"1000000000000\" "+
				"--uri=\"ipfs://...\" "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			maxSupplyStr, err := cmd.Flags().GetString(FlagMaxSupply)
			if err != nil {
				return err
			}
			maxSupply, ok := sdk.NewIntFromString(maxSupplyStr)
			if !ok {
				return fmt.Errorf("failed to parse max supply: %s", maxSupplyStr)
			}
			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return fmt.Errorf("the uri field is invalid")
			}

			msg := &tokentypes.MsgIssue{
				Symbol:    symbol,
				Name:      name,
				MaxSupply: maxSupply,
				Authority: authority.String(),
				URI:       uri,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsIssue)
	_ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagMaxSupply)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdDisableMint implements the edit fan token mintable command
func GetCmdDisableMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "disable-mint [denom]",
		Long: "Disable Mint of an existing fantoken.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken disable-mint <denom> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress().String()

			msg := tokentypes.NewMsgDisableMint(args[0], authority)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsDisableMint)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint [denom]",
		Long: "Mint fan tokens to a specified address.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken mint <denom> "+
				"--recipient=<recipient>"+
				"--amount=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress().String()

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return fmt.Errorf("invalid amount %s", amount)
			}

			addr, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}
			if addr != "" {
				if _, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			msg := tokentypes.NewMsgMint(
				addr, strings.TrimSpace(args[0]), authority, amount,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMint)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [denom]",
		Long: "Burn fantoken.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken burn <denom> "+
				"--amount=<amount> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress().String()

			amountStr, err := cmd.Flags().GetString(FlagAmount)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return fmt.Errorf("invalid amount %s", amount)
			}

			addr, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}
			if len(addr) > 0 {
				if _, err = sdk.AccAddressFromBech32(addr); err != nil {
					return err
				}
			}

			denom := strings.TrimSpace(args[0])

			msg := tokentypes.NewMsgBurn(denom, owner, amount)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMint)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferAuthority implements the transfer fan token authority command
func GetCmdTransferAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "transfer [denom]",
		Long: "Transfer the authority of a fan token to a new authority.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken transfer <denom> "+
				"--recipient=<recipient> "+
				"--from=<key-name> "+
				"--chain-id=<chain-id> "+
				"--fees=<fee>",
			version.AppName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			srcAuthority := clientCtx.GetFromAddress().String()

			dstAuthority, err := cmd.Flags().GetString(FlagDstAuthority)
			if err != nil {
				return err
			}
			if _, err := sdk.AccAddressFromBech32(dstAuthority); err != nil {
				return err
			}

			msg := tokentypes.NewMsgTransferAuthority(strings.TrimSpace(args[0]), srcAuthority, dstAuthority)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsTransferAuthority)
	_ = cmd.MarkFlagRequired(FlagDstAuthority)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
