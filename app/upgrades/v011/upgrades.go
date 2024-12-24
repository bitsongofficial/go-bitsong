package v011

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV11UpgradeHandler(mm *module.Manager, configurator module.Configurator,
	bpm upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		logger.Info("Updating fantoken fees")
		ftParams := keepers.FanTokenKeeper.GetParamSet(sdkCtx)
		ftParams.IssueFee.Denom = appparams.DefaultBondDenom
		ftParams.MintFee.Denom = appparams.DefaultBondDenom
		ftParams.BurnFee.Denom = appparams.DefaultBondDenom
		keepers.FanTokenKeeper.SetParamSet(sdkCtx, ftParams)

		logger.Info("Updating merkledrop fees")
		// mParams := keepers.MerkledropKeeper.GetParamSet(ctx)
		// mParams.CreationFee.Denom = appparams.DefaultBondDenom
		// keepers.MerkledropKeeper.SetParamSet(ctx, mParams)

		return newVM, err
	}
}
