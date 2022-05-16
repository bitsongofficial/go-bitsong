package cli

import (
	"github.com/spf13/cobra"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
)

// GetQueryCmd returns the query commands for the marketplace module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the marketplace module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand()

	return queryCmd
}
