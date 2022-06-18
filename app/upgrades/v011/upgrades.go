package v011

import (
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	fantokenkeeper "github.com/bitsongofficial/go-bitsong/x/fantoken/keeper"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator,
	ftk *fantokenkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		params := fantokentypes.DefaultParams()
		params.IssueFee.Denom = appparams.DefaultBondDenom
		params.MintFee.Denom = appparams.DefaultBondDenom
		params.BurnFee.Denom = appparams.DefaultBondDenom
		params.TransferFee.Denom = appparams.DefaultBondDenom

		ftk.SetParamSet(ctx, params)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
