package v010

import (
	"github.com/bitsongofficial/go-bitsong/v018/app/keepers"
	appparams "github.com/bitsongofficial/go-bitsong/v018/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
)

func CreateV10UpgradeHandler(mm *module.Manager, configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		keepers.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

		fromVM := make(map[string]uint64)
		for moduleName := range mm.Modules {
			fromVM[moduleName] = 1
		}
		// delete new modules from the map, for _new_ modules as to not skip InitGenesis
		delete(fromVM, authz.ModuleName)
		delete(fromVM, feegrant.ModuleName)
		delete(fromVM, packetforwardtypes.ModuleName)

		// make fromVM[authtypes.ModuleName] = 2 to skip the first RunMigrations for auth (because from version 2 to migration version 2 will not migrate)
		fromVM[authtypes.ModuleName] = 2

		// the first RunMigrations, which will migrate all the old modules except auth module
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		// now update auth version back to 1, to make the second RunMigrations includes only auth
		newVM[authtypes.ModuleName] = 1

		// Proposal #5
		// Force an update of validator min commission
		validators := keepers.StakingKeeper.GetAllValidators(ctx)
		minCommissionRate := sdk.NewDecWithPrec(5, 2)
		for _, v := range validators {
			if v.Commission.Rate.LT(minCommissionRate) {
				if v.Commission.MaxRate.LT(minCommissionRate) {
					v.Commission.MaxRate = minCommissionRate
				}

				v.Commission.Rate = minCommissionRate
				v.Commission.UpdateTime = ctx.BlockHeader().Time

				// call the before-modification hook since we're about to update the commission
				// staking.BeforeValidatorModified(ctx, v.GetOperator())

				keepers.StakingKeeper.SetValidator(ctx, v)
			}
		}

		// Proposal #6
		// Mint BTSGs for Cassini-Bridge
		multisigWallet, err := sdk.AccAddressFromBech32(CassiniMultiSig)
		if err != nil {
			return nil, err
		}
		mintCoins := sdk.NewCoins(sdk.NewCoin(appparams.DefaultBondDenom, sdk.NewInt(CassiniMintAmount)))

		// mint coins
		if err := keepers.BankKeeper.MintCoins(ctx, minttypes.ModuleName, mintCoins); err != nil {
			return nil, err
		}

		if err := keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, multisigWallet, mintCoins); err != nil {
			return nil, err
		}

		// RunMigrations twice is just a way to make auth module's migrates after staking
		return mm.RunMigrations(ctx, configurator, newVM)
	}
}
