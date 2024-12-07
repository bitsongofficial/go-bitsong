// source: https://github.com/DA0-DA0/polytone/blob/main/tests/strangelove/suite.go
package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"
	w "github.com/CosmWasm/wasmvm/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type Suite struct {
	t        *testing.T
	reporter *testreporter.RelayerExecReporter
	ctx      context.Context

	ChainA SuiteChain
	ChainB SuiteChain

	Relayer ibc.Relayer
	PathAB  string
}

type SuiteChain struct {
	Ibc    ibc.Chain
	Cosmos *cosmos.CosmosChain
	User   ibc.Wallet

	Note   string
	Voice  string
	Tester string
}

func NewPolytoneSuite(t *testing.T) Suite {
	if testing.Short() {
		t.Skip()
	}

	var (
		ctx                  = context.Background()
		client, network      = interchaintest.DockerSetup(t)
		rep                  = testreporter.NewNopReporter()
		eRep                 = rep.RelayerExecReporter(t)
		chainID_A, chainID_B = "chain-a", "chain-b"
		chainA, chainB       *cosmos.CosmosChain
	)

	// base config which all networks will use as defaults.
	baseCfg := ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "bitsongd",
		ChainID:             "", // change this for each
		Images:              []ibc.DockerImage{BitsongImage},
		Bin:                 "bitsongd",
		Bech32Prefix:        "bitsong",
		Denom:               "ubtsg",
		CoinType:            "639",
		GasPrices:           "0ubtsg",
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      btsgEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(defaultGenesisKV),
	}

	// Set specific chain ids for each so they are their own unique networks
	baseCfg.ChainID = chainID_A
	configA := baseCfg

	baseCfg.ChainID = chainID_B
	configB := baseCfg

	// Create chain factory with multiple Bitsong individual networks.
	numVals := 1
	numFullNodes := 0

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "bitsong",
			ChainConfig:   configA,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "bitsong",
			ChainConfig:   configB,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chainA, chainB = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	).Build(t, client, network)

	const pathAB = "ab"

	ic := interchaintest.NewInterchain().
		AddChain(chainA).
		AddChain(chainB).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainA,
			Chain2:  chainB,
			Relayer: r,
			Path:    pathAB,
		})

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	userFunds := sdkmath.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chainA, chainB)

	// abChan, err := ibc.GetTransferChannel(ctx, r, eRep, chainID_A, chainID_B)
	// require.NoError(t, err)

	// baChan := abChan.Counterparty

	// Start the relayer on all paths
	err = r.StartRelayer(ctx, eRep, pathAB)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Get original account balances
	userA, userB := users[0], users[1]
	t.Logf("userA: %s", userA)

	suite := Suite{
		t:        t,
		reporter: eRep,
		ctx:      ctx,
		ChainA: SuiteChain{
			Ibc:    chainA,
			Cosmos: chainA,
			User:   userA,
		},
		ChainB: SuiteChain{
			Ibc:    chainB,
			Cosmos: chainB,
			User:   userB,
		},
		Relayer: r,
		PathAB:  pathAB,
	}

	suite.SetupChain(&suite.ChainA)
	suite.SetupChain(&suite.ChainB)
	return suite
}

func (s *Suite) SetupChain(chain *SuiteChain) {
	user := chain.User
	cc := chain.Cosmos
	noteId, err := cc.StoreContract(s.ctx, user.KeyName(), "contracts/polytone_note.wasm")
	if err != nil {
		s.t.Fatal(err)
	}
	voiceId, err := cc.StoreContract(s.ctx, user.KeyName(), "contracts/polytone_voice.wasm")
	if err != nil {
		s.t.Fatal(err)
	}
	proxyId, err := cc.StoreContract(s.ctx, user.KeyName(), "contracts/polytone_proxy.wasm")
	if err != nil {
		s.t.Fatal(err)
	}

	testerId, err := cc.StoreContract(s.ctx, user.KeyName(), "contracts/polytone_tester.wasm")
	if err != nil {
		s.t.Fatal(err)
	}

	proxyUint, err := strconv.Atoi(proxyId)
	if err != nil {
		s.t.Fatal(err)
	}

	chain.Note = s.Instantiate(cc, user, noteId, NoteInstantiate{
		BlockMaxGas: 100_000_000,
	})
	chain.Voice = s.Instantiate(cc, user, voiceId, VoiceInstantiate{
		ProxyCodeId:     uint64(proxyUint),
		BlockMaxGas:     100_000_000,
		ContractAddrLen: 32,
	})
	chain.Tester = s.Instantiate(cc, user, testerId, TesterInstantiate{})
}

func (s *Suite) Instantiate(chain *cosmos.CosmosChain, user ibc.Wallet, codeId string, msg any) string {
	str, err := json.Marshal(msg)
	if err != nil {
		s.t.Fatal(err)
	}

	address, err := chain.InstantiateContract(s.ctx, user.KeyName(), codeId, string(str), true)
	if err != nil {
		s.t.Fatal(err)
	}
	return address
}

func (s *Suite) CreateChannel(initModule string, tryModule string, initChain, tryChain *SuiteChain, path string) (initChannel, tryChannel string, err error) {
	initStartChannels := s.QueryOpenChannels(initChain)
	err = s.Relayer.CreateChannel(s.ctx, s.reporter, path, ibc.CreateChannelOptions{
		SourcePortName: "wasm." + initModule,
		DestPortName:   "wasm." + tryModule,
		Order:          ibc.Unordered,
		Version:        "polytone-1",
	})
	if err != nil {
		return
	}
	err = testutil.WaitForBlocks(s.ctx, 10, initChain.Ibc, tryChain.Ibc)
	if err != nil {
		return
	}

	initChannels := s.QueryOpenChannels(initChain)

	if len(initChannels) == len(initStartChannels) {
		err = errors.New("no new channels created")
		return
	}

	initChannel = initChannels[len(initChannels)-1].ChannelID
	tryChannel = initChannels[len(initChannels)-1].Counterparty.ChannelID
	return
}

const CHANNEL_STATE_OPEN = "STATE_OPEN"
const CHANNEL_STATE_TRY = "STATE_TRYOPEN"
const CHANNEL_STATE_INIT = "STATE_INIT"

func (s *Suite) QueryOpenChannels(chain *SuiteChain) []ibc.ChannelOutput {
	eq := s.QueryChannelsInState(chain, CHANNEL_STATE_OPEN)
	fmt.Println("QueryChannelsInState", eq)
	return eq
}

func (s *Suite) QueryChannelsInState(chain *SuiteChain, state string) []ibc.ChannelOutput {
	channels, err := s.Relayer.GetChannels(s.ctx, s.reporter, chain.Ibc.Config().ChainID)
	if err != nil {
		s.t.Fatal(err)
	}
	openChannels := []ibc.ChannelOutput{}
	for index := range channels {
		if channels[index].State == state {
			openChannels = append(openChannels, channels[index])
		}
	}
	return openChannels
}

func (s *Suite) RoundtripMessage(note string, chain *SuiteChain, msg NoteExecute) (Callback, error) {
	callbacksStart := s.QueryTesterCallbackHistory(&s.ChainA).History

	marshalled, err := json.Marshal(msg)
	if err != nil {
		return Callback{}, err
	}
	_, err = chain.Cosmos.ExecuteContract(s.ctx, chain.User.KeyName(), note, string(marshalled))
	if err != nil {
		return Callback{}, err
	}
	// wait for packet to relay.
	err = testutil.WaitForBlocks(s.ctx, 10, s.ChainA.Ibc, s.ChainB.Ibc)
	if err != nil {
		return Callback{}, err
	}
	callbacksEnd := s.QueryTesterCallbackHistory(&s.ChainA).History
	if len(callbacksEnd) == len(callbacksStart) {
		return Callback{}, errors.New("no new callback")
	}
	callback := callbacksEnd[len(callbacksEnd)-1]
	require.Equal(
		s.t,
		chain.User.Address(),
		callback.Initiator,
	)
	require.Equal(s.t, "aGVsbG8K", callback.InitiatorMsg)
	return callback.Result, nil
}

func (s *Suite) RoundtripExecute(note string, chain *SuiteChain, msgs []w.CosmosMsg) (Callback, error) {
	msg := NoteExecuteMsg{
		Msgs:           msgs,
		TimeoutSeconds: 100,
		Callback: &CallbackRequest{
			Receiver: chain.Tester,
			Msg:      "aGVsbG8K",
		},
	}
	return s.RoundtripMessage(note, chain, NoteExecute{
		Execute: &msg,
	})
}

func (s *Suite) RoundtripQuery(note string, chain *SuiteChain, msgs []w.CosmosMsg) (Callback, error) {
	msg := NoteQuery{
		Msgs:           msgs,
		TimeoutSeconds: 100,
		Callback: CallbackRequest{
			Receiver: chain.Tester,
			Msg:      "aGVsbG8K",
		},
	}
	return s.RoundtripMessage(note, chain, NoteExecute{
		Query: &msg,
	})
}

func (s *Suite) QueryTesterCallbackHistory(chain *SuiteChain) HistoryResponse {
	var response DataWrappedHistoryResponse
	query := TesterQuery{
		History: &Empty{},
	}
	err := chain.Cosmos.QueryContract(s.ctx, chain.Tester, query, &response)
	if err != nil {
		s.t.Fatal(err)
	}
	return response.Data
}
