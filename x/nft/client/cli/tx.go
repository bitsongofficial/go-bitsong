package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

// NewTxCmd returns the transaction commands for the nft module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nft transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateNFT(),
		GetCmdTransferNFT(),
		GetCmdSignMetadata(),
		GetCmdUpdateMetadata(),
		GetCmdUpdateMetadataAuthority(),
		GetCmdCreateCollection(),
		GetCmdVerifyCollection(),
		GetCmdUnverifyCollection(),
		GetCmdUpdateCollectionAuthority(),
	)

	return txCmd
}

func GetCmdCreateNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-nft",
		Long: "Create an nft from provided params",
		Example: fmt.Sprintf(
			`$ %s tx nft create-nft
				--name="Punk10"
				--symbol="PUNK"
				--uri="https://punk.com/10"
				--seller-fee-basis-points=100
				--creators="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47"
				--creator-shares="10"
				--mutable=false
				--update-authority="bitsong13m350fvnk3s6y5n8ugxhmka277r0t7cw48ru47"`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			updateAuthority, err := cmd.Flags().GetString(FlagUpdateAuthority)
			if err != nil {
				return err
			}
			symbol, err := cmd.Flags().GetString(FlagSymbol)
			if err != nil {
				return err
			}
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			uri, err := cmd.Flags().GetString(FlagUri)
			if err != nil {
				return err
			}
			sellerFeeBasisPoints, err := cmd.Flags().GetUint32(FlagSellerFeeBasisPoints)
			if err != nil {
				return err
			}
			isMutable, err := cmd.Flags().GetBool(FlagMutable)
			if err != nil {
				return err
			}

			creatorAccsStr, err := cmd.Flags().GetString(FlagCreators)
			if err != nil {
				return err
			}
			creatorAccs := []string{}
			if creatorAccsStr != "" {
				creatorAccs = strings.Split(creatorAccsStr, ",")
			}
			creatorSharesStr, err := cmd.Flags().GetString(FlagCreatorShares)
			if err != nil {
				return err
			}
			creatorShareStrs := []string{}
			if creatorSharesStr != "" {
				creatorShareStrs = strings.Split(creatorSharesStr, ",")
			}
			creators := []*types.Creator{}
			for index, creatorAcc := range creatorAccs {
				shareStr := creatorShareStrs[index]
				share, err := strconv.Atoi(shareStr)
				if err != nil {
					return err
				}
				creators = append(creators, &types.Creator{
					Address: creatorAcc,
					Share:   uint32(share),
				})
			}

			msg := types.NewMsgCreateNFT(clientCtx.GetFromAddress(), updateAuthority, types.Data{
				Name:                 name,
				Symbol:               symbol,
				Uri:                  uri,
				SellerFeeBasisPoints: sellerFeeBasisPoints,
				Creators:             creators,
			}, false, isMutable)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateNFT())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdTransferNFT() *cobra.Command {
	// MsgTransferNFT
}

func GetCmdSignMetadata() *cobra.Command {
	// MsgSignMetadata
}
func GetCmdUpdateMetadata() *cobra.Command {
	// MsgUpdateMetadata

}

func GetCmdUpdateMetadataAuthority() *cobra.Command {
	// MsgUpdateMetadataAuthority
}

func GetCmdCreateCollection() *cobra.Command {
	// MsgCreateCollection
}

func GetCmdVerifyCollection() *cobra.Command {
	// MsgVerifyCollection
}

func GetCmdUnverifyCollection() *cobra.Command {
	// MsgUnverifyCollection
}

func GetCmdUpdateCollectionAuthority() *cobra.Command {
	// MsgUpdateCollectionAuthority
}
