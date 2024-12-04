package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	encparams "github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distr.AppModuleBasic{},
	gov.NewAppModuleBasic(getGovProposalHandlers()),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	consensus.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	packetforward.AppModuleBasic{},
	fantoken.AppModuleBasic{},
	wasm.AppModuleBasic{},
)

func appModules(
	app *BitsongApp,
	encodingConfig encparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler
	return []module.AppModule{

		genutil.NewAppModule(
			app.AppKeepers.AccountKeeper, app.AppKeepers.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		bank.NewAppModule(appCodec, app.AppKeepers.BankKeeper, app.AppKeepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.AppKeepers.CapabilityKeeper, false),
		gov.NewAppModule(appCodec, &app.AppKeepers.GovKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.AppKeepers.MintKeeper, app.AppKeepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.AppKeepers.SlashingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.AppKeepers.DistrKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.AppKeepers.UpgradeKeeper),
		evidence.NewAppModule(app.AppKeepers.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.FeeGrantKeeper, app.interfaceRegistry),
		fantoken.NewAppModule(appCodec, app.AppKeepers.FanTokenKeeper, app.AppKeepers.BankKeeper),
		authzmodule.NewAppModule(appCodec, app.AppKeepers.AuthzKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.AppKeepers.IBCKeeper),
		params.NewAppModule(app.AppKeepers.ParamsKeeper),
		transfer.NewAppModule(app.AppKeepers.TransferKeeper),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		packetforward.NewAppModule(app.AppKeepers.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		crisis.NewAppModule(app.AppKeepers.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
	}
}

func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName, capabilitytypes.ModuleName, crisistypes.ModuleName, govtypes.ModuleName,
		stakingtypes.ModuleName, ibctransfertypes.ModuleName, ibcexported.ModuleName, packetforwardtypes.ModuleName,
		authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		minttypes.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, authz.ModuleName, wasmtypes.ModuleName,
		feegrant.ModuleName, paramstypes.ModuleName, vestingtypes.ModuleName, fantokentypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName, ibctransfertypes.ModuleName, ibcexported.ModuleName,
		packetforwardtypes.ModuleName, feegrant.ModuleName, authz.ModuleName, capabilitytypes.ModuleName, authtypes.ModuleName,
		banktypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName, minttypes.ModuleName, genutiltypes.ModuleName, wasmtypes.ModuleName,
		evidencetypes.ModuleName, paramstypes.ModuleName, upgradetypes.ModuleName, vestingtypes.ModuleName, fantokentypes.ModuleName,
	}
}

func orderInitBlockers() []string {
	return []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		fantokentypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authtypes.ModuleName,
		genutiltypes.ModuleName,
		packetforwardtypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		wasmtypes.ModuleName,
	}
}

func simulationModules(
	app *BitsongApp,
	encodingConfig encparams.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler

	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.AppKeepers.BankKeeper, app.AppKeepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		// fantoken.NewAppModule(appCodec, app.FanTokenKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper),
		capability.NewAppModule(appCodec, *app.AppKeepers.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AppKeepers.AuthzKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.AppKeepers.GovKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.AppKeepers.MintKeeper, app.AppKeepers.AccountKeeper, nil, app.GetSubspace(stakingtypes.ModuleName)), // todo: replace nil w/ inflation reward calculation function
		staking.NewAppModule(appCodec, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.AppKeepers.DistrKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.AppKeepers.SlashingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.AppKeepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		params.NewAppModule(app.AppKeepers.ParamsKeeper),
		evidence.NewAppModule(app.AppKeepers.EvidenceKeeper),
		ibc.NewAppModule(app.AppKeepers.IBCKeeper),
		transfer.NewAppModule(app.AppKeepers.TransferKeeper),
		wasm.NewAppModule(appCodec, &app.AppKeepers.WasmKeeper, app.AppKeepers.StakingKeeper, app.AppKeepers.AccountKeeper, app.AppKeepers.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasm.ModuleName)),
	}
}
