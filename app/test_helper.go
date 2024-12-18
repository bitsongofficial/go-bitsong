package app

import (
	"encoding/json"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/bitsongofficial/go-bitsong/app/helpers"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/client"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/stretchr/testify/require"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// SimApp testing.
const (
	SimAppChainID = "testing"
)

// EmptyBaseAppOptions is a stub implementing AppOptions
type EmptyBaseAppOptions struct{}

// Get implements AppOptions
func (ao EmptyBaseAppOptions) Get(_ string) interface{} {
	return nil
}

var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// Setup initializes a new BitsongApp
func Setup(t *testing.T) *BitsongApp {
	t.Helper()

	privVal := helpers.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(appparams.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}

	app := SetupWithGenesisAccounts(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	return app
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} { return nil }

func setup(t *testing.T, withGenesis bool, opts ...wasmkeeper.Option) (*BitsongApp, GenesisState) {
	db := dbm.NewMemDB()
	nodeHome := t.TempDir()
	snapshotDir := filepath.Join(nodeHome, "data", "snapshots")
	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	require.NoError(t, err)
	t.Cleanup(func() { snapshotDB.Close() })
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	require.NoError(t, err)

	// var emptyWasmOpts []wasm.Option
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = nodeHome // ensure unique folder

	app := NewBitsongApp(log.NewNopLogger(), db, nil, true, EmptyAppOptions{}, opts, bam.SetChainID("testing"), bam.SetSnapshot(snapshotStore, snapshottypes.SnapshotOptions{KeepRecent: 2}))
	if withGenesis {
		return app, NewDefaultGenesisState(app.AppCodec())
	}
	return app, GenesisState{}
}

// SetupWithGenesisAccounts initializes a new SimApp with the provided genesis
// accounts and possible balances.
func SetupWithGenesisAccounts(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *BitsongApp {
	t.Helper()

	btsgApp, genesisState := setup(t, true)
	genesisState = genesisStateWithValSet(t, btsgApp, genesisState, valSet, genAccs, balances...)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	btsgApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
			ChainId:         SimAppChainID,
			Time:            time.Now().UTC(),
			InitialHeight:   1,
		},
	)

	// commit genesis changes
	btsgApp.Commit()
	btsgApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            SimAppChainID,
		Height:             btsgApp.LastBlockHeight() + 1,
		AppHash:            btsgApp.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
		Time:               time.Now().UTC(),
	}})

	return btsgApp
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, txCfg client.TxConfig, app *bam.BaseApp, header tmproto.Header, msgs []sdk.Msg,
	accNums, accSeqs []uint64, expSimPass, expPass bool, priv ...cryptotypes.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {
	tx, err := simtestutil.GenSignedMockTx(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		txCfg,
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
		simtestutil.DefaultGenTxGas,
		SimAppChainID,
		accNums,
		accSeqs,
		priv...,
	)
	require.NoError(t, err)
	txBytes, err := txCfg.TxEncoder()(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	_, res, err := app.Simulate(txBytes)

	if expSimPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	gInfo, res, err := app.SimDeliver(txCfg.TxEncoder(), tx)

	if expPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return gInfo, res, err
}

func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	return simtestutil.CreateTestPubKeys(numPubKeys)
}

func CheckBalance(t *testing.T, app *BitsongApp, addr sdk.AccAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, balances.IsEqual(app.AppKeepers.BankKeeper.GetAllBalances(ctxCheck, addr)))
}

func genesisStateWithValSet(t *testing.T,
	app *BitsongApp, genesisState GenesisState,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) GenesisState {
	codec := app.AppCodec()

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), math.LegacyOneDec()))

	}

	defaultStParams := stakingtypes.DefaultParams()
	stParams := stakingtypes.NewParams(
		defaultStParams.UnbondingTime,
		defaultStParams.MaxValidators,
		defaultStParams.MaxEntries,
		defaultStParams.HistoricalEntries,
		appparams.DefaultBondDenom,
		defaultStParams.MinCommissionRate, // 5%
	)

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(appparams.DefaultBondDenom, bondAmt.MulRaw(int64(len(valSet.Validators))))},
	})

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)
	// println("genesisStateWithValSet bankState:", string(genesisState[banktypes.ModuleName]))

	return genesisState
}
