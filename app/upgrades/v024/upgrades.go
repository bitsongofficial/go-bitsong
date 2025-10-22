package v024

import (
	"context"
	"fmt"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

		// Update consensus params in order to safely enable comet pruning
		consensusParams, err := k.ConsensusParamsKeeper.ParamsStore.Get(sdkCtx)
		if err != nil {
			return nil, err
		}
		consensusParams.Evidence.MaxAgeNumBlocks = 756_000
		consensusParams.Evidence.MaxAgeDuration = time.Second * 1_814_400 // 21 days (in seconds)
		err = k.ConsensusParamsKeeper.ParamsStore.Set(sdkCtx, consensusParams)
		if err != nil {
			return nil, err
		}

		// update mint keepers blocks_per_year to reflect new block speed
		mp, err := k.MintKeeper.Params.Get(sdkCtx)
		if err != nil {
			return nil, err
		}
		mp.BlocksPerYear = 13148719 // @ 31556925 seconds per tropical year (365 days, 5 hours, 48 mins, 45 seconds)
		err = k.MintKeeper.Params.Set(sdkCtx, mp)
		if err != nil {
			return nil, err
		}
		// retain signed blocks duration given new block speeds
		p, err := k.SlashingKeeper.GetParams(sdkCtx)
		if err != nil {
			return nil, err
		}
		p.SignedBlocksWindow = 25_000 /// ~16.67 hours
		err = k.SlashingKeeper.SetParams(sdkCtx, p)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}
