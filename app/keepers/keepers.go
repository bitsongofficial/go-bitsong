package keepers

import (
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

	cadancekeeper "github.com/bitsongofficial/go-bitsong/x/cadance/keeper"
	cadancetypes "github.com/bitsongofficial/go-bitsong/x/cadance/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	icq "github.com/cosmos/ibc-apps/modules/async-icq/v8"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v8/types"
	ibcwasmkeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/keeper"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

var (
	wasmCapabilities = []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
		"cosmwasm_1_3",
		"cosmwasm_1_4",
		"cosmwasm_2_0",
		"cosmwasm_2_1",
		"bitsong",
	}
)

// module account permissions
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:            {authtypes.Burner},
	ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	fantokentypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	wasmtypes.ModuleName:           {authtypes.Burner},
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
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper
	ScopedICQKeeper      capabilitykeeper.ScopedKeeper

	// keepers
	AccountKeeper        *authkeeper.AccountKeeper
	BankKeeper           bankkeeper.Keeper
	AuthzKeeper          authzkeeper.Keeper
	StakingKeeper        *stakingkeeper.Keeper
	DistrKeeper          distrkeeper.Keeper
	SlashingKeeper       slashingkeeper.Keeper
	MintKeeper           mintkeeper.Keeper
	IBCKeeper            *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCHooksKeeper       *ibchookskeeper.Keeper
	TransferKeeper       ibctransferkeeper.Keeper
	ICQKeeper            *icqkeeper.Keeper
	EvidenceKeeper       evidencekeeper.Keeper
	FanTokenKeeper       fantokenkeeper.Keeper
	WasmKeeper           wasmkeeper.Keeper
	CadanceKeeper        cadancekeeper.Keeper
	IBCFeeKeeper         ibcfeekeeper.Keeper
	IBCWasmClientKeeper  *ibcwasmkeeper.Keeper
	FeeGrantKeeper       feegrantkeeper.Keeper
	GovKeeper            govkeeper.Keeper
	ContractKeeper       *wasmkeeper.PermissionedKeeper
	SmartAccountKeeper   *smartaccountkeeper.Keeper
	AuthenticatorManager *authenticator.AuthenticatorManager

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
	dataDir string,
	wasmDir string,
	wasmConfig wasmtypes.WasmConfig,
	ibcWasmConfig ibcwasmtypes.WasmConfig,
) AppKeepers {
	appKeepers := AppKeepers{}
	Bech32Prefix := "bitsong"
	// TODO: bitsong specific btsgPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()
	keys := appKeepers.GetKVStoreKey()
	tkeys := appKeepers.GetTransientStoreKey()

	appKeepers.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]),
		govModAddress,
		runtime.EventService{},
	)
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper.ParamsStore)

	// grant capabilities for the ibc and ibc-transfer modules
	// & add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	appKeepers.ScopedICQKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icqtypes.ModuleName)

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
		appCodec, runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]), authtypes.ProtoBaseAccount, maccPerms, addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()), Bech32Prefix,
		govModAddress,
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

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]), appKeepers.AccountKeeper,
	)
	smartAccountKeeper := smartaccountkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[smartaccounttypes.StoreKey],
		govModuleAddr,
		appKeepers.GetSubspace(smartaccounttypes.ModuleName),
		appKeepers.AuthenticatorManager,
		appKeepers.FeeGrantKeeper,
	)
	appKeepers.SmartAccountKeeper = &smartAccountKeeper

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]), appCodec, bApp.MsgServiceRouter(), appKeepers.AccountKeeper,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper, govModAddress,
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		stakingKeeper, authtypes.FeeCollectorName, govModAddress,
	)

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, cdc, runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]), stakingKeeper, govModAddress,
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[minttypes.StoreKey]), stakingKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper, authtypes.FeeCollectorName, govModAddress,
	)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistrKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)

	appKeepers.StakingKeeper = stakingKeeper

	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibcexported.StoreKey], appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper, appKeepers.UpgradeKeeper, appKeepers.ScopedIBCKeeper, govModAddress,
	)

	// Configure the ibchooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		appKeepers.keys[ibchookstypes.StoreKey],
	)
	appKeepers.IBCHooksKeeper = &hooksKeeper

	ibcWasmClientKeeper := ibcwasmkeeper.NewKeeperWithConfig(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[ibcwasmtypes.StoreKey]),
		appKeepers.IBCKeeper.ClientKeeper,
		govModAddress,
		ibcWasmConfig,
		bApp.GRPCQueryRouter(),
	)
	appKeepers.IBCWasmClientKeeper = &ibcWasmClientKeeper

	// Setup the ICS4Wrapper used by the hooks middleware
	wasmHooks := ibchooks.NewWasmHooks(appKeepers.IBCHooksKeeper, &appKeepers.WasmKeeper, Bech32Prefix)
	appKeepers.Ics20WasmHooks = &wasmHooks
	appKeepers.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		keys[packetforwardtypes.StoreKey],
		nil, // transfer keeper starts nil, gets set below with SetTransferKeeper
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.BankKeeper,
		// The ICS4Wrapper is replaced by the HooksICS4Wrapper instead of the channel so that sending can be overridden by the middleware
		appKeepers.HooksICS4Wrapper,
		govModAddress,
	)

	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.PacketForwardKeeper,
		appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper, govModAddress,
	)

	appKeepers.PacketForwardKeeper.SetTransferKeeper(appKeepers.TransferKeeper)

	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	const middlewareTimeoutRetry = 0
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, &appKeepers.HooksICS4Wrapper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		appKeepers.PacketForwardKeeper,
		middlewareTimeoutRetry,
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
	)

	// forwardTimeout time.Duration,
	// ICQ Keeper
	icqKeeper := icqkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icqtypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedICQKeeper,
		bApp.GRPCQueryRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.ICQKeeper = &icqKeeper
	// Create Async ICQ module
	icqModule := icq.NewIBCModule(*appKeepers.ICQKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(icqtypes.ModuleName, icqModule)
	// Note: sealing occurs after new wasmkeeper

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]), appKeepers.StakingKeeper, appKeepers.SlashingKeeper, addresscodec.NewBech32Codec(sdk.Bech32PrefixAccAddr),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	appKeepers.FanTokenKeeper = fantokenkeeper.NewKeeper(
		appCodec,
		keys[fantokentypes.StoreKey],
		appKeepers.GetSubspace(fantokentypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistrKeeper,
		BlockedAddrs(),
	)

	// Stargate Queries
	acceptedStargateQueries := wasmkeeper.AcceptedQueries{
		// ibc
		"/ibc.core.client.v1.Query/ClientState":    &ibcclienttypes.QueryClientStateResponse{},
		"/ibc.core.client.v1.Query/ConsensusState": &ibcclienttypes.QueryConsensusStateResponse{},
		"/ibc.core.connection.v1.Query/Connection": &ibcconnectiontypes.QueryConnectionResponse{},

		// distribution
		"/cosmos.distribution.v1beta1.Query/DelegationRewards": &distrtypes.QueryDelegationRewardsResponse{},

		// staking
		"/cosmos.staking.v1beta1.Query/Delegation":          &stakingtypes.QueryDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Redelegations":       &stakingtypes.QueryRedelegationsResponse{},
		"/cosmos.staking.v1beta1.Query/UnbondingDelegation": &stakingtypes.QueryUnbondingDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Validator":           &stakingtypes.QueryValidatorResponse{},
		"/cosmos.staking.v1beta1.Query/Params":              &stakingtypes.QueryParamsResponse{},
		"/cosmos.staking.v1beta1.Query/Pool":                &stakingtypes.QueryPoolResponse{},

		// fantoken
		"/bitsong.fantoken.v1beta1.Query/Params":    &fantokentypes.QueryParamsResponse{},
		"/bitsong.fantoken.v1beta1.Query/FanToken":  &fantokentypes.QueryFanTokenResponse{},
		"/bitsong.fantoken.v1beta1.Query/FanTokens": &fantokentypes.QueryFanTokensResponse{},
	}

	querierOpts := wasmkeeper.WithQueryPlugins(
		&wasmkeeper.QueryPlugins{Stargate: wasmkeeper.AcceptListStargateQuerier(acceptedStargateQueries, bApp.GRPCQueryRouter(), appCodec)})
	wasmOpts = append(wasmOpts, querierOpts)

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[wasmtypes.StoreKey]),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		wasmCapabilities,
		govModAddress,
		wasmOpts...,
	)

	// set the contract keeper for ICS4Wrappers (IBC Middleware)
	appKeepers.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(&appKeepers.WasmKeeper)
	appKeepers.Ics20WasmHooks.ContractKeeper = &appKeepers.WasmKeeper

	appKeepers.CadanceKeeper = cadancekeeper.NewKeeper(
		appKeepers.keys[cadancetypes.StoreKey],
		appCodec,
		appKeepers.WasmKeeper,
		appKeepers.ContractKeeper,
		govModAddress,
	)

	// register CosmWasm authenticator
	appKeepers.AuthenticatorManager.RegisterAuthenticator(
		authenticator.NewCosmwasmAuthenticator(appKeepers.ContractKeeper, appKeepers.AccountKeeper, appCodec))

	// wire wasm to IBC
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.ChannelKeeper))

	// Seal the router
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	// register the proposal types
	govRouter := govtypesv1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypesv1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(fantokentypes.RouterKey, fantoken.NewProposalHandler(appKeepers.FanTokenKeeper))

	govConfig := govtypes.DefaultConfig()
	govConfig.MaxMetadataLen = 10200
	appKeepers.GovKeeper = *govkeeper.NewKeeper(
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
	paramsKeeper.Subspace(icqtypes.ModuleName)
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
