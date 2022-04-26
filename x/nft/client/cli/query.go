package cli

import (
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

// GetQueryCmd returns the query commands for the nft module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the nft module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand()

	return queryCmd
}
