package app

import (
	"io"
	"os"

	rewardTypes "github.com/bitsongofficial/go-bitsong/x/reward/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	cmint "github.com/cosmos/cosmos-sdk/x/mint"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	desmosibc "github.com/bitsongofficial/go-bitsong/x/ibc_desmos"
	"github.com/bitsongofficial/go-bitsong/x/mint"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/bitsongofficial/go-bitsong/x/track"

	"github.com/bitsongofficial/go-bitsong/x/reward"
)

const appName = "GoBitsong"

var (
	// DefaultCLIHome represents the default home directory for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.bitsongcli")

	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.bitsongd")

	// ModuleBasics is in charge of setting up basic module elements
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		supply.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		cmint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler,
			distr.ProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},

		reward.AppModuleBasic{},
		track.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		cmint.ModuleName:          {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		track.ModuleName:          {supply.Burner},
		reward.ModuleName:         nil,
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName: true,
	}
)

var _ simapp.App = (*GoBitsong)(nil)

// Extended ABCI application
type GoBitsong struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// sdk keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	cmintKeeper    cmint.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	ibcKeeper      ibc.Keeper

	// bitsong keepers
	rewardKeeper    reward.Keeper
	trackKeeper     track.Keeper
	desmosIBCKeeper desmosibc.Keeper

	// the module manager
	mm *module.Manager
}

// NewBitsongApp returns a reference to an initialized GoBitsong.
func NewBitsongApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *GoBitsong {
	// First define the top level codec that will be shared by the different modules
	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	cdc := codecstd.MakeCodec(ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey, bank.StoreKey,
		supply.StoreKey, cmint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, ibc.StoreKey,

		reward.StoreKey, track.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &GoBitsong{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// Init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[cmint.ModuleName] = app.paramsKeeper.Subspace(cmint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	app.subspaces[reward.ModuleName] = app.paramsKeeper.Subspace(reward.DefaultParamspace)
	app.subspaces[track.ModuleName] = app.paramsKeeper.Subspace(track.DefaultParamspace)

	// Add keepers
	app.accountKeeper = auth.NewAccountKeeper(
		appCodec, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		appCodec, keys[bank.ModuleName], app.accountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	app.supplyKeeper = supply.NewKeeper(
		appCodec, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		appCodec, keys[staking.StoreKey], app.bankKeeper, app.supplyKeeper, app.subspaces[staking.ModuleName],
	)
	app.cmintKeeper = cmint.NewKeeper(
		appCodec, keys[cmint.StoreKey], app.subspaces[cmint.ModuleName], &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName,
	)
	app.distrKeeper = distr.NewKeeper(
		appCodec, keys[distr.StoreKey], app.subspaces[distr.ModuleName], app.bankKeeper, &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashing.NewKeeper(
		appCodec, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.crisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))
	app.govKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], app.subspaces[gov.ModuleName],
		app.supplyKeeper, &stakingKeeper, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	app.ibcKeeper = ibc.NewKeeper(app.cdc, keys[ibc.StoreKey], stakingKeeper)

	// Custom modules
	app.rewardKeeper = reward.NewKeeper(
		app.cdc, keys[rewardTypes.StoreKey], app.subspaces[reward.ModuleName],
		app.supplyKeeper, app.trackKeeper, app.bankKeeper,
	)
	app.trackKeeper = track.NewKeeper(
		app.cdc, app.keys[track.ModuleName], "track", // TODO: Change this
		app.stakingKeeper, app.accountKeeper, app.supplyKeeper,
		app.subspaces[track.ModuleName],
	)
	app.mintKeeper = mint.NewKeeper(app.rewardKeeper, app.supplyKeeper)
	app.desmosIBCKeeper = desmosibc.NewKeeper(app.cdc, app.ibcKeeper.ChannelKeeper)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper, app.supplyKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.bankKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.bankKeeper, app.supplyKeeper),
		mint.NewAppModule(cmint.NewAppModule(app.cmintKeeper, app.supplyKeeper), app.cmintKeeper, app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.bankKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.bankKeeper, app.supplyKeeper),
		ibc.NewAppModule(app.ibcKeeper),

		reward.NewAppModule(app.rewardKeeper, app.supplyKeeper, app.bankKeeper),
		track.NewAppModule(app.trackKeeper, app.bankKeeper),
		desmosibc.NewAppModule(app.desmosIBCKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(cmint.ModuleName, distr.ModuleName, slashing.ModuleName, staking.ModuleName)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, cmint.ModuleName, supply.ModuleName,
		crisis.ModuleName, genutil.ModuleName,

		track.ModuleName, reward.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(app.accountKeeper, app.supplyKeeper, app.ibcKeeper, ante.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	authvesting.RegisterCodec(cdc)

	return cdc
}

// Name returns the name of the App
func (app *GoBitsong) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *GoBitsong) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *GoBitsong) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *GoBitsong) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	res := app.mm.InitGenesis(ctx, app.cdc, genesisState)

	// Set Historical infos in InitChain to ignore genesis params
	stakingParams := staking.DefaultParams()
	stakingParams.HistoricalEntries = 1000
	app.stakingKeeper.SetParams(ctx, stakingParams)

	return res
}

// LoadHeight loads a particular height
func (app *GoBitsong) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *GoBitsong) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *GoBitsong) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *GoBitsong) SimulationManager() *module.SimulationManager {
	// TODO
	return nil
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *GoBitsong) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[supply.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}
