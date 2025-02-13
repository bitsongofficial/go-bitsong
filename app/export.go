package app

import (
	"encoding/json"
	"fmt"
	"log"

	"cosmossdk.io/math"
	v020 "github.com/bitsongofficial/go-bitsong/app/upgrades/v020"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *BitsongApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
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
func (app *BitsongApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

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

	/* Just to be safe, assert the invariants on current state. */
	// app.AppKeepers.CrisisKeeper.AssertInvariants(ctx)

	/* Handle fee distribution state. */
	app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		dels := app.AppKeepers.StakingKeeper.GetValidatorDelegations(ctx, val.GetOperator())

		for _, del := range dels {
			ctx.Logger().Info(fmt.Sprintf("del_info: %q %v", val.GetOperator().String(), del.GetDelegatorAddr().String()))
			if del.Shares.LTE(math.LegacyZeroDec()) {
				ctx.Logger().Info(fmt.Sprintf("removing negative delegations: %q %v", val.GetOperator().String(), del.GetDelegatorAddr().String()))
				// remove reward information from distribution store
				if app.AppKeepers.DistrKeeper.HasDelegatorStartingInfo(ctx, val.GetOperator(), sdk.AccAddress(del.DelegatorAddress)) {
					app.AppKeepers.DistrKeeper.DeleteDelegatorStartingInfo(ctx, val.GetOperator(), sdk.AccAddress(del.DelegatorAddress))
				}
				// remove delegation from staking store
				if err := app.AppKeepers.StakingKeeper.RemoveDelegation(ctx, del); err != nil {
					panic(err)
				}
			} else {
				if !app.AppKeepers.DistrKeeper.HasDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr()) {
					panic(distrtypes.ErrEmptyDelegationDistInfo)
				}
				endingPeriod := app.AppKeepers.DistrKeeper.IncrementValidatorPeriod(ctx, val)
				rewardsRaw, patched := v020.CustomCalculateDelegationRewards(ctx, &app.AppKeepers, val, del, endingPeriod)
				outstanding := app.AppKeepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, del.GetValidatorAddr())

				if patched {
					ctx.Logger().Info("~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~")
					ctx.Logger().Info(fmt.Sprintf("PATCHED: %q %v", val.GetOperator().String(), del.GetDelegatorAddr().String()))
					err := v020.V018ManualDelegationRewardsPatch(ctx, rewardsRaw, outstanding, &app.AppKeepers, val, del, endingPeriod)
					if err != nil {
						panic(err)
					}
				}
			}

		}

		return false
	})

	/*  ensure no delegations exist without starting info*/
	for _, del := range app.AppKeepers.StakingKeeper.GetAllDelegations(ctx) {
		if !app.AppKeepers.DistrKeeper.HasDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr()) {
			panic(distrtypes.ErrEmptyDelegationDistInfo)
		}

		/* ensure all rewards are patched */
		val := app.AppKeepers.StakingKeeper.Validator(ctx, del.GetValidatorAddr())
		endingPeriod := app.AppKeepers.DistrKeeper.IncrementValidatorPeriod(ctx, val)
		/* will error if still broken */
		app.AppKeepers.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
		// ctx.Logger().Info(fmt.Sprintf("delegator reward: %q %v %q", val.GetOperator().String(), del.GetDelegatorAddr().String(), reward))
	}

	// withdraw all validator commission
	app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {

		_, err := app.AppKeepers.DistrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		if err != nil {
			ctx.Logger().Info(fmt.Sprintf("attempted to withdraw commission from validator with none, skipping: %q", val.GetOperator().String()))
			return false
		}
		return false
	})

	// withdraw all delegator rewards
	dels := app.AppKeepers.StakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		_, _ = app.AppKeepers.DistrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	}

	// clear validator slash events
	app.AppKeepers.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.AppKeepers.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.AppKeepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.AppKeepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
		feePool := app.AppKeepers.DistrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.AppKeepers.DistrKeeper.SetFeePool(ctx, feePool)

		app.AppKeepers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		app.AppKeepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr)
		app.AppKeepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr)
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.AppKeepers.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.AppKeepers.StakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.AppKeepers.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.AppKeepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[stakingtypes.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()[1:]
		addr := sdk.ValAddress(key)
		validator, found := app.AppKeepers.StakingKeeper.GetValidator(ctx, addr)
		if !found {
			ctx.Logger().Info(fmt.Sprintf("expected validator, not found: %q. removing key from store...", addr.String()))
			store.Delete(key)
			counter++
			continue
		}
		ctx.Logger().Info("-==-=-=-==---=-=-=-=-=--=-=-=-")
		ctx.Logger().Info(fmt.Sprintf("found: %q,  %q", addr.String(), validator.OperatorAddress))

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.AppKeepers.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()
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

// /* remove any remaining validator keys from store.This runs after we retrieve all current validators from staking keeper store,
//  preventing us from deleting active validators store. */
// store := sdkCtx.KVStore(k.GetKey(stakingtypes.StoreKey))
// iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
// counter := int16(0)

// for ; iter.Valid(); iter.Next() {
// 	key := iter.Key()[1:]
// 	addr := sdk.ValAddress(key)
// 	validator, err := k.StakingKeeper.GetValidator(sdkCtx, addr)
// 	if err != nil {
// 		sdkCtx.Logger().Info(fmt.Sprintf("expected validator, not found: %q", addr.String()))
// 		store.Delete(key)
// 		counter++
// 		continue
// 	} else {
// 		sdkCtx.Logger().Info("-==-=-=-==---=-=-=-=-=--=-=-=-")
// 		sdkCtx.Logger().Info(fmt.Sprintf("found: %q", validator.OperatorAddress))
// 	}
// 	counter++
// }
