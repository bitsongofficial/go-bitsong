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
