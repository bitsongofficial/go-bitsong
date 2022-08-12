package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
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
		GetCmdUpdateCandyMachine(),
		GetCmdCloseCandyMachine(),
		GetCmdMintNFT(),
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

			collId, err := cmd.Flags().GetUint64(nftcli.FlagCollectionId)
			if err != nil {
				return err
			}
			price, err := cmd.Flags().GetUint64(FlagPrice)
			if err != nil {
				return err
			}

			treasury, err := cmd.Flags().GetString(FlagTreasury)
			if err != nil {
				return err
			}
			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return err
			}
			endSettingsTypeStr, err := cmd.Flags().GetString(FlagEndSettingsType)
			if err != nil {
				return err
			}
			endSettingsValue, err := cmd.Flags().GetUint64(FlagEndSettingsValue)
			if err != nil {
				return err
			}

			goLiveDate, err := cmd.Flags().GetUint64(FlagGoLiveDate)
			if err != nil {
				return err
			}

			metadataBaseUrl, err := cmd.Flags().GetString(FlagMetadataBaseUrl)
			if err != nil {
				return err
			}

			mutable, err := cmd.Flags().GetBool(nftcli.FlagMutable)
			if err != nil {
				return err
			}
			sellerFeeBasisPoints, err := cmd.Flags().GetUint32(nftcli.FlagSellerFeeBasisPoints)
			if err != nil {
				return err
			}

			creators, err := nftcli.CollectCreatorsData(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCandyMachine(
				clientCtx.GetFromAddress(),
				types.CandyMachine{
					CollId:     collId,
					Price:      price,
					Treasury:   treasury,
					Denom:      denom,
					GoLiveDate: goLiveDate,
					EndSettings: types.EndSettings{
						EndType: types.EndSettingType(types.EndSettingType_value[endSettingsTypeStr]),
						Value:   endSettingsValue,
					},
					Authority:            clientCtx.GetFromAddress().String(),
					MetadataBaseUrl:      metadataBaseUrl,
					Mutable:              mutable,
					SellerFeeBasisPoints: sellerFeeBasisPoints,
					Creators:             creators,
				},
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

func GetCmdUpdateCandyMachine() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update-candymachine [flags]",
		Long: "Update candy machine from provided params",
		Example: fmt.Sprintf(
			`$ %s tx candymachine update-candymachine`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collId, err := cmd.Flags().GetUint64(nftcli.FlagCollectionId)
			if err != nil {
				return err
			}
			price, err := cmd.Flags().GetUint64(FlagPrice)
			if err != nil {
				return err
			}

			treasury, err := cmd.Flags().GetString(FlagTreasury)
			if err != nil {
				return err
			}
			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return err
			}
			endSettingsTypeStr, err := cmd.Flags().GetString(FlagEndSettingsType)
			if err != nil {
				return err
			}
			endSettingsValue, err := cmd.Flags().GetUint64(FlagEndSettingsValue)
			if err != nil {
				return err
			}

			goLiveDate, err := cmd.Flags().GetUint64(FlagGoLiveDate)
			if err != nil {
				return err
			}

			metadataBaseUrl, err := cmd.Flags().GetString(FlagMetadataBaseUrl)
			if err != nil {
				return err
			}

			mutable, err := cmd.Flags().GetBool(nftcli.FlagMutable)
			if err != nil {
				return err
			}
			sellerFeeBasisPoints, err := cmd.Flags().GetUint32(nftcli.FlagSellerFeeBasisPoints)
			if err != nil {
				return err
			}

			creators, err := nftcli.CollectCreatorsData(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCandyMachine(
				clientCtx.GetFromAddress(),
				types.CandyMachine{
					CollId:     collId,
					Price:      price,
					Treasury:   treasury,
					Denom:      denom,
					GoLiveDate: goLiveDate,
					EndSettings: types.EndSettings{
						EndType: types.EndSettingType(types.EndSettingType_value[endSettingsTypeStr]),
						Value:   endSettingsValue,
					},
					Authority:            clientCtx.GetFromAddress().String(),
					MetadataBaseUrl:      metadataBaseUrl,
					Mutable:              mutable,
					SellerFeeBasisPoints: sellerFeeBasisPoints,
					Creators:             creators,
				},
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

func GetCmdCloseCandyMachine() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "close-candymachine [collection_id] [flags]",
		Long: "Close candy machine from provided params",
		Args: cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`$ %s tx candymachine close-candymachine 1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgCloseCandyMachine(
				clientCtx.GetFromAddress(),
				uint64(collId),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint-nft [collection_id] [nft_name] [flags]",
		Long: "Mint nft from provided params",
		Args: cobra.ExactArgs(2),
		Example: fmt.Sprintf(
			`$ %s tx candymachine close-candymachine 1 punk1`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collId, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(
				clientCtx.GetFromAddress(),
				uint64(collId),
				args[1],
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
