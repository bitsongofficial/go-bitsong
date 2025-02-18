package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tmcfg "github.com/cometbft/cometbft/config"
	cmttypes "github.com/cometbft/cometbft/types"

	"github.com/cometbft/cometbft/libs/cli"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagSeed defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))

	return err
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// P2P
			seeds := []string{
				"ade4d8bc8cbe014af6ebdf3cb7b1e9ad36f412c0@seeds.polkachu.com:16056",      // Polkachu
				"8542cd7e6bf9d260fef543bc49e59be5a3fa90740@seed.publicnode.com:26656",    // Allnodes team
				"8defec7d0eec97f507411e02fd2634e3efc997a2@bitsong-seed.panthea.eu:41656", // Panthea EU
			}
			config.P2P.Seeds = strings.Join(seeds, ",")
			config.P2P.MaxNumInboundPeers = 80
			config.P2P.MaxNumOutboundPeers = 60

			// Mempool
			config.Mempool.Size = 10000

			// State Sync
			config.StateSync.TrustPeriod = 112 * time.Hour
			config.BlockSync.Version = "v0"

			// // Consensus
			// config.Consensus.TimeoutCommit = 1500 * time.Millisecond // 1.5s

			//  other
			config.Moniker = args[0]
			config.SetRoot(clientCtx.HomeDir)

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				mnemonic, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			genFile := config.GenesisFile()
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			var toPrint printInfo
			isMainnet := chainID == "" || chainID == "bitsong-2b"
			genesisFileDownloadFailed := false
			if isMainnet {
				// download mainnet genesis file, if fail make new one
				err := downloadGenesis(config)
				if err != nil {
					fmt.Println("Failed to download genesis file, using a random chain ID and genesis file for local testing")
					genesisFileDownloadFailed = true
					chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
				} else {
					// Set chainID to bitsong-2b in the case of a blank chainID
					chainID = "bitsong-2b"
					// We dont print the app state for mainnet nodes because it's massive
					fmt.Println("Not printing app state for mainnet node due to verbosity")
					toPrint = newPrintInfo(config.Moniker, chainID, nodeID, "", nil)
				}
			}

			if genesisFileDownloadFailed || !isMainnet {
				var genDoc genutiltypes.AppGenesis

				appStateJSON, err := json.MarshalIndent(mbm.DefaultGenesis(clientCtx.Codec), "", " ")
				if err != nil {
					return errors.Wrap(err, "Failed to marshall default genesis state")
				}

				if _, err := os.Stat(genFile); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					_, genDocFromFile, err := genutiltypes.GenesisStateFromGenFile(genFile)
					if err != nil {
						return fmt.Errorf("failed to unmarshal genesis state: %w", err)
					}
					genDoc = *genDocFromFile
				}

				genDoc.Consensus = &genutiltypes.ConsensusGenesis{}
				genDoc.ChainID = chainID
				genDoc.Consensus.Params = cmttypes.DefaultConsensusParams()
				genDoc.Consensus.Validators = nil
				genDoc.AppState = appStateJSON

				if err = genutil.ExportGenesisFile(&genDoc, genFile); err != nil {
					return errors.Wrap(err, "Failed to export gensis file")
				}
				toPrint = newPrintInfo(config.Moniker, chainID, nodeID, "", appStateJSON)
			}

			tmcfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

// downloadGenesis downloads the genesis file from a predefined URL and writes it to the genesis file path specified in the config.
// It creates an HTTP client to send a GET request to the genesis file URL. If the request is successful, it reads the response body
// and writes it to the destination genesis file path. If any step in this process fails, it generates the default genesis.
//
// Parameters:
// - config: A pointer to a tmcfg.Config object that contains the configuration, including the genesis file path.
//
// Returns:
// - An error if the download or file writing fails, otherwise nil.
func downloadGenesis(config *tmcfg.Config) error {
	// URL of the genesis file to download
	genesisURL := "https://raw.githubusercontent.com/bitsongofficial/networks/refs/heads/master/bitsong-2b/genesis.json?download"

	// Determine the destination path for the genesis file
	genFilePath := config.GenesisFile()

	// Create a new HTTP client with a 30-second timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", genesisURL, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request for genesis file")
	}

	// Send the request
	fmt.Println("Downloading genesis file from", genesisURL)
	fmt.Println("If the download is not successful in 30 seconds, we will gracefully continue and the default genesis file will be used")
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to download genesis file")
	}
	defer resp.Body.Close()

	// Check if the HTTP request was successful
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed to download genesis file: HTTP status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read genesis file response body")
	}

	// Write the body to the destination genesis file
	err = os.WriteFile(genFilePath, body, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write genesis file to destination")
	}

	return nil
}
