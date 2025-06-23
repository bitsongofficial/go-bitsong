package v023

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func CreateV023UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		sdkCtx := sdk.UnwrapSDKContext(context)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// Run migrations first
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(sdkCtx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// apply logic patch
		err = CustomV023PatchLogic(sdkCtx, k)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}

// transfer all funds from protocolpool module account to distirbution module account
func CustomV023PatchLogic(ctx sdk.Context, k *keepers.AppKeepers) error {

	// perform patch on both of the accounts the module uses, just to be safe
	protocolPoolAddr := protocolpooltypes.ModuleName
	protocolPoolEscrowAddr := protocolpooltypes.ProtocolPoolEscrowAccount

	targetModule := distributiontypes.ModuleName

	if err := transferModuleBalance(ctx, k, protocolPoolAddr, targetModule); err != nil {
		return err
	}

	if err := transferModuleBalance(ctx, k, protocolPoolEscrowAddr, targetModule); err != nil {
		return err
	}

	return nil
}

func transferModuleBalance(ctx sdk.Context, k *keepers.AppKeepers, moduleName, targetModule string) error {
	moduleAddr := k.AccountKeeper.GetModuleAddress(moduleName)
	balances, err := k.BankKeeper.AllBalances(ctx, types.NewQueryAllBalancesRequest(moduleAddr, nil, false))
	if err != nil {
		return err
	}

	// no need to call bankkeeper if module has no balance
	if balances.Balances.Len() == 0 {
		return nil
	}
	// send tokens from module to module
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, moduleName, targetModule, balances.Balances)
	if err != nil {
		return err
	}
	// now we need to update the feepool value (mimics logic used in fundCommunityPool, but we must reimplement due to sending tokens from module to module)
	feePool, err := k.DistrKeeper.FeePool.Get(ctx)
	if err != nil {
		return err
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(balances.Balances...)...)
	return k.DistrKeeper.FeePool.Set(ctx, feePool)
}
