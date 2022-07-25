package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
)

// GetQueryCmd returns the query commands for the candymachine module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the candymachine module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryCandyMachines(),
	)

	return queryCmd
}

func GetCmdQueryCandyMachines() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "candymachines [flags]",
		Long: "Query candymachines.",
		Example: fmt.Sprintf(
			`$ %s query candymachine candymachines`, version.AppName),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)

			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CandyMachines(context.Background(), &types.QueryCandyMachinesRequest{})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
