package v024

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pelletier/go-toml/v2"

	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV024Upgrade(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers, homepath string) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		sdkCtx := sdk.UnwrapSDKContext(context)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// Run migrations first
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(sdkCtx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// retrieve config.toml
		appConfigPath := filepath.Join(homepath, "config", "config.toml")
		configBytes, err := os.ReadFile(appConfigPath)
		if err != nil {
			return nil, err
		}
		// unmarshal file
		var config map[string]interface{}
		if err := toml.Unmarshal(configBytes, &config); err != nil {
			return nil, err
		}

		// update block speed to 2.4s
		if consensus, ok := config["consensus"].(map[string]interface{}); ok {
			consensus["timeout_commit"] = "2400ms"  // 2.4s
			consensus["timeout_propose"] = "2400ms" // 2.4s
		}
		// apply changes to config file
		updatedBytes, err := toml.Marshal(config)
		os.WriteFile(appConfigPath, updatedBytes, 0644)

		// update mint keepers blocks_per_year to reflect new block speed
		mp, err := k.MintKeeper.Params.Get(sdkCtx)
		if err != nil {
			return nil, err
		}
		mp.BlocksPerYear = 13148719 // @ 31556925 seconds per tropical year (365 days, 5 hours, 48 mins, 45 seconds)
		k.MintKeeper.Params.Set(sdkCtx, mp)

		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}
