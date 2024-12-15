package v020

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateV020UpgradeHandler(mm *module.Manager, configurator module.Configurator, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		ctx = sdk.UnwrapSDKContext(ctx)
		ctx.Logger().Info(`
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		V0182 UPGRADE manually claims delegation rewards for all users. 
		This will refresh the delegation information to the upgrade block.
		This prevents the error from occuring in the future.
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		`)

		// manually claim rewards by calling keeper functions
		for _, validator := range k.StakingKeeper.GetAllValidators(ctx) {
			for _, del := range k.StakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator()) {
				valAddr := del.GetValidatorAddr()
				val := k.StakingKeeper.Validator(ctx, valAddr)

				// check existence of delegator starting info
				if !k.DistrKeeper.HasDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr()) {
					return nil, distrtypes.ErrEmptyDelegationDistInfo
				}

				// end current period and calculate rewards
				endingPeriod := k.DistrKeeper.IncrementValidatorPeriod(ctx, val)
				rewardsRaw := customCalculateDelegationRewards(ctx, k, val, del, endingPeriod)
				outstanding := k.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, del.GetValidatorAddr())

				// defensive edge case may happen on the very final digits
				// of the decCoins due to operation order of the distribution mechanism.
				rewards := rewardsRaw.Intersect(outstanding)
				if !rewards.IsEqual(rewardsRaw) {
					logger := k.DistrKeeper.Logger(ctx)
					logger.Info(
						"rounding error withdrawing rewards from validator",
						"delegator", del.GetDelegatorAddr().String(),
						"validator", val.GetOperator().String(),
						"got", rewards.String(),
						"expected", rewardsRaw.String(),
					)
				}

				// truncate reward dec coins, return remainder to community pool
				finalRewards, remainder := rewards.TruncateDecimal()

				// add coins to user account
				if !finalRewards.IsZero() {
					ctx.Logger().Info("finalRewards", finalRewards)
					withdrawAddr := k.DistrKeeper.GetDelegatorWithdrawAddr(ctx, del.GetDelegatorAddr())
					err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, withdrawAddr, finalRewards)
					if err != nil {
						return nil, err
					}
				}

				// update the outstanding rewards and the community pool only if the
				// transaction was successful
				k.DistrKeeper.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), distrtypes.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
				feePool := k.DistrKeeper.GetFeePool(ctx)
				feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
				k.DistrKeeper.SetFeePool(ctx, feePool)

				// decrement reference count of starting period
				startingInfo := k.DistrKeeper.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
				startingPeriod := startingInfo.PreviousPeriod
				customDecrementReferenceCount(ctx, k, del.GetValidatorAddr(), startingPeriod)

				// remove delegator starting info
				k.DistrKeeper.DeleteDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())

				if finalRewards.IsZero() {
					baseDenom, _ := sdk.GetBaseDenom()
					if baseDenom == "" {
						baseDenom = sdk.DefaultBondDenom
					}

					// Note, we do not call the NewCoins constructor as we do not want the zero
					// coin removed.
					finalRewards = sdk.Coins{sdk.NewCoin(baseDenom, math.ZeroInt())}
					ctx.Logger().Info("finalRewards", finalRewards)
				}

				// reinitialize the delegation
				// period has already been incremented - we want to store the period ended by this delegation action
				previousPeriod := k.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddr).Period - 1

				// increment reference count for the period we're going to track
				incrementReferenceCount(ctx, k, valAddr, previousPeriod)

				validator := k.StakingKeeper.Validator(ctx, valAddr)
				delegation := k.StakingKeeper.Delegation(ctx, del.GetDelegatorAddr(), valAddr)

				// calculate delegation stake in tokens
				// we don't store directly, so multiply delegation shares * (tokens per share)
				// note: necessary to truncate so we don't allow withdrawing more rewards than owed
				stake := CustommTokensFromSharesTruncated(validator.GetTokens(), delegation.GetShares(), validator.GetDelegatorShares())

				// save new delegator starting info to kv store
				k.DistrKeeper.SetDelegatorStartingInfo(ctx, valAddr, del.GetDelegatorAddr(), distrtypes.NewDelegatorStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
			}
		}

		// // confirm patch has been applied by querying rewards again for each delegation
		// for _, del := range k.StakingKeeper.GetAllDelegations(ctx) {
		// 	valAddr := del.GetValidatorAddr()
		// 	val := k.StakingKeeper.Validator(ctx, valAddr)
		// 	// calculate rewards
		// 	k.DistrKeeper.CalculateDelegationRewards(ctx, val, del, uint64(ctx.BlockHeight()))
		// }

		ctx.Logger().Info(`
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-
		Upgrade V018 Patch complete. 
		All delegation rewards claimed and startingInfo set to this block height
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

func customCalculateDelegationRewards(ctx sdk.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins) {
	// fetch starting info for delegation
	startingInfo := k.DistrKeeper.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake
	startingHeight := startingInfo.Height
	endingHeight := uint64(ctx.BlockHeight())
	if endingHeight > startingHeight {
		k.DistrKeeper.IterateValidatorSlashEventsBetween(ctx, del.GetValidatorAddr(), startingHeight, endingHeight,
			func(height uint64, event distrtypes.ValidatorSlashEvent) (stop bool) {
				endingPeriod := event.ValidatorPeriod
				if endingPeriod > startingPeriod {
					rewards = rewards.Add(customCalculateDelegationRewardsBetween(ctx, k, val, startingPeriod, endingPeriod, stake)...)
					stake = stake.MulTruncate(math.LegacyOneDec().Sub(event.Fraction))
					startingPeriod = endingPeriod
				}
				return false
			},
		)
	}
	currentStake := val.TokensFromShares(del.GetShares())

	if stake.GT(currentStake) {
		marginOfErr := currentStake.Mul(sdk.NewDecWithPrec(12, 3)) // 1.2%
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
		} else {
			// ok := CalculateRewardsForSlashedDelegators(ctx, k, val, del, currentStake, SLASHED_DELEGATORS)
			// if ok {
			// 	stake = currentStake
			// } else {
			// }
			panic(fmt.Sprintln("current stake is not delgator from slashed validator, and is more than maximum margin of error"))
		}
	}
	// calculate rewards for final period
	rewards = rewards.Add(customCalculateDelegationRewardsBetween(ctx, k, val, startingPeriod, endingPeriod, stake)...)
	return rewards
}

func customCalculateDelegationRewardsBetween(ctx sdk.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI,
	startingPeriod, endingPeriod uint64, stake sdk.Dec,
) (rewards sdk.DecCoins) {
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, val.GetOperator(), startingPeriod)
	ending := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, val.GetOperator(), endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	return
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func customDecrementReferenceCount(ctx sdk.Context, k *keepers.AppKeepers, valAddr sdk.ValAddress, period uint64) {
	historical := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {

		k.DistrKeeper.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}

// increment the reference count for a historical rewards value
func incrementReferenceCount(ctx sdk.Context, k *keepers.AppKeepers, valAddr sdk.ValAddress, period uint64) {
	historical := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// calculate the token worth of provided shares, truncated
func CustommTokensFromSharesTruncated(t math.Int, ds math.LegacyDec, shares sdk.Dec) math.LegacyDec {
	return (shares.MulInt(t)).QuoTruncate(ds)
}
