package v013

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bitsongofficial/go-bitsong/v018/app/keepers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV13UpgradeHandler(mm *module.Manager, configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		ctx.Logger().Info("Initialize wasm params...")
		params := keepers.WasmKeeper.GetParams(ctx)
		params.CodeUploadAccess = wasmtypes.AllowNobody
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
		keepers.WasmKeeper.SetParams(ctx, params)

		return newVM, err
	}
}
