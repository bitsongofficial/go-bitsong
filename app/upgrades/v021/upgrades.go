package v021

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV021UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// remove unbonded validators data from x/dist store as staking hook did not due to v018 bug.

		// setup vote extension
		consensusParams, err := k.ConsensusParamsKeeper.ParamsStore.Get(ctx)
		if err != nil {
			return nil, err
		}
		// Update consensus params in order to safely enable comet pruning
		// consensusParams.Evidence.MaxAgeNumBlocks = 1_000_000
		// consensusParams.Evidence.MaxAgeDuration = time.Second * 1209600
		consensusParams.Abci.VoteExtensionsEnableHeight = 69
		err = k.ConsensusParamsKeeper.ParamsStore.Set(ctx, consensusParams)
		if err != nil {
			return nil, err
		}
		// setup icq params
		// setup ibchook params
		// setup clock params
		// setup smart account params
		// setup wasm client

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
