package v013

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV13UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return newVM, err
		}

		logger.Info("Initialize wasm params...")
		params := keepers.WasmKeeper.GetParams(ctx)
		params.CodeUploadAccess = wasmtypes.AllowNobody
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
		keepers.WasmKeeper.SetParams(ctx, params)

		return newVM, err
	}
}
