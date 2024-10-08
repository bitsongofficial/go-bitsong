package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

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
		GetCmdMint(),
		GetCmdBurn(),
		GetCmdDisableMint(),
		GetCmdSetAuthority(),
		GetCmdSetMinter(),
		GetCmdSetUri(),
		// GetCmdUpdateFantokenFees(),
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
				Minter:    authority.String(),
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
		Use:   "mint [amount][denom]",
		Short: "Mint fan tokens to a specified address.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken mint [amount][denom] "+
				"--recipient=<recipient>"+
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

			minter := clientCtx.GetFromAddress().String()

			rcpt, err := cmd.Flags().GetString(FlagRecipient)
			if err != nil {
				return err
			}
			if rcpt != "" {
				if _, err = sdk.AccAddressFromBech32(rcpt); err != nil {
					return err
				}
			}

			coin, err := sdk.ParseCoinNormalized(strings.TrimSpace(args[0]))
			if err != nil {
				return err
			}

			msg := fantokentypes.NewMsgMint(rcpt, coin, minter)

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
		Use:   "burn [amount][denom]",
		Short: "Burn fantoken.",
		Example: fmt.Sprintf(
			"$ %s tx fantoken burn [amount][denom] "+
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

			coin, err := sdk.ParseCoinNormalized(strings.TrimSpace(args[0]))
			if err != nil {
				return err
			}

			msg := fantokentypes.NewMsgBurn(coin, owner)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSetAuthority implements the transfer fan token authority command
func GetCmdSetAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-authority [denom]",
		Short: "Set the authority of a fan token to a new authority",
		Example: fmt.Sprintf(
			"$ %s tx fantoken set-authority <denom> "+
				"--new-authority=<new-authority> "+
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

			oldAuthority := clientCtx.GetFromAddress().String()

			newAuthority, _ := cmd.Flags().GetString(FlagNewAuthority)

			if len(strings.TrimSpace(newAuthority)) > 0 {
				if _, err := sdk.AccAddressFromBech32(newAuthority); err != nil {
					return err
				}
			}

			msg := fantokentypes.NewMsgSetAuthority(strings.TrimSpace(args[0]), oldAuthority, newAuthority)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSetAuthority)
	//_ = cmd.MarkFlagRequired(FlagNewAuthority)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSetMinter implements the transfer fan token authority command
func GetCmdSetMinter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-minter [denom]",
		Short: "Set the minter of a fan token to a new minter",
		Example: fmt.Sprintf(
			"$ %s tx fantoken set-minter <denom> "+
				"--new-minter=<new-minter> "+
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

			oldMinter := clientCtx.GetFromAddress().String()

			newMinter, err := cmd.Flags().GetString(FlagNewMinter)
			if err != nil {
				return err
			}
			if _, err := sdk.AccAddressFromBech32(newMinter); err != nil {
				return err
			}

			msg := fantokentypes.NewMsgSetMinter(strings.TrimSpace(args[0]), oldMinter, newMinter)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSetMinter)
	_ = cmd.MarkFlagRequired(FlagNewMinter)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSetUri implements the transfer fan token authority command
func GetCmdSetUri() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-uri [denom]",
		Short: "Set the uri of the fantoken",
		Example: fmt.Sprintf(
			"$ %s tx fantoken set-uri <denom> "+
				"--uri=<uri> "+
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

			uri, err := cmd.Flags().GetString(FlagURI)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress().String()

			msg := fantokentypes.NewMsgSetUri(strings.TrimSpace(args[0]), uri, authority)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsSetUri)
	_ = cmd.MarkFlagRequired(FlagURI)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdUpdateFantokenFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fantoken-fees [proposal-file]",
		Short: "Submit an update fantoken fees proposal.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update fantoken fees proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal update-fantoken-fees <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "Update Fantoken Fees Proposal",
  "description": "update the current fees",
  "issue_fee": "1000000ubtsg",
  "mint_fee": "1000000ubtsg",
  "burn_fee": "1000000ubtsg",
  "deposit": "500000000ubtsg"
}
`, version.AppName,
			),
		),
		Example: fmt.Sprintf(
			"$ %s tx gov submit-proposal update-fantoken-fees [proposal-file] "+
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

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			content := fantokentypes.NewUpdateFeesProposal(proposal.Title, proposal.Description, issueFee, mintFee, burnFee)

			msg, err := v1beta1.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

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
