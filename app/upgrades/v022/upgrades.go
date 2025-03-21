package v022

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func CreateV022UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
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
		err = CustomV022PatchLogic(sdkCtx, k, false)

		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}

// / Fetches all validators from x/staking, and query rewards for each of their delegations.
func CustomV022PatchLogic(ctx sdk.Context, k *keepers.AppKeepers, simulated bool) error {

	err := CustomValPatch(ctx, k, simulated)
	if err != nil {
		return err
	}

	return nil
}

// withdraws all rewards to delegations for this val.
func CustomValPatch(sdkCtx sdk.Context, k *keepers.AppKeepers, simulated bool) error {
	condJSON := ConditionalJSON{
		PatchDelegationCount:     0,
		PatchedHistRewards:       make([]distrtypes.ValidatorHistoricalRewardsRecord, 0),
		PatchedDelegation:        make([]PatchedDelegation, 0),
		ZeroSharesDelegation:     make([]ZeroSharesDelegation, 0),
		NilDelegationCalculation: make([]NilDelegationCalculation, 0),
		DistSlashStore:           DistrSlashObject{},
	}

	allVals, err := k.StakingKeeper.GetAllValidators(sdkCtx)
	if err != nil {
		return err
	}

	for _, val := range allVals {
		valAddr, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.OperatorAddress)
		if err != nil {
			return err
		}

		dels, err := k.StakingKeeper.GetValidatorDelegations(sdkCtx, sdk.ValAddress(valAddr))
		if err != nil {
			return err
		}
		/* patch all delegators rewards */
		for _, del := range dels {
			//  if delegation has 0 shares, remove from store: https://github.com/permissionlessweb/bitsong-cosmos-sdk/blob/92e7f67f9601e5f2dc7b1daebd24e8a37efbcc72/x/staking/keeper/delegation.go#L1026
			if del.GetShares().LTE(math.LegacyZeroDec()) {
				// remove from store
				err = k.StakingKeeper.RemoveDelegation(sdkCtx, stakingtypes.NewDelegation(del.GetDelegatorAddr(), del.GetValidatorAddr(), del.GetShares()))
				condJSON.ZeroSharesDelegation = append(condJSON.ZeroSharesDelegation, ZeroSharesDelegation{
					OperatorAddress:  del.ValidatorAddress,
					DelegatorAddress: del.DelegatorAddress,
				})
				if err != nil {
					return err
				}
				continue
			}
			val, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
			if err != nil {
				return err
			}
			delAddr, err := k.AccountKeeper.AddressCodec().StringToBytes(del.GetDelegatorAddr())
			if err != nil {
				return err
			}
			validator, err := k.StakingKeeper.Validator(sdkCtx, val)
			if err != nil {
				return err
			}

			// end current period and calculate rewards
			endingPeriod, err := k.DistrKeeper.IncrementValidatorPeriod(sdkCtx, validator)
			if err != nil {
				return err
			}

			delegationRewards, patched := CustomCalculateDelegationRewards(sdkCtx, k, validator, del, endingPeriod)
			if patched {
				condJSON.PatchDelegationCount++
				condJSON.PatchedDelegation = append(condJSON.PatchedDelegation, PatchedDelegation{
					OperatorAddress:   del.ValidatorAddress,
					DelegatorAddress:  del.DelegatorAddress,
					PatchedDelegation: delegationRewards.AmountOf("ubtsg").String(),
				})
				if delegationRewards == nil {
					condJSON.NilDelegationCalculation = append(condJSON.NilDelegationCalculation, NilDelegationCalculation{
						OperatorAddress:  del.ValidatorAddress,
						DelegatorAddress: del.DelegatorAddress,
					})
				}
				_, err := CustomWithdrawDelegationRewards(sdkCtx, k, validator, del, endingPeriod)

				if err != nil {
					return err

				}

				// reinitialize the delegation
				err = customInitializeDelegation(sdkCtx, *k, val, delAddr)
				if err != nil {
					return err
				}
			}

		}

		if val.OperatorAddress == PatchVal1 || val.OperatorAddress == PatchVal2 {
			rewards, err := k.DistrKeeper.GetValidatorCurrentRewards(sdkCtx, valAddr)
			if err != nil {
				return err
			}
			fmt.Printf("rewards: %v\n", rewards)
			// set ghost slash event
			err = k.DistrKeeper.SetValidatorSlashEvent(sdkCtx, valAddr, 1, rewards.Period, distrtypes.NewValidatorSlashEvent(rewards.Period, math.LegacySmallestDec()))
			if err != nil {
				return err
			}
		}
	}

	k.DistrKeeper.IterateValidatorHistoricalRewards(sdkCtx,
		func(val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
			if val.String() == PatchVal1 || val.String() == PatchVal2 {
				// print to logs and update validator reward reference to 1
				condJSON.PatchedHistRewards = append(condJSON.PatchedHistRewards, distrtypes.ValidatorHistoricalRewardsRecord{
					ValidatorAddress: val.String(),
					Period:           period,
					Rewards:          rewards,
				})
				// update reference count to 1
				// rewards.ReferenceCount = 1
				// err := k.DistrKeeper.SetValidatorHistoricalRewards(sdkCtx, val, period, rewards)
				// if err != nil {
				// 	panic(err)
				// }
			}
			return false
		},
	)

	// count slashes for all current validators
	// Create a map to track validators and their slash events
	validatorSlashMap := make(map[string][]Slash)
	slashCount := uint64(0)
	k.DistrKeeper.IterateValidatorSlashEvents(sdkCtx,
		func(val sdk.ValAddress, height uint64, vse distrtypes.ValidatorSlashEvent) (stop bool) {
			valAddr := val.String()
			event := Slash{
				Height:   height,
				Fraction: vse.Fraction.String(),
				Period:   vse.ValidatorPeriod,
			}
			validatorSlashMap[valAddr] = append(validatorSlashMap[valAddr], event)
			slashCount++

			return false
		})

	condJSON.DistSlashStore.DistrSlashEvent = make([]map[string][]Slash, 0, len(validatorSlashMap))
	for valAddr, slashes := range validatorSlashMap {
		condJSON.DistSlashStore.DistrSlashEvent = append(condJSON.DistSlashStore.DistrSlashEvent,
			map[string][]Slash{valAddr: slashes})
	}
	condJSON.DistSlashStore.SlashEventCount = slashCount

	PrintConditionalJsonLogs(condJSON, "upgradeHandlerDebug.json")

	return nil
}

func PrintConditionalJsonLogs(condJSON ConditionalJSON, fileName string) error {
	// Marshal the ConditionalJSON object to JSON
	jsonBytes, err := json.MarshalIndent(condJSON, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal and indent conditional logs, continuing with upgrade... %s\n", fileName)
	}
	// Write the JSON to a file
	err = os.WriteFile(fileName, jsonBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write debugging log, continuing with upgrade... %s\n", fileName)
	}
	fmt.Printf("Wrote conditionals to %s\n", fileName)
	return nil
}

func CustomCalculateDelegationRewards(ctx context.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, patched bool) {
	patched = false
	// fetch starting info for delegation
	addrCodec := k.AccountKeeper.AddressCodec()
	delAddr, err := addrCodec.StringToBytes(del.GetDelegatorAddr())
	if err != nil {
		panic(err)
	}

	valAddr, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(del.GetValidatorAddr())
	if err != nil {
		panic(err)
	}
	startingInfo, err := k.DistrKeeper.GetDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
	if err != nil {
		panic(err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if startingInfo.Height == uint64(sdkCtx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake

	startingHeight := startingInfo.Height
	endingHeight := uint64(sdkCtx.BlockHeight())
	if endingHeight > startingHeight {
		k.DistrKeeper.IterateValidatorSlashEventsBetween(ctx, valAddr, startingHeight, endingHeight,
			func(height uint64, event distrtypes.ValidatorSlashEvent) (stop bool) {
				endingPeriod := event.ValidatorPeriod
				if endingPeriod > startingPeriod {
					rewards = rewards.Add(customCalculateDelegationRewardsBetween(sdkCtx, k, val, startingPeriod, endingPeriod, stake)...)
					// Note: It is necessary to truncate so we don't allow withdrawing
					// more rewards than owed.
					stake = stake.MulTruncate(math.LegacyOneDec().Sub(event.Fraction))
					startingPeriod = endingPeriod
				}
				return false
			},
		)
	}

	currentStake := val.TokensFromShares(del.GetShares())

	if stake.GT(currentStake) {
		marginOfErr := currentStake.Mul(math.LegacyNewDecWithPrec(95, 3)) // 0.095
		if stake.LTE(currentStake.Add(marginOfErr)) {
			fmt.Printf("PATCH APPLIED TO DELEGATOR: %v\n FOR VALIDATOR:  %v\n", del.GetDelegatorAddr(), del.GetValidatorAddr())
			fmt.Printf("distribution store: %v\n", stake)
			fmt.Printf("staking store: %v\n", currentStake)
			fmt.Printf("startingInfo.PreviousPeriod: %v\n", startingInfo.PreviousPeriod)
			stake = currentStake
			patched = true
		} else {
			fmt.Printf("marginOfErr: %v\n", marginOfErr)
			panic(fmt.Sprintf("calculated final stake for delegator %s,to validator %s, greater than current stake"+
				"\n\tfinal stake:\t%s"+
				"\n\tcurrent stake:\t%s",
				del.GetDelegatorAddr(), del.GetValidatorAddr(), stake, currentStake))
		}
	}
	// calculate rewards for final period
	rewards = rewards.Add(customCalculateDelegationRewardsBetween(ctx, k, val, startingPeriod, endingPeriod, stake)...)
	return rewards, patched
}

func CustomWithdrawDelegationRewards(ctx context.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (sdk.Coins, error) {
	addrCodec := k.AccountKeeper.AddressCodec()
	delAddr, err := addrCodec.StringToBytes(del.GetDelegatorAddr())
	if err != nil {
		return nil, err
	}
	valAddr, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(del.GetValidatorAddr())
	if err != nil {
		return nil, err
	}
	// check existence of delegator starting info
	hasInfo, err := k.DistrKeeper.HasDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
	if err != nil {
		return nil, err
	}
	if !hasInfo {
		return nil, distrtypes.ErrEmptyDelegationDistInfo
	}

	// custom calc reward for use elsewhere
	customCalculatedrewards, patched := CustomCalculateDelegationRewards(ctx, k, val, del, endingPeriod)
	if !patched {
		panic("applying patch to delegation that should not have been patched, panic!")
	}
	outstanding, err := k.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, sdk.ValAddress(valAddr))
	if err != nil {
		return nil, err
	}
	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := customCalculatedrewards.Intersect(outstanding)
	if !rewards.Equal(customCalculatedrewards) {
		logger := k.DistrKeeper.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from validator",
			"delegator", del.GetDelegatorAddr(),
			"validator", val.GetOperator(),
			"got", rewards.String(),
			"expected", customCalculatedrewards.String(),
		)
	}

	// truncate reward dec coins, return remainder to community pool
	finalRewards, remainder := rewards.TruncateDecimal()
	// add coins to user account
	if !finalRewards.IsZero() {
		withdrawAddr, err := k.DistrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return nil, err
		}

		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, withdrawAddr, finalRewards)

		if err != nil {
			senderAddr := k.AccountKeeper.GetModuleAddress(distrtypes.ModuleName)
			distrBal := k.BankKeeper.GetAllBalances(ctx, senderAddr)
			fmt.Printf("distrBal: %v\n", distrBal)
			fmt.Printf("delAddr: %v\n", del.GetDelegatorAddr())
			fmt.Printf("senderAddr.String(): %v\n", senderAddr.String())
			// panic("distribution module has less than what it should have")
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	err = k.DistrKeeper.SetValidatorOutstandingRewards(ctx, sdk.ValAddress(valAddr), distrtypes.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	if err != nil {
		return nil, err
	}

	feePool, err := k.DistrKeeper.FeePool.Get(ctx)
	if err != nil {
		return nil, err
	}

	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	err = k.DistrKeeper.FeePool.Set(ctx, feePool)
	if err != nil {
		return nil, err
	}
	// decrement reference count of starting period
	startingInfo, err := k.DistrKeeper.GetDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
	if err != nil {
		return nil, err
	}
	startingPeriod := startingInfo.PreviousPeriod
	err = customDecrementReferenceCount(ctx, k, sdk.ValAddress(valAddr), startingPeriod)
	if err != nil {
		return nil, err
	}
	// remove delegator starting info
	err = k.DistrKeeper.DeleteDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
	if err != nil {
		return nil, err
	}
	if finalRewards.IsZero() {
		baseDenom, _ := sdk.GetBaseDenom()
		if baseDenom == "" {
			baseDenom = sdk.DefaultBondDenom
		}

		// Note, we do not call the NewCoins constructor as we do not want the zero
		// coin removed.
		finalRewards = sdk.Coins{sdk.NewCoin(baseDenom, math.ZeroInt())}
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			distrtypes.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, finalRewards.String()),
			sdk.NewAttribute(distrtypes.AttributeKeyValidator, val.GetOperator()),
			sdk.NewAttribute(distrtypes.AttributeKeyDelegator, del.GetDelegatorAddr()),
		),
	)

	return finalRewards, nil
}

func customCalculateDelegationRewardsBetween(ctx context.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI,
	startingPeriod, endingPeriod uint64, stake math.LegacyDec,
) (rewards sdk.DecCoins) {
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}
	// fmt.Printf("startingPeriod: %v\n", startingPeriod)
	// fmt.Printf("endingPeriod: %v\n", endingPeriod)
	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}
	valBz, err := k.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	if err != nil {
		panic(err)
	}

	// return staking * (ending - starting)
	starting, err := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valBz, startingPeriod)
	if err != nil {
		panic(err)
	}
	ending, err := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valBz, endingPeriod)
	if err != nil {
		panic(err)
	}

	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// fmt.Printf("starting: %v\n", ending)
	// fmt.Printf("ending: %v\n", starting)
	// fmt.Printf("stake: %v\n", stake)
	// fmt.Printf("difference: %v\n", difference)
	// fmt.Printf("rewards: %v\n", rewards)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	return rewards
}

// increment the reference count for a historical rewards value
func customIncrementReferenceCount(ctx context.Context, k keepers.AppKeepers, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		fmt.Printf("valAddr.String(): %v\n", valAddr.String())
		fmt.Printf("historical.ReferenceCount: %v\n", historical.ReferenceCount)
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func customDecrementReferenceCount(ctx context.Context, k *keepers.AppKeepers, valAddr sdk.ValAddress, period uint64) error {
	historical, _ := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.DistrKeeper.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	}
	return k.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)

}

// initialize starting info for a new delegation
func customInitializeDelegation(ctx context.Context, k keepers.AppKeepers, val sdk.ValAddress, del sdk.AccAddress) error {
	// period has already been incremented - we want to store the period ended by this delegation action
	valCurrentRewards, err := k.DistrKeeper.GetValidatorCurrentRewards(ctx, val)
	if err != nil {
		return err
	}
	previousPeriod := valCurrentRewards.Period - 1

	// increment reference count for the period we're going to track
	err = customIncrementReferenceCount(ctx, k, val, previousPeriod)
	if err != nil {
		return err
	}

	validator, err := k.StakingKeeper.Validator(ctx, val)
	if err != nil {
		return err
	}

	delegation, err := k.StakingKeeper.Delegation(ctx, del, val)
	if err != nil {
		return err
	}

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.DistrKeeper.SetDelegatorStartingInfo(ctx, val, del, distrtypes.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))
}
