package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
)

// NewTxCmd returns the transaction commands for the nft module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        nfttypes.ModuleName,
		Short:                      "nft transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()

	return txCmd
}
