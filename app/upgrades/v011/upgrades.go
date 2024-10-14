package v011

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/keepers"
	appparams "github.com/bitsongofficial/go-bitsong/v018/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV11UpgradeHandler(mm *module.Manager, configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		ctx.Logger().Info("Updating fantoken fees")
		ftParams := keepers.FanTokenKeeper.GetParamSet(ctx)
		ftParams.IssueFee.Denom = appparams.DefaultBondDenom
		ftParams.MintFee.Denom = appparams.DefaultBondDenom
		ftParams.BurnFee.Denom = appparams.DefaultBondDenom
		keepers.FanTokenKeeper.SetParamSet(ctx, ftParams)

		ctx.Logger().Info("Updating merkledrop fees")
		// mParams := keepers.MerkledropKeeper.GetParamSet(ctx)
		// mParams.CreationFee.Denom = appparams.DefaultBondDenom
		// keepers.MerkledropKeeper.SetParamSet(ctx, mParams)

		return newVM, err
	}
}
