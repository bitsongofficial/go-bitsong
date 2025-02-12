package v016

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV16UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, keepers *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
