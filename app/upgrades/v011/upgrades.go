package v011

import (
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	fantokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator,
	ftk *fantokenkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		ctx.Logger().Info("Updating fantoken fees")
		params := ftk.GetParamSet(ctx)
		params.IssueFee.Denom = appparams.DefaultBondDenom
		params.MintFee.Denom = appparams.DefaultBondDenom
		params.BurnFee.Denom = appparams.DefaultBondDenom
		ftk.SetParamSet(ctx, params)

		return newVM, err
	}
}
