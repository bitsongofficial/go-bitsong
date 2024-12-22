package e2e

import (
	"context"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/docker/docker/client"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	ibclocalhost "github.com/cosmos/ibc-go/v7/modules/light-clients/09-localhost"
)

var (
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"
	Denom            = "ubtsg"

	BitsongE2eRepo  = "ghcr.io/bitsongofficial/go-bitsong-e2e"
	BitsongMainRepo = "ghcr.io/bitsongofficial/go-bitsong"

	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "main"

	bitsongRepo, bitsongVersion = GetDockerImageInfo()

	BitsongImage = ibc.DockerImage{
		Repository: bitsongRepo,
		Version:    bitsongVersion,
		UidGid:     "1025:1025",
	}

	// SDK v47 Genesis
	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: Denom,
		},
	}

	baseCfg = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "bitsong",
		ChainID:             "bitsong-2",
		Images:              []ibc.DockerImage{BitsongImage},
		Bin:                 "bitsongd",
		Bech32Prefix:        "bitsong",
		Denom:               "ubtsg",
		CoinType:            "118",
		GasPrices:           "0ubtsg",
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      btsgEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
	}
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("bitsong", "bitsong")
	sdk.GetConfig().SetBech32PrefixForValidator("bitsongvaloper", "bitsong")
	sdk.GetConfig().SetBech32PrefixForConsensusNode("bitsongvalcons", "bitsong")
	sdk.GetConfig().SetCoinType(118)
}

// btsgEncoding registers the Bitsong specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func btsgEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	ibclocalhost.RegisterInterfaces(cfg.InterfaceRegistry)
	wasmtypes.RegisterInterfaces(cfg.InterfaceRegistry)
	// fantokentypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

// CreateChain generates a new chain with a custom image (useful for upgrades)
func CreateChain(t *testing.T, numVals, numFull int, img ibc.DockerImage) []ibc.Chain {
	cfg := baseCfg
	cfg.Images = []ibc.DockerImage{img}
	return CreateICTestBitsongChainCustomConfig(t, numVals, numFull, cfg)
}

// CreateThisBranchWithValsAndFullNodes generates this branch's chain (ex: from the commit), with a set of validators and full nodes.
func CreateICTestBitsongChain(t *testing.T, numVals, numFull int) []ibc.Chain {
	return CreateChain(t, numVals, numFull, BitsongImage)
}

func CreateICTestBitsongChainCustomConfig(t *testing.T, numVals, numFull int, config ibc.ChainConfig) []ibc.Chain {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "bitsong",
			ChainConfig:   config,
			NumValidators: &numVals,
			NumFullNodes:  &numFull,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	// chain := chains[0].(*cosmos.CosmosChain)
	return chains
}

func BuildInitialChain(t *testing.T, chains []ibc.Chain) (*interchaintest.Interchain, context.Context, *client.Client, string) {
	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()

	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err := ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	return ic, ctx, client, network
}
