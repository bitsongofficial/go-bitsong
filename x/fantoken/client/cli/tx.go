package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

// NewTxCmd returns the transaction commands for the fantoken module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        fantokentypes.ModuleName,
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
		GetCmdUpdateFeesProposal(),
	)

	return txCmd
}

// GetCmdIssue implements the issue fantoken command
func GetCmdIssue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new fantoken.",
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

			msg := &fantokentypes.MsgIssue{
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
		Use:   "disable-mint [denom]",
		Short: "Disable Mint of an existing fantoken.",
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

			msg := fantokentypes.NewMsgDisableMint(args[0], authority)
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
		Use:   "mint [denom]",
		Short: "Mint fan tokens to a specified address.",
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

			msg := fantokentypes.NewMsgMint(
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
		Use:   "burn [denom]",
		Short: "Burn fantoken.",
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

			msg := fantokentypes.NewMsgBurn(denom, owner, amount)

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
		Use:   "transfer-authority [denom]",
		Short: "Transfer the authority of a fan token to a new authority",
		Example: fmt.Sprintf(
			"$ %s tx fantoken transfer-authority <denom> "+
				"--dst-authority=<dst-authority> "+
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

			msg := fantokentypes.NewMsgTransferAuthority(strings.TrimSpace(args[0]), srcAuthority, dstAuthority)

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

func GetCmdUpdateFeesProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fees-proposal [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update fees proposal.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update fees proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal update-fees-proposal <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Update Fantoken Fees Proposal",
  "description": "update the current fees",
  "issue_fee": "1000000ubtsg",
  "mint_fee": "1000000ubtsg",
  "burn_fee": "1000000ubtsg",
  "transfer_fee": "1000000ubtsg",
  "deposit": "500000000ubtsg"
}
`, version.AppName,
			),
		),
		Example: fmt.Sprintf(
			"$ %s tx fantoken update-fees-proposal [proposal-file] "+
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

			proposal, err := parseUpdateFeesProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			issueFee, err := sdk.ParseCoinNormalized(proposal.IssueFee)
			if err != nil {
				return err
			}

			mintFee, err := sdk.ParseCoinNormalized(proposal.MintFee)
			if err != nil {
				return err
			}

			burnFee, err := sdk.ParseCoinNormalized(proposal.BurnFee)
			if err != nil {
				return err
			}

			transferFee, err := sdk.ParseCoinNormalized(proposal.TransferFee)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := fantokentypes.NewUpdateFeesProposal(proposal.Title, proposal.Description, issueFee, mintFee, burnFee, transferFee)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseUpdateFeesProposal(cdc codec.JSONCodec, proposalFile string) (fantokentypes.UpdateFeesProposalWithDeposit, error) {
	proposal := fantokentypes.UpdateFeesProposalWithDeposit{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
