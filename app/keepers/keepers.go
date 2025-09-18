package keepers

import (
	"fmt"

	"github.com/spf13/cast"

	"cosmossdk.io/x/feegrant"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/params"

	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"
	fantokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/bitsongofficial/go-bitsong/x/smart-account/authenticator"
	smartaccountkeeper "github.com/bitsongofficial/go-bitsong/x/smart-account/keeper"
	smartaccounttypes "github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	protocolpoolkeeper "github.com/cosmos/cosmos-sdk/x/protocolpool/keeper"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	wasmvm "github.com/CosmWasm/wasmvm/v3"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10/types"

	"github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v10/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	ibcwlc "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10"
	ibcwlckeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/keeper"
	ibcwlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
)

func ExtendedBuiltInCapabilities() []string {
	return append(wasmkeeper.BuiltInCapabilities(), "bitsong", "cosmwasm_3_0")
}

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:                  nil,
	distrtypes.ModuleName:                       nil,
	minttypes.ModuleName:                        {authtypes.Minter},
	stakingtypes.BondedPoolName:                 {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:              {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                         {authtypes.Burner},
	ibctransfertypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	fantokentypes.ModuleName:                    {authtypes.Minter, authtypes.Burner},
	wasmtypes.ModuleName:                        {authtypes.Burner},
	protocolpooltypes.ModuleName:                nil,
	protocolpooltypes.ProtocolPoolEscrowAccount: nil,
}

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	ParamsKeeper          paramskeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ConsensusParamsKeeper *consensusparamkeeper.Keeper

	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	// keepers
	AccountKeeper        *authkeeper.AccountKeeper
	BankKeeper           bankkeeper.Keeper
	AuthzKeeper          *authzkeeper.Keeper
	StakingKeeper        *stakingkeeper.Keeper
	DistrKeeper          *distrkeeper.Keeper
	SlashingKeeper       *slashingkeeper.Keeper
	MintKeeper           *mintkeeper.Keeper
	IBCKeeper            *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCHooksKeeper       *ibchookskeeper.Keeper
	TransferKeeper       *ibctransferkeeper.Keeper
	EvidenceKeeper       *evidencekeeper.Keeper
	FanTokenKeeper       *fantokenkeeper.Keeper
	WasmKeeper           *wasmkeeper.Keeper
	IBCWasmClientKeeper  *ibcwlckeeper.Keeper
	FeeGrantKeeper       *feegrantkeeper.Keeper
	GovKeeper            *govkeeper.Keeper
	ContractKeeper       *wasmkeeper.PermissionedKeeper
	SmartAccountKeeper   *smartaccountkeeper.Keeper
	AuthenticatorManager *authenticator.AuthenticatorManager
	ProtocolPoolKeeper   *protocolpoolkeeper.Keeper

	// Middleware wrapper
	Ics20WasmHooks      *ibchooks.WasmHooks
	HooksICS4Wrapper    ibchooks.ICS4Middleware
	PacketForwardKeeper *packetforwardkeeper.Keeper
}

func NewAppKeepers(
	appCodec codec.Codec,
	encodingConfig appparams.EncodingConfig,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	wasmDir string,
	wasmConfig wasmtypes.NodeConfig,
	ibcWasmConfig ibcwlctypes.WasmConfig,
) AppKeepers {
	appKeepers := AppKeepers{}
	Bech32Prefix := "bitsong"
	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	wasmCapabilities := ExtendedBuiltInCapabilities()
	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()
	keys := appKeepers.GetKVStoreKey()
	tkeys := appKeepers.GetTransientStoreKey()

	appKeepers.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	consensusParamsKeeper := consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		govModAddress,
		runtime.EventService{},
	)
	appKeepers.ConsensusParamsKeeper = &consensusParamsKeeper
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper.ParamsStore)

	// grant capabilities for the ibc and ibc-transfer modules
	// & add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[crisistypes.StoreKey]),
		invCheckPeriod, appKeepers.BankKeeper, authtypes.FeeCollectorName,
		govModAddress,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
	)

	// get skipUpgradeHeights from the app options
	// TODO: update to get from confix
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]), appCodec, homePath, bApp, govModAddress,
	)

	// add keepers
	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		Bech32Prefix,
		govModAddress,
		authkeeper.WithUnorderedTransactions(false),
	)
	appKeepers.AccountKeeper = &accountKeeper

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]), appKeepers.AccountKeeper, BlockedAddrs(),
		govModAddress, bApp.Logger(),
	)

	// Initialize authenticators
	appKeepers.AuthenticatorManager = authenticator.NewAuthenticatorManager()
	appKeepers.AuthenticatorManager.InitializeAuthenticators([]authenticator.Authenticator{
		authenticator.NewSignatureVerification(appKeepers.AccountKeeper),
		authenticator.NewMessageFilter(encodingConfig),
		authenticator.NewAllOf(appKeepers.AuthenticatorManager),
		authenticator.NewAnyOf(appKeepers.AuthenticatorManager),
		authenticator.NewPartitionedAnyOf(appKeepers.AuthenticatorManager),
		authenticator.NewPartitionedAllOf(appKeepers.AuthenticatorManager),
	})
	govModuleAddr := appKeepers.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	feegrantKeeper := feegrantkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]), appKeepers.AccountKeeper,
	)
	appKeepers.FeeGrantKeeper = &feegrantKeeper

	smartAccountKeeper := smartaccountkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[smartaccounttypes.StoreKey],
		govModuleAddr,
		appKeepers.GetSubspace(smartaccounttypes.ModuleName),
		appKeepers.AuthenticatorManager,
		*appKeepers.FeeGrantKeeper,
	)
	appKeepers.SmartAccountKeeper = &smartAccountKeeper

	authzKeeper := authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]), appCodec, bApp.MsgServiceRouter(), appKeepers.AccountKeeper,
	)
	appKeepers.AuthzKeeper = &authzKeeper

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper, govModAddress,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	// TODO: implement custom mint function to split inflation rewards to incentivize cdn:
	protopoolKeeper := protocolpoolkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[protocolpooltypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.ProtocolPoolKeeper = &protopoolKeeper

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		stakingKeeper, authtypes.FeeCollectorName, govModAddress, // distrkeeper.WithExternalCommunityPool(appKeepers.ProtocolPoolKeeper),
	)
	appKeepers.DistrKeeper = &distrKeeper

	slashKeeper := slashingkeeper.NewKeeper(
		appCodec, cdc, runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]), stakingKeeper, govModAddress,
	)
	appKeepers.SlashingKeeper = &slashKeeper

	mintKeeper := mintkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[minttypes.StoreKey]), stakingKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper, authtypes.FeeCollectorName, govModAddress,
		// mintkeeper.WithMintFn(myCustomMintFunc), // Use custom minting function: https://github.com/cosmos/cosmos-sdk/blob/v0.53.0/UPGRADING.md?plain=1#L192
	)
	appKeepers.MintKeeper = &mintKeeper

	// register the staking hookshttps://github.com/cosmos/cosmos-sdk/blob/v0.53.0/UPGRADING.md?plain=1#L192
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistrKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)

	appKeepers.StakingKeeper = stakingKeeper

	fantokenKeeper := fantokenkeeper.NewKeeper(
		appCodec,
		keys[fantokentypes.StoreKey],
		appKeepers.GetSubspace(fantokentypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		BlockedAddrs(),
	)
	appKeepers.FanTokenKeeper = &fantokenKeeper

	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[ibcexported.StoreKey]), appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.UpgradeKeeper, govModAddress,
	)

	// Configure the ibchooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		appKeepers.keys[ibchookstypes.StoreKey],
	)
	appKeepers.IBCHooksKeeper = &hooksKeeper

	// Setup the ICS4Wrapper used by the hooks middleware
	wasmHooks := ibchooks.NewWasmHooks(appKeepers.IBCHooksKeeper, appKeepers.WasmKeeper, Bech32Prefix)
	appKeepers.Ics20WasmHooks = &wasmHooks
	appKeepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[packetforwardtypes.StoreKey]),
		nil, // transfer keeper starts nil, gets set below with SetTransferKeeper
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.BankKeeper,
		// The ICS4Wrapper is replaced by the HooksICS4Wrapper instead of the channel so that sending can be overridden by the middleware
		appKeepers.HooksICS4Wrapper,
		govModAddress,
	)

	transferKeeper := ibctransferkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[ibctransfertypes.StoreKey]),
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PacketForwardKeeper,
		appKeepers.IBCKeeper.ChannelKeeper, bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper, appKeepers.BankKeeper,
		govModAddress,
	)
	appKeepers.TransferKeeper = &transferKeeper

	appKeepers.PacketForwardKeeper.SetTransferKeeper(appKeepers.TransferKeeper)

	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	const middlewareTimeoutRetry = 0
	transferStack = transfer.NewIBCModule(*appKeepers.TransferKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, &appKeepers.HooksICS4Wrapper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		appKeepers.PacketForwardKeeper,
		middlewareTimeoutRetry,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]), appKeepers.StakingKeeper, appKeepers.SlashingKeeper, addresscodec.NewBech32Codec(sdk.Bech32PrefixAccAddr),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = evidenceKeeper

	acceptedStargateQueries := AcceptedQueries()
	querierOpts := wasmkeeper.WithQueryPlugins(
		&wasmkeeper.QueryPlugins{Stargate: wasmkeeper.AcceptListStargateQuerier(acceptedStargateQueries, bApp.GRPCQueryRouter(), appCodec)})
	wasmOpts = append(wasmOpts, querierOpts)

	// create wasmvm to use for both x/wasm and wasm-light-client
	wasmVm, err := wasmvm.NewVM(wasmDir, wasmCapabilities, 32, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	if err != nil {
		panic(fmt.Sprintf("failed to create bitsong wasm vm: %s", err))
	}

	lcWasmer, err := wasmvm.NewVM(ibcWasmConfig.DataDir, wasmCapabilities, 32, ibcWasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	if err != nil {
		panic(fmt.Sprintf("failed to create bitsong wasm vm for 08-wasm: %s", err))
	}

	appWasmKeeper := wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[wasmtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(*appKeepers.DistrKeeper),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeperV2,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmtypes.VMConfig{},
		wasmCapabilities,
		govModAddress,
		append(wasmOpts, wasmkeeper.WithWasmEngine(wasmVm))...,
	)
	appKeepers.WasmKeeper = &appWasmKeeper

	ibcWasmClientKeeper := ibcwlckeeper.NewKeeperWithVM(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[ibcwlctypes.StoreKey]),
		appKeepers.IBCKeeper.ClientKeeper,
		govModAddress,
		lcWasmer,
		bApp.GRPCQueryRouter(),
	)
	appKeepers.IBCWasmClientKeeper = &ibcWasmClientKeeper

	// set the contract keeper for ICS4Wrappers (IBC Middleware)
	appKeepers.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(appKeepers.WasmKeeper)
	appKeepers.Ics20WasmHooks.ContractKeeper = appKeepers.WasmKeeper

	// register CosmWasm authenticator
	appKeepers.AuthenticatorManager.RegisterAuthenticator(
		authenticator.NewCosmwasmAuthenticator(appKeepers.ContractKeeper, appKeepers.AccountKeeper, appCodec))

	// wire wasm to IBC

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.TransferKeeper, appKeepers.IBCKeeper.ChannelKeeper))
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	clientKeeper := appKeepers.IBCKeeper.ClientKeeper
	storeProvider := appKeepers.IBCKeeper.ClientKeeper.GetStoreProvider()

	// Add tendermint & ibcWasm light client routes
	tmLightClientModule := ibctm.NewLightClientModule(appCodec, storeProvider)
	ibcWasmLightClientModule := ibcwlc.NewLightClientModule(*appKeepers.IBCWasmClientKeeper, storeProvider)
	clientKeeper.AddRoute(ibctm.ModuleName, &tmLightClientModule)
	clientKeeper.AddRoute(ibcwlctypes.ModuleName, ibcWasmLightClientModule)

	// register the proposal types
	govRouter := govtypesv1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypesv1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(fantokentypes.RouterKey, fantoken.NewProposalHandler(appKeepers.FanTokenKeeper))

	govConfig := govtypes.DefaultConfig()
	govConfig.MaxMetadataLen = 10200
	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		stakingKeeper,
		appKeepers.DistrKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		govModAddress,
	)
	// Set legacy router for backwards compatibility with gov v1beta1
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	return appKeepers

}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	keytable := ibcclienttypes.ParamKeyTable()
	keytable.RegisterParamSet(&ibcconnectiontypes.Params{})

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keytable)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	paramsKeeper.Subspace(wasmtypes.ModuleName)
	paramsKeeper.Subspace(ibchookstypes.ModuleName)
	paramsKeeper.Subspace(smartaccounttypes.ModuleName).WithKeyTable(smartaccounttypes.ParamKeyTable())
	paramsKeeper.Subspace(fantokentypes.ModuleName)

	return paramsKeeper
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// BlockedAddrs returns all the app's module account addresses that are not allowed to receive tokens
func BlockedAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	// allow supplement pool amount to receive tokens
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// ModuleAccountAddrs provides a list of blocked module accounts from configuration in AppConfig
//
// Ported from WasmApp

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}
