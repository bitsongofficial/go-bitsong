package app

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"

	errorsmod "cosmossdk.io/errors"
	smartaccount "github.com/bitsongofficial/go-bitsong/x/smart-account"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"
	ibc_hooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10"
	ibcwlc "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10"
	ibcwlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

	ibcwasm "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"github.com/cosmos/cosmos-sdk/baseapp"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/protocolpool"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	wasmlckeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/keeper"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"

	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/codec/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/bytes"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	"github.com/bitsongofficial/go-bitsong/docs"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"

	v024 "github.com/bitsongofficial/go-bitsong/app/upgrades/v024"
	// unnamed import of statik for swagger UI support
	// _ "github.com/bitsongofficial/go-bitsong/swagger/statik"
	// performance profiling (pprof)
	// _ "net/http/pprof"
)

const appName = "BitsongApp"

// We pull these out so we can set them with LDFLAGS in the Makefile
var (
	NodeDir       = ".bitsongd"
	Bech32Prefix  = "bitsong"
	EmptyWasmOpts []wasmkeeper.Option
	// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "true"
	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	EnableSpecificProposals = ""

	Upgrades = []upgrades.Upgrade{
		// v010.Upgrade, v011.Upgrade, v013.Upgrade, v014.Upgrade,
		// v015.Upgrade, v016.Upgrade, v018.Upgrade, v020.Upgrade,
		// v021.Upgrade, v022.Upgrade, v023.Upgrade,
		v024.Upgrade,
	}
)

func init() {
	SetAddressPrefixes()
}

var (
	// DefaultNodeHome default home directories for Bitsongd
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

// SetAddressPrefixes builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts, validators, and consensus nodes and verifies that addreeses have correct format.
func SetAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	// This is copied from the cosmos sdk v0.43.0-beta1
	// source: https://github.com/cosmos/cosmos-sdk/blob/v0.43.0-beta1/types/address.go#L141
	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}

		if len(bytes) > address.MaxAddrLen {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}

		// TODO: Do we want to allow addresses of lengths other than 20 and 32 bytes?
		if len(bytes) != 20 && len(bytes) != 32 {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}

		return nil
	})
}

func GetWasmOpts(appOpts servertypes.AppOptions) []wasm.Option {
	var wasmOpts []wasmkeeper.Option
	if cast.ToBool(appOpts.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	// default wasm gas configuration.
	wasmOpts = append(wasmOpts, wasmkeeper.WithGasRegister(wasmtypes.NewWasmGasRegister(wasmtypes.DefaultGasRegisterConfig())))

	return wasmOpts
}

var (
	_ CosmosApp               = (*BitsongApp)(nil)
	_ servertypes.Application = (*BitsongApp)(nil)
)

// BitsongApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type BitsongApp struct {
	*baseapp.BaseApp
	keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// the module manager
	mm *module.Manager

	// simulation manager
	sm           *module.SimulationManager
	configurator module.Configurator
	homePath     string
}

// init sets DefaultNodeHome to default bitsongd install location.
func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".bitsongd")
}

// overrideWasmVariables overrides the wasm variables to:
//   - allow for larger wasm files
func overrideWasmVariables() {
	// Override Wasm size limitation from WASMD.
	wasmtypes.MaxWasmSize = 7 * 1024 * 1024 // 7mb wasm blob
	wasmtypes.MaxProposalWasmSize = wasmtypes.MaxWasmSize
}

// NewBitsongApp returns a reference to an initialized BitsongApp.
func NewBitsongApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	homePath string,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *BitsongApp {
	encodingConfig := MakeEncodingConfig()
	overrideWasmVariables()

	appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(appName, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	app := &BitsongApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		tkeys:             storetypes.NewTransientStoreKeys(paramstypes.TStoreKey),
	}

	app.homePath = homePath
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadNodeConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	ibcWasmConfig := wasmlctypes.WasmConfig{
		DataDir:               filepath.Join(homePath, "ibc_08-wasm"),
		SupportedCapabilities: keepers.ExtendedBuiltInCapabilities(),
		ContractDebugMode:     false,
	}

	// Setup keepers
	appKeepers := keepers.NewAppKeepers(
		appCodec,
		encodingConfig,
		bApp,
		legacyAmino,
		keepers.GetMaccPerms(),
		appOpts,
		wasmOpts,
		wasmDir,
		wasmConfig,
		ibcWasmConfig,
	)
	// add all light clients to app keepers
	clientKeeper := appKeepers.IBCKeeper.ClientKeeper
	storeProvider := appKeepers.IBCKeeper.ClientKeeper.GetStoreProvider()
	tmLightClientModule := ibctm.NewLightClientModule(appCodec, storeProvider)
	ibcWasmLightClientModule := ibcwlc.NewLightClientModule(*appKeepers.IBCWasmClientKeeper, storeProvider)
	clientKeeper.AddRoute(ibctm.ModuleName, &tmLightClientModule)
	clientKeeper.AddRoute(ibcwlctypes.ModuleName, ibcWasmLightClientModule)

	app.AppKeepers = appKeepers

	app.keys = app.GetKVStoreKey()

	// cosmos-sdk@v0.50: textual signature for ledger devices
	enabledSignModes := append(authtx.DefaultSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)
	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfigWithTextual, err := authtx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create textual tx config: %w", err))
	}
	app.txConfig = txConfigWithTextual

	// load state streaming if enabled
	if err := app.RegisterStreamingServices(appOpts, app.keys); err != nil {
		panic(err)
	}

	// upgrade info
	app.setupUpgradeStoreLoaders()

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(genutil.NewAppModule(
		app.AccountKeeper, app.StakingKeeper, app.BaseApp,
		encodingConfig.TxConfig,
	),
		auth.NewAppModule(appCodec, *app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(*app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(appCodec, *app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, addresscodec.NewBech32Codec(Bech32PrefixAccAddr)),
		wasm.NewAppModule(appCodec, app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		fantoken.NewAppModule(appCodec, app.FanTokenKeeper, app.BankKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibcwasm.NewAppModule(*app.IBCWasmClientKeeper),
		ibctm.NewAppModule(tmLightClientModule),
		params.NewAppModule(app.ParamsKeeper),
		consensus.NewAppModule(appCodec, *app.ConsensusParamsKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		transfer.NewAppModule(*app.TransferKeeper),
		ibc_hooks.NewAppModule(*app.AccountKeeper),
		protocolpool.NewAppModule(*app.ProtocolPoolKeeper, app.AccountKeeper, app.BankKeeper),
		smartaccount.NewAppModule(appCodec, *app.SmartAccountKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)))

	// NOTE: upgrade module is prioritized in preblock
	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
		authtypes.ModuleName,
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.mm.SetOrderEndBlockers(orderEndBlockers()...)
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)

	app.mm.RegisterInvariants(app.CrisisKeeper)

	// upgrade handlers
	app.configurator = module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	// register upgrade
	app.setupUpgradeHandlers(app.configurator)

	app.sm = module.NewSimulationManager(simulationModules(app, encodingConfig, skipGenesisInvariants)...)

	app.sm.RegisterStoreDecoders()

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc := getReflectionService()
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.keys)
	app.MountTransientStores(app.GetTransientStoreKey())

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			SmartAccount:      app.SmartAccountKeeper,
			IBCKeeper:         app.IBCKeeper,
			TxCounterStoreKey: runtime.NewKVStoreService(app.GetKey(wasmtypes.StoreKey)),
			WasmConfig:        wasmConfig,
			Cdc:               appCodec,
			TxEncoder:         app.txConfig.TxEncoder(),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(anteHandler)
	app.SetPostHandler(NewPostHandler(appCodec, app.SmartAccountKeeper, app.AccountKeeper, encodingConfig.TxConfig.SignModeHandler()))
	app.SetEndBlocker(app.EndBlocker)
	app.SetPrecommiter(app.Precommitter)
	app.SetPrepareCheckStater(app.PrepareCheckStater)

	if manager := app.SnapshotManager(); manager != nil {
		err = manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
		//  takes care of persisting the external state from wasm code when snapshot is created
		err = manager.RegisterExtensions(
			wasmlckeeper.NewWasmSnapshotter(app.CommitMultiStore(), app.IBCWasmClientKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
		ctx := app.BaseApp.NewUncachedContext(true, cmtproto.Header{})
		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}

		if err := app.IBCWasmClientKeeper.InitializePinnedCodes(ctx); err != nil {
			tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
	}

	app.sm.RegisterStoreDecoders()

	return app
}

// Name returns the name of the App
func (app *BitsongApp) Name() string { return app.BaseApp.Name() }

func (app *BitsongApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// PreBlocker application updates before each begin block.
func (app *BitsongApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	// Set gas meter to the free gas meter.
	// This is because there is currently non-deterministic gas usage in the
	// pre-blocker, e.g. due to hydration of in-memory data structures.
	//
	// Note that we don't need to reset the gas meter after the pre-blocker
	// because Go is pass by value.
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	mm := app.ModuleManager()
	return mm.PreBlock(ctx)
}

// BeginBlocker application updates every begin block
func (app *BitsongApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *BitsongApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// Precommitter application updates before the commital of a block after all transactions have been delivered.
func (app *BitsongApp) Precommitter(ctx sdk.Context) {
	mm := app.ModuleManager()
	if err := mm.Precommit(ctx); err != nil {
		panic(err)
	}
}

func (app *BitsongApp) PrepareCheckStater(ctx sdk.Context) {
	mm := app.ModuleManager()
	if err := mm.PrepareCheckState(ctx); err != nil {
		panic(err)
	}
}

// InitChainer application update at chain initialization
func (app *BitsongApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *BitsongApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *BitsongApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range keepers.GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *BitsongApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns Bitsong's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *BitsongApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Bitsong's InterfaceRegistry
func (app *BitsongApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *BitsongApp) ModuleManager() module.Manager {
	return *app.mm
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *BitsongApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *BitsongApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *BitsongApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *BitsongApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all app module routes with the provided
// API server.
func (app *BitsongApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	module.NewBasicManagerFromManager(app.mm, nil).RegisterGRPCGatewayRoutes(
		clientCtx,
		apiSvr.GRPCGatewayRouter,
	)

	// Register new tendermint queries routes from grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr)
	}

}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *BitsongApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *BitsongApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService implements the Application.RegisterNodeService method.
func (app *BitsongApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.BaseApp.GRPCQueryRouter(), cfg)
}

// SimulationManager implements the SimulationApp interface
func (app *BitsongApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *BitsongApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			storeUpgrades := upgrade.StoreUpgrades
			app.SetStoreLoader(
				upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades),
			)
		}
	}
}

func (app *BitsongApp) setupUpgradeHandlers(cfg module.Configurator) {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.mm,
				cfg,
				app.BaseApp,
				&app.AppKeepers,
				app.homePath,
			),
		)
	}
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, apiSvr *api.Server) error {
	staticSubDir, err := fs.Sub(docs.Docs, "static")
	if err != nil {
		return err
	}

	staticServer := http.FileServer(http.FS(staticSubDir))

	// Handle /swag without trailing slash - redirect to /swag/
	apiSvr.Router.HandleFunc("/swag", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swag/", http.StatusMovedPermanently)
	})

	apiSvr.Router.PathPrefix("/swag/").Handler(http.StripPrefix("/swag/", staticServer))

	return nil
}

// we cache the reflectionService to save us time within tests.
var cachedReflectionService *runtimeservices.ReflectionService = nil

func getReflectionService() *runtimeservices.ReflectionService {
	if cachedReflectionService != nil {
		return cachedReflectionService
	}
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	cachedReflectionService = reflectionSvc
	return reflectionSvc
}

// source: https://github.com/osmosis-labs/osmosis/blob/7b1a78d397b632247fe83f51867f319adf3a858c/app/app.go#L786
// one-liner: cd ../bitsong-snapshots && bitsongd comet unsafe-reset-all && cp ~/.bitsongd/data/priv_validator_state.json ~/.bitsongd/priv_validator_state.json && lz4 -c -d <bitsong-snapshot>.tar.lz4 | tar -x -C $HOME/.bitsongd && cp ~/.bitsongd/priv_validator_state.json ~/.bitsongd/data/priv_validator_state.json && cd ../go-bitsong && make install && bitsongd in-place-testnet test1 bitsong1mt3wj088jvurp3vlh2yfar6vqrqp0llnsj8lar bitsongvaloper1qxw4fjged2xve8ez7nu779tm8ejw92rv0vcuqr
func InitBitsongAppForTestnet(app *BitsongApp, newValAddr bytes.HexBytes, newValPubKey crypto.PubKey, newOperatorAddress, upgradeToTrigger, brokenVals string) *BitsongApp { // newValsPower []testnetserver.ValidatorInfo

	ctx := app.BaseApp.NewUncachedContext(true, cmtproto.Header{})
	pubkey := &ed25519.PubKey{Key: newValPubKey.Bytes()}
	pubkeyAny, err := types.NewAnyWithValue(pubkey)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// STAKING
	brokeValAddr, err := sdk.ValAddressFromBech32(brokenVals)
	if err != nil {
		tmos.Exit(err.Error())
	}
	retainedValidator, err := app.StakingKeeper.GetValidator(ctx, brokeValAddr)
	if err != nil {
		tmos.Exit(err.Error())
	}
	fmt.Printf("retainedValidator: %v\n", retainedValidator)

	retainedValDels, err := app.StakingKeeper.GetValidatorDelegations(ctx, brokeValAddr)
	if err != nil {
		tmos.Exit(err.Error())
	}
	fmt.Printf("retainedValDels: %v\n", retainedValDels)

	// Create Validator struct for our new validator.
	_, bz, err := bech32.DecodeAndConvert(newOperatorAddress)
	if err != nil {
		tmos.Exit(err.Error())
	}
	bech32Addr, err := bech32.ConvertAndEncode("bitsongvaloper", bz)
	if err != nil {
		tmos.Exit(err.Error())
	}
	newVal := stakingtypes.Validator{
		OperatorAddress: bech32Addr,
		ConsensusPubkey: pubkeyAny,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(900000000000000),
		DelegatorShares: math.LegacyMustNewDecFromStr("10000000"),
		Description: stakingtypes.Description{
			Moniker: "Testnet Validator",
		},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.05"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
			},
		},
		MinSelfDelegation: math.OneInt(),
	}

	// Remove all validators from power store
	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, err := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		tmos.Exit(err.Error())
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Remove all valdiators from last validators store
	iterator, err = app.StakingKeeper.LastValidatorsIterator(ctx)
	if err != nil {
		tmos.Exit(err.Error())
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	//  TODO: retain validator from store
	// Remove all validators from validators store
	iterator = storetypes.KVStorePrefixIterator(stakingStore, stakingtypes.ValidatorsKey)
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Remove all validators from unbonding queue
	iterator = storetypes.KVStorePrefixIterator(stakingStore, stakingtypes.ValidatorQueueKey)
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Add our validator to power and last validators store
	err = app.StakingKeeper.SetValidator(ctx, newVal)
	if err != nil {
		tmos.Exit(err.Error())
	}
	// Add retainedValidator to power and last validators store
	err = app.StakingKeeper.SetValidator(ctx, retainedValidator)
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, retainedValidator)
	if err != nil {
		tmos.Exit(err.Error())
	}

	err = app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	if err != nil {
		tmos.Exit(err.Error())
	}

	valAddr, err := sdk.ValAddressFromBech32(newVal.GetOperator())
	if err != nil {
		tmos.Exit(err.Error())
	}

	err = app.StakingKeeper.SetLastValidatorPower(ctx, valAddr, 1)
	if err != nil {
		tmos.Exit(err.Error())
	}

	if err := app.StakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddr); err != nil {
		panic(err)
	}

	// Initialize records for this validator across all distribution stores
	valAddr, err = sdk.ValAddressFromBech32(newVal.GetOperator())
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddr, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, valAddr, distrtypes.InitialValidatorAccumulatedCommission())
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.DistrKeeper.SetValidatorOutstandingRewards(ctx, valAddr, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
	if err != nil {
		tmos.Exit(err.Error())
	}

	// SLASHING
	// Set validator signing info for our new validator.
	newConsAddr := sdk.ConsAddress(newValAddr.Bytes())
	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	err = app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)
	if err != nil {
		tmos.Exit(err.Error())
	}

	newExpeditedVotingPeriod := time.Minute
	newVotingPeriod := time.Minute * 2

	govParams, err := app.GovKeeper.Params.Get(ctx)
	if err != nil {
		tmos.Exit(err.Error())
	}

	govParams.ExpeditedVotingPeriod = &newExpeditedVotingPeriod
	govParams.VotingPeriod = &newVotingPeriod
	govParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin("ubtsg", 100000000))
	govParams.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewInt64Coin("ubtsg", 150000000))
	err = app.GovKeeper.Params.Set(ctx, govParams)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// BANK
	//

	// Fund edgenet faucet

	// UPGRADE
	//

	if upgradeToTrigger != "" {
		upgradePlan := upgradetypes.Plan{
			Name:   upgradeToTrigger,
			Height: app.LastBlockHeight() + 10,
		}
		err = app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan)
		if err != nil {
			panic(err)
		}
	}
	return app
}
