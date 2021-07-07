package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bitsongofficial/chainmodules/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	cryptocodec "github.com/tendermint/tendermint/crypto/encoding"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"log"
	"time"
)

const (
	flagGenesisTime   = "genesis-time"
	flagInitialHeight = "initial-height"
	flagTestnetReplaceKeys = "testnet-replace-cons-keys"
)

func migrateCmd() *cobra.Command {
	cmd := cobra.Command{
		Use: "migrate [genesis-file]",
		Short: "Migrate Genesis File from v0.7 to v0.8",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.
Example:
$ %s migrate /path/to/genesis.json --chain-id=bitsong-2 --genesis-time=2021-07-01T12:00:00Z --initial-height=1000000
`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var ctx = client.GetClientContextFromCmd(cmd)

			genesisBlob, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			chainID, err := cmd.Flags().GetString(flags.FlagChainID)
			if err != nil {
				return err
			}

			genesisTime, err := cmd.Flags().GetString(flagGenesisTime)
			if err != nil {
				return err
			}

			initialHeight, err := cmd.Flags().GetInt64(flagInitialHeight)
			if err != nil {
				return err
			}

			genesis, err := tmtypes.GenesisDocFromJSON(genesisBlob)
			if err != nil {
				return err
			}

			var currentState genutiltypes.AppMap
			if err := json.Unmarshal(genesis.AppState, &currentState); err != nil {
				return err
			}

			migrateFn := cli.GetMigrationCallback("v0.40")
			if migrateFn == nil {
				return fmt.Errorf("sdk migration function is not available")
			}

			currentState = migrateFn(currentState, ctx)

			var (
				bankGenesis banktypes.GenesisState
				stakingGenesis stakingtypes.GenesisState
				ibcTransferGenesis = ibctransfertypes.DefaultGenesisState()
				ibcGenesis         = ibctypes.DefaultGenesisState()
				capabilityGenesis  = capabilitytypes.DefaultGenesis()
				evidenceGenesis    = evidencetypes.DefaultGenesisState()
			)
			ctx.JSONMarshaler.MustUnmarshalJSON(currentState[banktypes.ModuleName], &bankGenesis)
			ctx.JSONMarshaler.MustUnmarshalJSON(currentState[stakingtypes.ModuleName], &stakingGenesis)

			// TODO: replace BondDenom with new types
			bankGenesis.DenomMetadata = []banktypes.Metadata{
				{
					Description: "The BitSong’s Network native coin",
					Base: types.BondDenom,
					Display: types.BondDenom[1:],
					DenomUnits: []*banktypes.DenomUnit{
						{Denom: types.BondDenom[1:], Exponent: uint32(6), Aliases: []string{}},
						{Denom: types.BondDenom, Exponent: uint32(0), Aliases: []string{"microbtsg"}},
					},
				},
			}

			stakingGenesis.Params.HistoricalEntries = 10000

			ibcTransferGenesis.Params.ReceiveEnabled = false
			ibcTransferGenesis.Params.SendEnabled = false
			ibcGenesis.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}

			currentState[banktypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(&bankGenesis)
			currentState[ibctransfertypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(ibcTransferGenesis)
			currentState[ibchost.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(ibcGenesis)
			currentState[capabilitytypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(capabilityGenesis)
			currentState[evidencetypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(evidenceGenesis)
			currentState[stakingtypes.ModuleName] = ctx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)

			genesis.AppState, err = json.Marshal(currentState)
			if err != nil {
				return err
			}

			if genesisTime != "" {
				var t time.Time
				if err := t.UnmarshalText([]byte(genesisTime)); err != nil {
					return err
				}

				genesis.GenesisTime = t
			}
			if chainID != "" {
				genesis.ChainID = chainID
			}

			genesis.InitialHeight = initialHeight

			// Replace validator keys for testnet
			replacementKeys, _ := cmd.Flags().GetString(flagTestnetReplaceKeys)

			if replacementKeys != "" {
				genesis = loadKeydataFromFile(ctx, replacementKeys, genesis)
			}

			genesisBlob, err = tmjson.Marshal(genesis)
			if err != nil {
				return err
			}

			sortedGenesisBlob, err := sdk.SortJSON(genesisBlob)
			if err != nil {
				return err
			}

			fmt.Println(string(sortedGenesisBlob))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagChainID, "", "set chain id")
	cmd.Flags().String(flagGenesisTime, "", "set genesis time")
	cmd.Flags().Int64(flagInitialHeight, 0, "set the initial height")
	cmd.Flags().String(flagTestnetReplaceKeys, "", "Provide a JSON file to replace the consensus keys of validators")

	return &cmd
}

type replacementConfig struct {
	Name             string `json:"validator_name"`
	ValidatorAddress string `json:"validator_address"`
	ConsensusPubkey  string `json:"testnet_consensus_public_key"`
}

type replacementConfigs []replacementConfig

func (r *replacementConfigs) isReplacedValidator(validatorAddress string) (int, replacementConfig) {

	for i, replacement := range *r {
		if replacement.ValidatorAddress == validatorAddress {
			return i, replacement
		}
	}

	return -1, replacementConfig{}
}

func loadKeydataFromFile(clientCtx client.Context, replacementrJSON string, genDoc *tmtypes.GenesisDoc) *tmtypes.GenesisDoc {
	jsonReplacementBlob, err := ioutil.ReadFile(replacementrJSON)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read replacement keys from file %s", replacementrJSON))
	}

	var replacementKeys replacementConfigs

	err = json.Unmarshal(jsonReplacementBlob, &replacementKeys)

	if err != nil {
		log.Fatal("Could not unmarshal replacement keys ")
	}

	var state genutiltypes.AppMap
	if err := json.Unmarshal(genDoc.AppState, &state); err != nil {
		log.Fatal(errors.Wrap(err, "failed to JSON unmarshal initial genesis state"))
	}

	var stakingGenesis stakingtypes.GenesisState
	var slashingGenesis slashingtypes.GenesisState

	clientCtx.JSONMarshaler.MustUnmarshalJSON(state[stakingtypes.ModuleName], &stakingGenesis)
	clientCtx.JSONMarshaler.MustUnmarshalJSON(state[slashingtypes.ModuleName], &slashingGenesis)

	for i, val := range stakingGenesis.Validators {
		idx, replacement := replacementKeys.isReplacedValidator(val.OperatorAddress)

		if idx != -1 {

			toReplaceValConsAddress, _ := val.GetConsAddr()

			consPubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, replacement.ConsensusPubkey)

			if err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", consPubKey, err))
			}

			val.ConsensusPubkey, err = codectypes.NewAnyWithValue(consPubKey)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", consPubKey, err))
			}

			replaceValConsAddress, _ := val.GetConsAddr()
			protoReplaceValConsPubKey, _ := val.TmConsPublicKey()
			replaceValConsPubKey, _ := cryptocodec.PubKeyFromProto(protoReplaceValConsPubKey)

			for i, signingInfo := range slashingGenesis.SigningInfos {
				if signingInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.SigningInfos[i].Address = replaceValConsAddress.String()
					slashingGenesis.SigningInfos[i].ValidatorSigningInfo.Address = replaceValConsAddress.String()
				}
			}

			for i, missedInfo := range slashingGenesis.MissedBlocks {
				if missedInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.MissedBlocks[i].Address = replaceValConsAddress.String()
				}
			}

			for tmIdx, tmval := range genDoc.Validators {
				if tmval.Address.String() == replaceValConsAddress.String() {
					genDoc.Validators[tmIdx].Address = replaceValConsAddress.Bytes()
					genDoc.Validators[tmIdx].PubKey = replaceValConsPubKey

				}
			}
			stakingGenesis.Validators[i] = val

		}

	}
	state[stakingtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&stakingGenesis)
	state[slashingtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&slashingGenesis)

	genDoc.AppState, err = json.Marshal(state)

	if err != nil {
		log.Fatal("Could not marshal App State")
	}
	return genDoc
}