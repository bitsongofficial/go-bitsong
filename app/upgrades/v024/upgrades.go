package v024

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
)

func CreateV024UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(context)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		if _, found := vm[hyperlanetypes.ModuleName]; found {
			return nil, fmt.Errorf("%s already present in version map; its InitGenesis would be skipped", hyperlanetypes.ModuleName)
		}
		if _, found := vm[warptypes.ModuleName]; found {
			return nil, fmt.Errorf("%s already present in version map; its InitGenesis would be skipped", warptypes.ModuleName)
		}

		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))

		versionMap, err := mm.RunMigrations(sdkCtx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// materialize the hyperlane and warp module accounts at
		// the upgrade height. GetModuleAccount lazily creates and persists the account
		// the first time it is read; warp would otherwise only get created on its first
		// synthetic-token mint. Forcing creation here guarantees the accounts exist
		// exactly at the upgrade block.
		_ = k.AccountKeeper.GetModuleAccount(sdkCtx, hyperlanetypes.ModuleName)
		_ = k.AccountKeeper.GetModuleAccount(sdkCtx, warptypes.ModuleName)

		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		logger.Info("Hyperlane (core + warp) modules initialized")
		return versionMap, nil
	}
}
