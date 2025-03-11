package server

import (
	"io"
	"os"
	"path/filepath"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.WriteCloser, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}

// AddTestnetCreatorCommand allows chains to create a testnet from the state existing in their node's data directory.
func AddTestnetCreatorCommand(rootCmd *cobra.Command, appCreator types.AppCreator, addStartFlags types.ModuleInitFlags) {
	testnetCreateCmd := InPlaceTestnetCreator(appCreator)
	addStartFlags(testnetCreateCmd)
	rootCmd.AddCommand(testnetCreateCmd)
}
