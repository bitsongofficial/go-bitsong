package v020

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/app/keepers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV182UpgradeHandler(mm *module.Manager, configurator module.Configurator, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		ctx.Logger().Info(`
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		V0182 UPGRADE manually claims delegation rewards for all users. 
		This will refresh the delegation information to the upgrade block.
		This prevent the error from occuring in the future.
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		`)

		// manually claim rewards for all delegators slashed in validator
		for _, val := range k.StakingKeeper.GetAllValidators(ctx) {
			// get delegations
			val := sdk.ValAddress(val.GetOperator())
			delegation := k.StakingKeeper.GetValidatorDelegations(ctx, val)

			for _, del := range delegation {
				k.DistrKeeper.WithdrawDelegationRewards(ctx, del.GetDelegatorAddr(), val)
			}
		}
		// confirm calculations work as expected by checking rewards for every delgation.
		// This upgrade fails if any delegators still are impacted by v018 upgrade error.
		for _, del := range k.StakingKeeper.GetAllDelegations(ctx) {
			valAddr := del.GetValidatorAddr()
			val := k.StakingKeeper.Validator(ctx, valAddr)
			// manually claim reward
			k.DistrKeeper.WithdrawDelegationRewards(ctx, del.GetDelegatorAddr(), valAddr)
			// calculate rewards
			k.DistrKeeper.CalculateDelegationRewards(ctx, val, del, uint64(ctx.BlockHeight()))
		}

		ctx.Logger().Info(`
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		Upgrade V018 Patch complete. 
		All delegation rewards have been claimed. 
		Nodes are now able to regress upstream to the main cosmos-sdk module :)
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		`)

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))

		return versionMap, err
	}
}
