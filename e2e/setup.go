package e2e

import (
	"context"
	"fmt"
	"testing"

	// fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

var (
	coinType = "639"
	denom    = "ubtsg"

	BitsongE2ERepo  = "ghcr.io/bitsongofficial/go-bitsong-e2e"
	BitsongMainRepo = "ghcr.io/bitsongofficial/go-bitsong"

	UpgradeFromBitsongImage = ibc.DockerImage{
		Repository: "bitsong",
		Version:    "v0.17.0",
		UidGid:     "1025:1025",
	}

	CurrentBranchBitsongImage = ibc.DockerImage{
		Repository: "bitsong",
		Version:    "local",
		UidGid:     "1025:1025",
	}

	// // SDK v47 Genesis
	defaultGenesisKV = []cosmos.GenesisKV{
		// {
		// 	Key:   "app_state.gov.params.voting_period",
		// 	Value: VotingPeriod,
		// },
		// {
		// 	Key:   "app_state.gov.params.max_deposit_period",
		// 	Value: MaxDepositPeriod,
		// },
		// {
		// 	Key:   "app_state.gov.params.min_deposit.0.denom",
		// 	Value: Denom,
		// },
		// {
		// 	Key:   "app_state.feepay.params.enable_feepay",
		// 	Value: false,
		// },
	}

	bitsongCfg = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "bitsong",
		ChainID:             "bitsong-local-1",
		Images:              []ibc.DockerImage{UpgradeFromBitsongImage},
		Bin:                 "bitsongd",
		Bech32Prefix:        "bitsong",
		Denom:               denom,
		CoinType:            coinType,
		GasPrices:           fmt.Sprintf("0%s", denom),
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		SkipGenTx:           false,
		PreGenesis:          nil,
		EncodingConfig:      bitsongEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
		ConfigFileOverrides: nil,
	}
)

// bitsongEncoding registers the Bitsong specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func bitsongEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	// fantokentypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

// generates this branch image
func CreateThisBranchChain(t *testing.T, numVals, numFull int) []ibc.Chain {
	return CreateChain(t, numVals, numFull, CurrentBranchBitsongImage)
}

// generates custom chain version image
func CreateChain(t *testing.T, numVals, numFull int, img ibc.DockerImage) []ibc.Chain {
	cfg := bitsongCfg
	cfg.Images = []ibc.DockerImage{img}
	return CreateChainWithCustomConfig(t, numVals, numFull, cfg)
}

func CreateChainWithCustomConfig(t *testing.T, numVals, numFull int, config ibc.ChainConfig) []ibc.Chain {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "bitsong",
			ChainName:     "bitsong",
			Version:       config.Images[0].Version,
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
