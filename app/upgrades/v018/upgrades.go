package v018

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

func CreateV18UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, keepers *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		logger.Info(`
			; 
			;;
			;';. 
			;  ;;
			;   ;;
			;    ;;
			;    ;;
			;   ;'
			;  ' 
		,;;;,; 
		;;;;;;
		;;;;; 

		##     ##      #####          ##    #######  
		##     ##     ##   ##       ####   ##     ## 
		##     ##    ##     ##        ##   ##     ##
		##     ##    ##     ##        ##    ####### 
	 	 ##   ##     ##     ##        ##   ##     ##
		  ## ##       ##   ##  ###    ##   ##     ##
		   ###         #####   ###  ######  ####### 
	
			; 
			;;
			;';. 
			;  ;;
			;   ;;
			;    ;;
			;    ;;
			;   ;'
			;  ' 
		,;;;,; 
		;;;;;;
		;;;;; 
				
		`)

		// https://github.com/cosmos/cosmos-sdk/pull/12363/files
		// Set param key table for params module migration
		for _, subspace := range keepers.ParamsKeeper.GetSubspaces() {
			subspace := subspace

			var keyTable paramstypes.KeyTable
			switch subspace.Name() {
			case authtypes.ModuleName:
				keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
			case banktypes.ModuleName:
				keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
			case stakingtypes.ModuleName:
				keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck
			case distrtypes.ModuleName:
				keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
			case minttypes.ModuleName:
				keyTable = minttypes.ParamKeyTable()
			case slashingtypes.ModuleName:
				keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
			case govtypes.ModuleName:
				keyTable = govv1.ParamKeyTable() //nolint:staticcheck
			case crisistypes.ModuleName:
				keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
			// ibc types
			case ibctransfertypes.ModuleName:
				keyTable = ibctransfertypes.ParamKeyTable()
			// bitsong modules
			case fantokentypes.ModuleName:
				keyTable = fantokentypes.ParamKeyTable() //nolint:staticcheck
				if !subspace.HasKeyTable() {
					subspace.WithKeyTable(keyTable)
				}
			// // wasmd module
			// case wasmtypes.ModuleName:
			// 	keyTable = wasmtypes.ParamKeyTable() //nolint:staticcheck
			default:
				continue
			}

			if !subspace.HasKeyTable() {
				subspace.WithKeyTable(keyTable)
			}
		}

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		baseAppLegacySS := keepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(sdkCtx, baseAppLegacySS, &keepers.ConsensusParamsKeeper.ParamsStore)

		// manually cache current x/mint params for upgrade
		mintParams, _ := keepers.MintKeeper.Params.Get(sdkCtx)
		mintParams.MintDenom = appparams.MicroCoinUnit
		keepers.MintKeeper.Params.Set(ctx, mintParams)

		// manually cache current x/bank params for upgrade
		bankParams := keepers.BankKeeper.GetParams(ctx)
		bankParams.DefaultSendEnabled = true
		bankParams.SendEnabled = []*banktypes.SendEnabled{}
		logger.Info(fmt.Sprintf("bankParamaters: %v ", bankParams))
		keepers.BankKeeper.SetParams(ctx, bankParams)

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := keepers.IBCKeeper.ClientKeeper.GetParams(sdkCtx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		keepers.IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)

		// update gov params to use a 50% initial deposit ratio
		govParams, err := keepers.GovKeeper.Params.Get(sdkCtx)
		govParams.MinInitialDepositRatio = math.LegacyNewDec(50).Quo(math.LegacyNewDec(100)).String()
		if err := keepers.GovKeeper.Params.Set(sdkCtx, govParams); err != nil {
			return nil, err
		}

		return versionMap, err
	}
}
