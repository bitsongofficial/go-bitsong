package cli

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

const (
	FlagName        = "name"
	FlagMetadataURI = "metadata-uri"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	artistTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Artist transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	artistTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateArtist(cdc),
		GetCmdDeposit(cdc),
	)...)

	return artistTxCmd
}

// GetCmdCreateArtist implements the create artist command handler.
func GetCmdCreateArtist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new artist initialized with status nil",
		Long: strings.TrimSpace(fmt.Sprintf(`Create a new Artist initialized with status nil.
Example:
$ %s tx artist create --name="Freddy Mercury" --metadata-uri="QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u" --from mykey
`,
			version.ClientName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get flags
			flagName := viper.GetString(FlagName)               // Get artist name
			flagMetadataURI := viper.GetString(FlagMetadataURI) // Get artist metadata uri

			// Get params
			from := cliCtx.GetFromAddress() // Get owner

			// Build create artist message
			msg := types.NewMsgCreateArtist(flagName, flagMetadataURI, from)

			// Run basic validation
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagName, "", "the artist name")
	cmd.Flags().String(FlagMetadataURI, "", "the artist metadata uri")

	return cmd
}

func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [artist-id] [deposit]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens for an unverified artist",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a deposit for an unverified artist. You can
find the artist-id by running "%s query artist all".
Example:
$ %s tx artist deposit 1 10ubtsg --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the artist id is a uint
			artistID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("artist-id %s not a valid uint, please input a valid artist-id", args[0])
			}

			// Get depositor address
			from := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(from, artistID, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
