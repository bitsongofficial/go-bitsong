package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strconv"
)

// GetQueryCmd returns the query commands for the nft module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the merkledrop module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryMerkledrop(),
		GetCmdGenerateProofs(),
	)

	return queryCmd
}

func GetCmdQueryMerkledrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "detail [id]",
		Long:    "Query a merkledrop detail by id.",
		Example: fmt.Sprintf(`$ %s query merkledrop detail [id]`, version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Merkledrop(context.Background(), &types.QueryMerkledropRequest{
				Id: uint64(id),
			})
			res.Merkledrop.String()

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func printTree(tree Tree) {
	for _, row := range tree {
		for _, item := range row {
			fmt.Printf("%x  ", item[26:])
		}
		fmt.Println()
	}
}

func GetCmdGenerateProofs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate-proofs [in-file] [out-file]",
		Long:    "Generate proofs with user input",
		Example: fmt.Sprintf(`$ %s query merkledrop generate-proofs [in-file] [out-file]`, version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := client.GetClientQueryContext(cmd)
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
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
