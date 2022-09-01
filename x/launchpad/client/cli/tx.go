package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
	nftcli "github.com/bitsongofficial/go-bitsong/x/nft/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
)

// NewTxCmd returns the transaction commands for the launchpad module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "launchpad transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateLaunchPad(),
		GetCmdUpdateLaunchPad(),
		GetCmdCloseLaunchPad(),
		GetCmdMintNFT(),
		GetCmdMintNFTs(),
	)

	return txCmd
}

func GetCmdCreateLaunchPad() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-launchpad [flags]",
		Long: "Create launchpad from provided params",
		Example: fmt.Sprintf(
			`$ %s tx launchpad create-launchpad`,
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
			endTimestamp, err := cmd.Flags().GetUint64(FlagEndTimestamp)
			if err != nil {
				return err
			}
			maxMint, err := cmd.Flags().GetUint64(FlagMaxMint)
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

			shuffle, err := cmd.Flags().GetBool(FlagShuffle)
			if err != nil {
				return err
			}

			creators, err := nftcli.CollectCreatorsData(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateLaunchPad(
				clientCtx.GetFromAddress(),
				types.LaunchPad{
					CollId:               collId,
					Price:                price,
					Treasury:             treasury,
					Denom:                denom,
					GoLiveDate:           goLiveDate,
					EndTimestamp:         endTimestamp,
					MaxMint:              maxMint,
					Authority:            clientCtx.GetFromAddress().String(),
					MetadataBaseUrl:      metadataBaseUrl,
					Mutable:              mutable,
					SellerFeeBasisPoints: sellerFeeBasisPoints,
					Creators:             creators,
					Shuffle:              shuffle,
				},
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateLaunchPad())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdUpdateLaunchPad() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update-launchpad [flags]",
		Long: "Update launchpad from provided params",
		Example: fmt.Sprintf(
			`$ %s tx launchpad update-launchpad`,
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

			endTimestamp, err := cmd.Flags().GetUint64(FlagEndTimestamp)
			if err != nil {
				return err
			}
			maxMint, err := cmd.Flags().GetUint64(FlagMaxMint)
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

			shuffle, err := cmd.Flags().GetBool(FlagShuffle)
			if err != nil {
				return err
			}

			creators, err := nftcli.CollectCreatorsData(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateLaunchPad(
				clientCtx.GetFromAddress(),
				types.LaunchPad{
					CollId:               collId,
					Price:                price,
					Treasury:             treasury,
					Denom:                denom,
					GoLiveDate:           goLiveDate,
					EndTimestamp:         endTimestamp,
					MaxMint:              maxMint,
					Authority:            clientCtx.GetFromAddress().String(),
					MetadataBaseUrl:      metadataBaseUrl,
					Mutable:              mutable,
					SellerFeeBasisPoints: sellerFeeBasisPoints,
					Creators:             creators,
					Shuffle:              shuffle,
				},
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateLaunchPad())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdCloseLaunchPad() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "close-launchpad [collection_id] [flags]",
		Long: "Close launchpad from provided params",
		Args: cobra.ExactArgs(1),
		Example: fmt.Sprintf(
			`$ %s tx launchpad close-launchpad 1`,
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

			msg := types.NewMsgCloseLaunchPad(
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
			`$ %s tx launchpad close-launchpad 1 punk1`,
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

func GetCmdMintNFTs() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "mint-nfts [collection_id] [nft_count] [flags]",
		Long: "Mint nft from provided params",
		Args: cobra.ExactArgs(2),
		Example: fmt.Sprintf(
			`$ %s tx launchpad close-launchpad 1 2`,
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

			nftCount, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFTs(
				clientCtx.GetFromAddress(),
				uint64(collId),
				uint64(nftCount),
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
