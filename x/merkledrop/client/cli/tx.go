package cli

import (
	"encoding/json"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// NewTxCmd returns the transaction commands for the merkledrop module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "merkledrop transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdGenerate(),
		GetCmdCreate(),
		GetCmdClaim(),
		GetCmdWithdraw(),
	)

	return txCmd
}

func GetCmdGenerate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [file-json] [out-list-json]",
		Short: "Generate a merkledrop from json file",
		Long: `Generate a merkledrop from json file
Parameters:
	file-json: input file list
	out-list-json: output list with proofs
		`,
		Example: fmt.Sprintf(`
$ %s tx merkledrop create-from-file accounts.json out-list.json

where accounts.json contains
{
	"bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj2": "1000000ubtsg",
	"bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw": "2000000ubtsg",
	"bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw": "3000000ubtsg"
}

after the computation the out-list.json should be similar to this output
{
  "bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw": {
    "index": 2,
    "coin": "3000000ubtsg",
    "proof": [
      "3346fbddeb1d097311651f5615d3b2528a3893fb79b2ce40b740e6d470296d85"
    ]
  },
  "bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj2": {
    "index": 0,
    "coin": "1000000ubtsg",
    "proof": [
      "a258c32bee9b0bbb7a2d1999ab4698294844e7440aa6dcd067e0d5142fa20522",
      "7f0b92cc8318e4fb0db9052325b474e2eabb80d79e6e1abab92093d3a88fe029"
    ]
  },
  "bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw": {
    "index": 1,
    "coin": "2000000ubtsg",
    "proof": [
      "7a807e653a5d63556f46fd66a2ac9af6bddaa6864611e6b8da2ccf8389a91345",
      "7f0b92cc8318e4fb0db9052325b474e2eabb80d79e6e1abab92093d3a88fe029"
    ]
  }
}
`,
			version.AppName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			listBytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			var stringList map[string]string
			if err := json.Unmarshal(listBytes, &stringList); err != nil {
				log.Fatalf("Could not unmarshal json: %v", err)
			}

			accMap, err := AccountsFromMap(stringList)
			if err != nil {
				log.Fatalf("Could not get accounts from map")
			}

			tree, claimInfo, err := CreateDistributionList(accMap)
			if err != nil {
				log.Fatalf("Could not create distribution list: %v", err)
			}

			if _, err := createFile(args[1], claimInfo); err != nil {
				log.Fatalf("Could not create file: %v", err)
			}

			fmt.Println(fmt.Sprintf("Merkle Root: %x", tree.Root()))
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		Long: "Create a merkledrop from provided params",
		Example: fmt.Sprintf(`
$ %s tx merkledrop create \
	--merkle-root="98ac4ade3eae2e324922ee68c42976eeaecc39d558fcfc2206ec3ab0bad5a36b" \
	--amount=100000000000 \
	--denom=ubtsg \
	--start-height=1 \
	--end-height=10
`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			merkleRoot, err := cmd.Flags().GetString(FlagMerkleRoot)
			if err != nil {
				return err
			}

			startHeight, err := cmd.Flags().GetInt64(FlagStartHeight)
			if err != nil {
				return err
			}

			endHeight, err := cmd.Flags().GetInt64(FlagEndHeight)
			if err != nil {
				return err
			}

			amount, err := cmd.Flags().GetInt64(FlagAmount)
			if err != nil {
				return err
			}

			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", amount, denom))
			if err != nil {
				return err
			}

			msg := types.NewMsgCreate(clientCtx.GetFromAddress(), merkleRoot, startHeight, endHeight, coin)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagCreateMerkledrop())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "claim",
		Long: "Claim a merkledrop from provided params",
		Args: cobra.ExactArgs(1),
		Example: fmt.Sprintf(`
$ %s tx merkledrop claim [id] \
	--proofs="20245fe3fcdbf17069bc0de04e319296766a7138be5e5a27c6f5bc05e0c23de9,b8fedba5a18186d4fb92ffcf9924b408d6048aaeb76b10cad97cf6be4071b710" \
	--amount=1000
`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			merkledropId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			proofsStr, err := cmd.Flags().GetString(FlagProofs)
			if err != nil {
				return err
			}
			proofs := []string{}
			if proofsStr != "" {
				proofs = strings.Split(proofsStr, ",")
			}

			amount, err := cmd.Flags().GetInt64(FlagAmount)
			if err != nil {
				return err
			}

			index, err := cmd.Flags().GetUint64(FlagIndex)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaim(index, merkledropId, sdk.NewInt(amount), proofs, clientCtx.GetFromAddress())

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagClaimMerkledrop())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "withdraw",
		Long: "Withdraw funds from an expired merkledrop",
		Args: cobra.ExactArgs(1),
		Example: fmt.Sprintf(`
$ %s tx merkledrop withdraw [id]
`,
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			merkledropId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(clientCtx.GetFromAddress(), merkledropId)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagClaimMerkledrop())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
