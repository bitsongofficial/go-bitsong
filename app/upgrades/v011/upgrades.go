package v011

import (
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	fantokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	merkledropkeeper "github.com/bitsongofficial/go-bitsong/x/merkledrop/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator,
	ftk *fantokenkeeper.Keeper,
	mk *merkledropkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		ctx.Logger().Info("Updating fantoken fees")
		ftParams := ftk.GetParamSet(ctx)
		ftParams.IssueFee.Denom = appparams.DefaultBondDenom
		ftParams.MintFee.Denom = appparams.DefaultBondDenom
		ftParams.BurnFee.Denom = appparams.DefaultBondDenom
		ftk.SetParamSet(ctx, ftParams)

		ctx.Logger().Info("Updating merkledrop fees")
		mParams := mk.GetParamSet(ctx)
		mParams.CreationFee.Denom = appparams.DefaultBondDenom
		mk.SetParamSet(ctx, mParams)

		return newVM, err
	}
}
