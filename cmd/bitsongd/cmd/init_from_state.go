package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	errorsmod "cosmossdk.io/errors"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/cli"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/go-bip39"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctypes "github.com/cosmos/ibc-go/v7/modules/core/types"
	"github.com/spf13/cobra"
)

const (
	FlagOldMoniker         = "old-moniker"
	FlagOldAccountAddr     = "old-account-addr"
	FlagVotingPeriod       = "voting-period"
	FlagPruneIbc           = "prune-ibc"
	FlagIncreaseCoinAmount = "increase-coin-amount"
)

// InitFromStateCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitFromStateCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-from-state [moniker] state_exported.json [key-name]",
		Short: "Initialize private validator, p2p, genesis, application configuration and replace an exported state",
		Long:  `Initialize validators's, node's configuration files and replace an exported state.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// Override default settings in config.toml
			config.P2P.MaxNumInboundPeers = 100
			config.P2P.MaxNumOutboundPeers = 50
			config.Mempool.Size = 10000
			config.StateSync.TrustPeriod = 112 * time.Hour
			config.BlockSync.Version = "v0"

			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}

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

			_, pubKey, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)
			moniker := args[0]
			stateFile := args[1]
			keyName := args[2]
			oldMoniker, _ := cmd.Flags().GetString(FlagOldMoniker)
			oldAccountAddr, _ := cmd.Flags().GetString(FlagOldAccountAddr)
			votingPeriodFlag, _ := cmd.Flags().GetString(FlagVotingPeriod)
			votingPeriod, err := time.ParseDuration(votingPeriodFlag)
			if err != nil {
				return err
			}
			pruneIbc, _ := cmd.Flags().GetBool(FlagPruneIbc)
			increaseCoinAmount, _ := cmd.Flags().GetInt64(FlagIncreaseCoinAmount)

			// attempt to lookup address from Keybase if no address was provided
			kb, err := keyring.New(sdk.KeyringServiceName(), "test", clientCtx.HomeDir, bufio.NewReader(cmd.InOrStdin()), clientCtx.Codec)
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}
			clientCtx.Keyring = kb

			tmPubKey, _ := cryptocodec.ToTmPubKeyInterface(pubKey)
			genParams := StateExportParams{
				ChainID:            chainID,
				Moniker:            moniker,
				KeyName:            keyName,
				IncreaseCoinAmount: increaseCoinAmount,
				OldMoniker:         oldMoniker,
				TmPubKey:           tmPubKey,
				OldAccountAddress:  oldAccountAddr,
				VotingPeriod:       votingPeriod,
				PruneIBC:           pruneIbc,
				Overwrite:          overwrite,
				StateFile:          stateFile,
			}

			return initNodeFromState(config, clientCtx, genParams)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")

	cmd.Flags().String(FlagOldMoniker, "", "the validator moniker to replace")
	cmd.Flags().String(FlagOldAccountAddr, "", "the account address to replace")
	cmd.Flags().String(FlagVotingPeriod, "180s", "the voting period for the governance proposal")
	cmd.Flags().Bool(FlagPruneIbc, true, "prune the ibc state")
	cmd.Flags().Int64(FlagIncreaseCoinAmount, 1000000000000000, "increase the bonded token amount")

	return cmd
}

func initNodeFromState(config *tmcfg.Config, cliCtx client.Context, genParams StateExportParams) error {
	genFile := config.GenesisFile()

	if !genParams.Overwrite && tmos.FileExists(genFile) {
		return fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	genDoc, err := ConvertStateExport(cliCtx, genParams)
	if err != nil {
		return fmt.Errorf("failed to convert state export: %w", err)
	}

	// validate genesis state
	// TODO: fix this!
	/*if err = mbm.ValidateGenesis(cliCtx.Codec, cliCtx.TxConfig, appState); err != nil {
		return fmt.Errorf("error validating genesis file: %s", err.Error())
	}*/

	// save genesis
	if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
		return fmt.Errorf("failed to export gensis file: %s", err.Error())
	}

	config.BlockSyncMode = false
	tmcfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	fmt.Println("State imported successfully")

	return nil
}

type StateExportParams struct {
	StateFile          string
	ChainID            string
	Moniker            string
	KeyName            string
	IncreaseCoinAmount int64
	OldMoniker         string
	TmPubKey           tmcrypto.PubKey
	OldAccountAddress  string
	VotingPeriod       time.Duration
	PruneIBC           bool
	Overwrite          bool
}

func (s *StateExportParams) String() string {
	return fmt.Sprintf(`State Export Params:
  State File: %s
  Chain ID: %s
  Moniker: %s
  Key Name: %s
  Increase Coin Amount: %d
  Old Moniker: %s
  Old Account Address: %s
  Voting Period: %s
  Prune IBC: %t
  Overwrite: %t
`, s.StateFile, s.ChainID, s.Moniker, s.KeyName, s.IncreaseCoinAmount, s.OldMoniker, s.OldAccountAddress, s.VotingPeriod, s.PruneIBC, s.Overwrite)
}

func (s *StateExportParams) Validate() error {
	if s.StateFile == "" {
		return fmt.Errorf("state file cannot be empty")
	}

	if s.Moniker == "" {
		return fmt.Errorf("moniker cannot be empty")
	}

	if s.KeyName == "" {
		return fmt.Errorf("key name cannot be empty")
	}

	if s.OldMoniker == "" {
		return fmt.Errorf("--%s cannot be empty", FlagOldMoniker)
	}

	if s.OldAccountAddress == "" {
		return fmt.Errorf("--%s cannot be empty", FlagOldAccountAddr)
	}

	if s.TmPubKey == nil {
		return fmt.Errorf("tm pub key cannot be empty")
	}

	if s.VotingPeriod == 0 {
		return fmt.Errorf("voting period cannot be zero")
	}

	return nil
}

func ConvertStateExport(clientCtx client.Context, params StateExportParams) (*tmtypes.GenesisDoc, error) {
	if !tmos.FileExists(params.StateFile) {
		return nil, fmt.Errorf("%s does not exist", params.StateFile)
	}

	// print state export params
	fmt.Println(params.String())

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	stateBz, err := ioutil.ReadFile(params.StateFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read state export file: %w", err)
	}

	// replace account
	newAccount, err := fetchKey(clientCtx.Keyring, params.KeyName)
	if err != nil {
		return nil, fmt.Errorf("couldn't find key %s", params.KeyName)
	}

	// Update the pub_key
	str := fmt.Sprintf(`"address"\s*:\s*"(%s)".*?"key"\s*:\s*"(.*?)"`, params.OldAccountAddress)
	re := regexp.MustCompile(str)
	match := re.FindStringSubmatch(string(stateBz))
	if len(match) > 1 {
		oldPubKey := match[2]
		pubkey, _ := newAccount.GetPubKey()
		newPubKey := base64.StdEncoding.EncodeToString(pubkey.Bytes())
		stateBz = bytes.Replace(stateBz, []byte(oldPubKey), []byte(newPubKey), -1)
	} else {
		panic("pub_key not found")
	}
	// Update address
	addr, _ := newAccount.GetAddress()
	stateBz = bytes.Replace(stateBz, []byte(params.OldAccountAddress), []byte(addr.String()), -1)

	genDoc := tmtypes.GenesisDoc{}
	err = tmjson.Unmarshal(stateBz, &genDoc)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal state export file: %w", err)
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return nil, err
	}

	appState, err := genutiltypes.GenesisStateFromGenDoc(genDoc)
	if err != nil {
		return nil, err
	}

	// the logic is based on testnetify.py

	// Replace chain-id
	genDoc.ChainID = params.ChainID

	// Update gov module
	// var govGenState = gov.ExportGenesis(clientCtx, Govkeeper)
	// clientCtx.Codec.MustUnmarshalJSON(appState[govtypes.ModuleName], &govGenState)
	// govGenState.VotingParams.VotingPeriod = params.VotingPeriod
	// govGenStateBz, err := clientCtx.Codec.MarshalJSON(&govGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	// appState[govtypes.ModuleName] = govGenStateBz

	// Prune IBC
	if params.PruneIBC {
		var ibcGenState ibctypes.GenesisState
		clientCtx.Codec.MustUnmarshalJSON(appState["ibc"], &ibcGenState)

		ibcGenState.ChannelGenesis.AckSequences = []ibcchanneltypes.PacketSequence{}
		ibcGenState.ChannelGenesis.Acknowledgements = []ibcchanneltypes.PacketState{}
		ibcGenState.ChannelGenesis.Channels = []ibcchanneltypes.IdentifiedChannel{}
		ibcGenState.ChannelGenesis.Commitments = []ibcchanneltypes.PacketState{}
		ibcGenState.ChannelGenesis.Receipts = []ibcchanneltypes.PacketState{}
		ibcGenState.ChannelGenesis.RecvSequences = []ibcchanneltypes.PacketSequence{}
		ibcGenState.ChannelGenesis.SendSequences = []ibcchanneltypes.PacketSequence{}
		ibcGenState.ChannelGenesis.NextChannelSequence = uint64(1)

		ibcGenState.ClientGenesis.Clients = []ibcclienttypes.IdentifiedClientState{}
		ibcGenState.ClientGenesis.ClientsConsensus = ibcclienttypes.ClientsConsensusStates{}
		ibcGenState.ClientGenesis.ClientsMetadata = []ibcclienttypes.IdentifiedGenesisMetadata{}
		ibcGenState.ClientGenesis.NextClientSequence = uint64(1)

		ibcGenStateBz, err := clientCtx.Codec.MarshalJSON(&ibcGenState)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ibc genesis state: %w", err)
		}
		appState["ibc"] = ibcGenStateBz
	}

	// Impersonate validator
	var oldValidator tmtypes.GenesisValidator

	// Update tendermint validator data
	for i, _ := range genDoc.Validators {
		if genDoc.Validators[i].Name == params.OldMoniker {
			oldValidator = genDoc.Validators[i]
			validator := &genDoc.Validators[i]

			// Replace validator data
			validator.PubKey = params.TmPubKey
			validator.Address = params.TmPubKey.Address()
			validator.Power = validator.Power + (params.IncreaseCoinAmount / 1000000)
		}
	}
	if oldValidator.Name == "" {
		return nil, fmt.Errorf("validator to replace %s not found", params.OldMoniker)
	}

	// Update staking module
	var stakingGenState stakingtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appState[stakingtypes.ModuleName], &stakingGenState)
	operatorAddr := ""
	for i, _ := range stakingGenState.Validators {
		validator := &stakingGenState.Validators[i]
		if validator.Description.Moniker == params.OldMoniker {
			valPubKey, err := cryptocodec.FromTmPubKeyInterface(params.TmPubKey)
			if err != nil {
				return nil, fmt.Errorf("failed to convert validator pubkey: %w", err)
			}

			newConsensusPubKey, err := codectypes.NewAnyWithValue(valPubKey)
			if err != nil {
				return nil, err
			}
			validator.ConsensusPubkey = newConsensusPubKey

			operatorAddr = validator.OperatorAddress
			validator.OperatorAddress = sdk.ValAddress(params.TmPubKey.Address()).String()
			validator.DelegatorShares = validator.DelegatorShares.Add(sdk.NewDec(params.IncreaseCoinAmount))
			validator.Tokens = validator.Tokens.Add(sdk.NewInt(params.IncreaseCoinAmount))
		}
	}

	// Update total power
	for i, _ := range stakingGenState.LastValidatorPowers {
		validatorPower := &stakingGenState.LastValidatorPowers[i]
		if validatorPower.Address == operatorAddr {
			validatorPower.Power = validatorPower.Power + (params.IncreaseCoinAmount / 1000000)
		}
	}
	stakingGenState.LastTotalPower = stakingGenState.LastTotalPower.Add(sdk.NewInt(params.IncreaseCoinAmount / 1000000))

	// Update self delegation on operator address
	for i, _ := range stakingGenState.Delegations {
		delegation := &stakingGenState.Delegations[i]
		new_account, err := newAccount.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed: %w", err)
		}
		if delegation.DelegatorAddress == new_account.String() {
			delegation.Shares = delegation.Shares.Add(sdk.NewDec(params.IncreaseCoinAmount))
		}
	}

	stakingGenStateBz, err := clientCtx.Codec.MarshalJSON(&stakingGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// Update genesis['app_state']['distribution']['delegator_starting_infos'] on operator address
	var distrGenState distrtypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appState[distrtypes.ModuleName], &distrGenState)
	for i, _ := range distrGenState.DelegatorStartingInfos {
		delegatorStartingInfo := &distrGenState.DelegatorStartingInfos[i]
		new_account, err := newAccount.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed: %w", err)
		}
		if delegatorStartingInfo.DelegatorAddress == new_account.String() {
			delegatorStartingInfo.StartingInfo.Stake = delegatorStartingInfo.StartingInfo.Stake.Add(sdk.NewDec(params.IncreaseCoinAmount))
		}
	}
	distrGenStateBz, err := clientCtx.Codec.MarshalJSON(&distrGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distrtypes.ModuleName] = distrGenStateBz

	// Update bank module
	var bankGenState banktypes.GenesisState
	clientCtx.Codec.MustUnmarshalJSON(appState[banktypes.ModuleName], &bankGenState)

	for i, _ := range bankGenState.Balances {
		balance := &bankGenState.Balances[i]

		// Add 1 BN ubtsg to bonded_tokens_pool module address
		bondedPool := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
		if balance.Address == bondedPool.String() {
			for balanceIdx, _ := range balance.Coins {
				coin := &balance.Coins[balanceIdx]
				if coin.Denom == stakingGenState.Params.BondDenom {
					coin.Amount = coin.Amount.Add(sdk.NewInt(params.IncreaseCoinAmount))
				}
			}
		}
	}

	// Update bank balance
	for i, _ := range bankGenState.Supply {
		supply := &bankGenState.Supply[i]
		if supply.Denom == stakingGenState.Params.BondDenom {
			supply.Amount = supply.Amount.Add(sdk.NewInt(params.IncreaseCoinAmount))
		}
	}

	bankGenStateBz, err := clientCtx.Codec.MarshalJSON(&bankGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	// Replace old validator_address with the new one
	appStateJSON = bytes.Replace(appStateJSON, []byte(operatorAddr), []byte(sdk.ValAddress(params.TmPubKey.Address()).String()), -1)
	appStateJSON = bytes.Replace(appStateJSON, []byte(sdk.ConsAddress(oldValidator.PubKey.Address()).String()), []byte(sdk.ConsAddress(params.TmPubKey.Address()).String()), -1)

	genDoc.AppState = appStateJSON

	return &genDoc, nil
}

func fetchKey(kb keyring.Keyring, keyref string) (keyring.Record, error) {
	// firstly check if the keyref is a key name of a key registered in a keyring.
	k, err := kb.Key(keyref)
	// if the key is not there or if we have a problem with a keyring itself then we move to a
	// fallback: searching for key by address.

	if err == nil || !errorsmod.IsOf(err, sdkerr.ErrIO, sdkerr.ErrKeyNotFound) {
		return *k, err
	}

	accAddr, err := sdk.AccAddressFromBech32(keyref)
	if err != nil {
		return *k, err
	}

	k, err = kb.KeyByAddress(accAddr)
	return *k, errorsmod.Wrap(err, "Invalid key")
}
