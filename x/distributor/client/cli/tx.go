package cli

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	FlagName = "name"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	albumTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Distributor transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	albumTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateDistributor(cdc),
		GetCmdSubmitProposal(cdc),
	)...)

	return albumTxCmd
}

func GetCmdCreateDistributor(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create --name [name]",
		Short: "create new distributor",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Distributor.
The distributor details must be supplied via a JSON file.
Example:
$ %s tx album create --name [distributor-name] --from=<key_or_address>
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagName := viper.GetString(FlagName) // Get name

			// Get address
			from := cliCtx.GetFromAddress()

			// Build create track message
			msg := types.NewMsgCreateDistributor(flagName, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagName, "", "the distributor name")

	return cmd
}

func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-verify [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a distributor verify proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a distributor verify proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx distributor submit-verify <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "BitSong Distributor",
  "description": "Please, verify us. BTSG Topic: https://btsg.community/......",
  "address":  "", 
  "deposit": [
    {
      "denom": "ubtsg",
      "amount": "10000"
    }
  ]
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseDistributorVerifyProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewDistributorVerifyProposal(proposal.Title, proposal.Description, proposal.Address)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
