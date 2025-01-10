package v021

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

func CreateV021UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// remove unbonded validators data from x/dist store as staking hook did not due to v018 bug.

		// setup vote extension
		consensusParams, err := k.ConsensusParamsKeeper.ParamsStore.Get(sdkCtx)
		if err != nil {
			return nil, err
		}
		// Update consensus params in order to safely enable comet pruning
		// consensusParams.Evidence.MaxAgeNumBlocks = 1_000_000
		// consensusParams.Evidence.MaxAgeDuration = time.Second * 1209600
		consensusParams.Abci.VoteExtensionsEnableHeight = 69
		err = k.ConsensusParamsKeeper.ParamsStore.Set(sdkCtx, consensusParams)
		if err != nil {
			return nil, err
		}
		// setup icq params
		setICQParams(sdkCtx, k.ICQKeeper)
		// setup ibchook params
		// setup clock params
		// setup smart account params
		// setup wasm client
		// https://github.com/cosmos/ibc-go/blob/main/docs/docs/03-light-clients/04-wasm/03-integration.md
		params := k.IBCKeeper.ClientKeeper.GetParams(sdkCtx)
		params.AllowedClients = append(params.AllowedClients, wasmlctypes.Wasm)
		k.IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}

func setICQParams(ctx sdk.Context, icqKeeper *icqkeeper.Keeper) {
	icqparams := icqtypes.DefaultParams()
	// icqparams.AllowQueries = wasmbinding.GetStargateWhitelistedPaths()
	// Adding SmartContractState query to allowlist
	icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
	//nolint:errcheck
	icqKeeper.SetParams(ctx, icqparams)
}
