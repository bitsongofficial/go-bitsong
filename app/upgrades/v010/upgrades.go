package v010

import (
	appparams "github.com/bitsongofficial/go-bitsong/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator,
	bank bankkeeper.Keeper,
	ibc *ibckeeper.Keeper,
	staking *stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ibc.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

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
		validators := staking.GetAllValidators(ctx)
		minCommissionRate := sdk.NewDecWithPrec(5, 2)
		for _, v := range validators {
			if v.Commission.Rate.LT(minCommissionRate) {
				if v.Commission.MaxRate.LT(minCommissionRate) {
					v.Commission.MaxRate = minCommissionRate
				}

				v.Commission.Rate = minCommissionRate
				v.Commission.UpdateTime = ctx.BlockHeader().Time

				// call the before-modification hook since we're about to update the commission
				staking.BeforeValidatorModified(ctx, v.GetOperator())

				staking.SetValidator(ctx, v)
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
		if err := bank.MintCoins(ctx, minttypes.ModuleName, mintCoins); err != nil {
			return nil, err
		}

		if err := bank.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, multisigWallet, mintCoins); err != nil {
			return nil, err
		}

		// RunMigrations twice is just a way to make auth module's migrates after staking
		return mm.RunMigrations(ctx, configurator, newVM)
	}
}
