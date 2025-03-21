package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const (
	FlagHeight           = "height"
	FlagForZeroHeight    = "for-zero-height"
	FlagJailAllowedAddrs = "jail-allowed-addrs"
	FlagNewOperatorAddr  = "modules-to-export"
)

// ExportCmd dumps app state to JSON.
func CustomExportCmd(appExporter servertypes.AppExporter, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-export",
		Short: "Debug v022 upgradeHandler logic, via simulated export state.",
		Args:  cobra.NoArgs, // 		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// Set testnet keys to be used by the application.
			// This is done to prevent changes to existing start API.
			serverCtx.Viper.Set(server.KeyIsTestnet, true)

			homeDir, _ := cmd.Flags().GetString(flags.FlagHome)
			config.SetRoot(homeDir)

			if _, err := os.Stat(config.GenesisFile()); os.IsNotExist(err) {
				return err
			}

			db, err := openDB(config.RootDir, server.GetAppDBBackend(serverCtx.Viper))
			if err != nil {
				return err
			}

			if appExporter == nil {
				if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "WARNING: App exporter not defined. Returning genesis file."); err != nil {
					return err
				}

				// Open file in read-only mode so we can copy it to stdout.
				// It is possible that the genesis file is large,
				// so we don't need to read it all into memory
				// before we stream it out.
				f, err := os.OpenFile(config.GenesisFile(), os.O_RDONLY, 0)
				if err != nil {
					return err
				}
				defer f.Close()

				if _, err := io.Copy(cmd.OutOrStdout(), f); err != nil {
					return err
				}

				return nil
			}

			traceWriterFile, _ := cmd.Flags().GetString("trace-store")
			traceWriter, err := openTraceWriter(traceWriterFile)

			if err != nil {
				return err
			}

			// todo: add list of validators to retain power for testnet
			height, _ := cmd.Flags().GetInt64(FlagHeight)
			jailAllowedAddrs, _ := cmd.Flags().GetStringSlice(FlagJailAllowedAddrs)
			modulesToExport, _ := cmd.Flags().GetStringSlice(server.FlagModulesToExport)
			outputDocument, _ := cmd.Flags().GetString(flags.FlagOutputDocument)

			exported, err := appExporter(serverCtx.Logger, db, traceWriter, height, true, jailAllowedAddrs, serverCtx.Viper, modulesToExport)
			if err != nil {
				return fmt.Errorf("error exporting state: %w", err)
			}

			appGenesis, err := genutiltypes.AppGenesisFromFile(serverCtx.Config.GenesisFile())
			if err != nil {
				return err
			}

			// set current binary version
			appGenesis.AppName = version.AppName
			appGenesis.AppVersion = version.Version

			appGenesis.AppState = exported.AppState
			appGenesis.InitialHeight = exported.Height
			appGenesis.Consensus = genutiltypes.NewConsensusGenesis(exported.ConsensusParams, exported.Validators)

			out, err := json.Marshal(appGenesis)
			if err != nil {
				return err
			}

			if outputDocument == "" {
				// Copy the entire genesis file to stdout.
				_, err := io.Copy(cmd.OutOrStdout(), bytes.NewReader(out))
				return err
			}

			if err = appGenesis.SaveAs(outputDocument); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Int64(FlagHeight, -1, "Export state from a particular height (-1 means latest height)")
	cmd.Flags().StringSlice(FlagJailAllowedAddrs, []string{}, "Comma-separated list of operator addresses of jailed validators to unjail")
	cmd.Flags().StringSlice(server.FlagModulesToExport, []string{}, "Comma-separated list of modules to export. If empty, will export all modules")
	cmd.Flags().String(flags.FlagOutputDocument, "", "Exported state is written to the given file instead of STDOUT")

	return cmd
}

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
