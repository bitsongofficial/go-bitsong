package app

// to run: bitsongd custom-export --output-document v0.21.5-export.json
// defaults to for zero height
import (
	"encoding/json"
	"fmt"
	"log"

	storetypes "cosmossdk.io/store/types"
	v022 "github.com/bitsongofficial/go-bitsong/app/upgrades/v022"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

//   - validator slash event:
//   - historical reference:
//
// x/slashing:
//   - validator signing info:
//   - historical reference:
func (app *BitsongApp) CustomExportAppStateAndValidators(

	forZeroHeight bool, jailAllowedAddrs []string,

) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true)

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.customTestUpgradeHandlerLogicViaExport(ctx, jailAllowedAddrs)
	}

	genState, _ := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.AppKeepers.StakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//
//	in favour of export at a block height
func (app *BitsongApp) customTestUpgradeHandlerLogicViaExport(ctx sdk.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

	condJSON := v022.ConditionalJSON{
		PatchDelegationCount: 0,
		ZeroSharesDelegation: make([]v022.ZeroSharesDelegation, 0),
		PatchedDelegation:    make([]v022.PatchedDelegation, 0),
		DistSlashStore:       v022.DistrSlashObject{},
	}

	// check if there is a allowed address list
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)
	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		allowedAddrsMap[addr] = true

	}
	// // debug module
	err := v022.CustomV022PatchLogic(ctx, &app.AppKeepers, true)
	if err != nil {
		panic(err)
	}

	// simapp export for 0 height state logic:
	// withdraw all validator commission
	err = app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.AppKeepers.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		_, _ = app.AppKeepers.DistrKeeper.WithdrawValidatorCommission(ctx, valBz)
		return false
	})
	if err != nil {
		panic(err)
	}

	// count slashes  for all current  validators
	slashCount := uint64(0)
	app.AppKeepers.DistrKeeper.IterateValidatorSlashEvents(ctx,
		func(val sdk.ValAddress, _ uint64, vse distrtypes.ValidatorSlashEvent) (stop bool) {
			var existing *v022.DistrSlashEvent
			for i, obj := range condJSON.DistSlashStore.DistrSlashEvent {
				if obj.Val == val.String() {
					existing = &condJSON.DistSlashStore.DistrSlashEvent[i]
					break
				}
			} // Create a new DistrSlashEvent from the current event
			event := v022.Slash{
				Fraction: vse.Fraction.String(),
				Period:   vse.ValidatorPeriod,
			}

			// If existing object found, append the event and increment count
			if existing != nil {
				existing.Slashes = append(existing.Slashes, event)
				existing.SlashEventCount++
			} else {
				// Create a new DistrSlashObject
				newSlashObj := v022.DistrSlashEvent{
					Val:             val.String(),
					SlashEventCount: 1,
					Slashes:         []v022.Slash{event},
				}
				condJSON.DistSlashStore.DistrSlashEvent = append(condJSON.DistSlashStore.DistrSlashEvent, newSlashObj)
			}
			slashCount++
			return false
		})
	condJSON.DistSlashStore.SlashEventCount = slashCount

	// Marshal the ConditionalJSON object to JSON
	v022.PrintConditionalJsonLogs(condJSON, "conditional.json")
	// /* Just to be safe, assert the invariants on current state. */
	app.AppKeepers.CrisisKeeper.AssertInvariants(ctx)

	app.AppKeepers.StakingKeeper.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) (stop bool) {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		// delAddr := sdk.AccAddress(del.DelegatorAddress)
		val, err := app.AppKeepers.StakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			panic(err)
		}
		endingPeriod, err := app.AppKeepers.DistrKeeper.IncrementValidatorPeriod(ctx, val)
		if err != nil {
			panic(err)
		}
		rewardsRaw, patched := v022.CustomCalculateDelegationRewards(ctx, &app.AppKeepers, val, del, endingPeriod)

		if patched {
			condJSON.PatchedDelegation = append(condJSON.PatchedDelegation, v022.PatchedDelegation{
				OperatorAddress:   val.OperatorAddress,
				DelegatorAddress:  del.DelegatorAddress,
				PatchedDelegation: rewardsRaw.AmountOf("ubtsg").String(),
			})

			//export logic omits assertion of error, as shown here: https://github.com/permissionlessweb/bitsong-cosmos-sdk/blob/92e7f67f9601e5f2dc7b1daebd24e8a37efbcc72/simapp/export.go#L106
			_, err := v022.CustomWithdrawDelegationRewards(ctx, &app.AppKeepers, val, del, endingPeriod)

			if err != nil {
				fmt.Printf("err: %v\n", err)
				return false
			}

		} else {
			fmt.Printf("val.OperatorAddress: %v\n", val.OperatorAddress)
			fmt.Printf("del.DelegatorAddress: %v\n", del.DelegatorAddress)
			return false

		}
		if val.OperatorAddress == "bitsongvaloper1qxw4fjged2xve8ez7nu779tm8ejw92rv0vcuqr" ||
			val.OperatorAddress == "bitsongvaloper1xnc32z84cc9vwftvv4w0v02a2slug3tjt6qyct" {

			// // end current period and calculate rewards
			// endingPeriod, err := app.AppKeepers.DistrKeeper.IncrementValidatorPeriod(ctx, val)
			// if err != nil {
			// 	panic(err)
			// }

			// panic if we cannot calculate rewards for impacted validators
			_, err = app.AppKeepers.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
			if err != nil {
				panic(err)
			}

		}
		return false

	})

	// withdraw all delegator rewards
	dels, err := app.AppKeepers.StakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		panic(err)
	}

	// Marshal the ConditionalJSON object to JSON
	v022.PrintConditionalJsonLogs(condJSON, "conditional.json")

	// clear validator slash events
	app.AppKeepers.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.AppKeepers.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// reinitialize all validators
	err = app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.AppKeepers.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps, err := app.AppKeepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valBz)
		if err != nil {
			panic(err)
		}
		feePool, err := app.AppKeepers.DistrKeeper.FeePool.Get(ctx)
		if err != nil {
			panic(err)
		}
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		if err := app.AppKeepers.DistrKeeper.FeePool.Set(ctx, feePool); err != nil {
			panic(err)
		}

		if err := app.AppKeepers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, valBz); err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// reinitialize all delegations
	for _, del := range dels {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr := sdk.MustAccAddressFromBech32(del.DelegatorAddress)

		if err := app.AppKeepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr); err != nil {
			// never called as BeforeDelegationCreated always returns nil
			panic(fmt.Errorf("error while incrementing period: %w", err))
		}

		if err := app.AppKeepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr); err != nil {
			// never called as AfterDelegationModified always returns nil
			panic(fmt.Errorf("error while creating a new delegation period record: %w", err))
		}
	}
	// iterate through redelegations, reset creation height
	app.AppKeepers.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		err = app.AppKeepers.StakingKeeper.SetRedelegation(ctx, red)
		if err != nil {
			panic(err)
		}
		return false
	})
	// iterate through unbonding delegations, reset creation height
	app.AppKeepers.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		err = app.AppKeepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		if err != nil {
			panic(err)
		}
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.GetKey(stakingtypes.StoreKey))
	iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
		validator, err := app.AppKeepers.StakingKeeper.GetValidator(ctx, addr)
		if err != nil {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.AppKeepers.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	if err := iter.Close(); err != nil {
		app.Logger().Error("error while closing the key-value store reverse prefix iterator: ", err)
		return
	}

	_, err = app.AppKeepers.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/* Handle slashing state. */

	// reset start height on signing infos
	app.AppKeepers.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.AppKeepers.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)

}
