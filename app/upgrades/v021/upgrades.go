package v021

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	cadancetypes "github.com/bitsongofficial/go-bitsong/x/cadance/types"
	sca "github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"

	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
)

func CreateV021UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(context)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// Run migrations first
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(sdkCtx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// reapplies v018 patch after removing delegations with 0 power, letting us revert back upstream to cosmos-sdk library
		vals, _ := k.StakingKeeper.GetAllValidators(sdkCtx)
		for _, val := range vals {
			valAddr := sdk.ValAddress(val.OperatorAddress)
			dels, _ := k.StakingKeeper.GetValidatorDelegations(sdkCtx, valAddr)
			for _, del := range dels {
				if del.Shares.LTE(math.LegacyZeroDec()) {
					sdkCtx.Logger().Info(fmt.Sprintf("removing negative delegation from store: %q %v", val.GetOperator(), del.GetDelegatorAddr())) // remove reward information from distribution store
					if exists, _ := k.DistrKeeper.HasDelegatorStartingInfo(sdkCtx, valAddr, sdk.AccAddress(del.DelegatorAddress)); exists {
						sdkCtx.Logger().Info("delegation info found, deleting...")
						if err := k.DistrKeeper.DeleteDelegatorStartingInfo(sdkCtx, valAddr, sdk.AccAddress(del.DelegatorAddress)); err != nil {
							return nil, err
						}
						sdkCtx.Logger().Info("removed negative delegation from store successfully!")
					}
					// remove delegation from staking store
					sdkCtx.Logger().Info("~-~-~~-~-~~-~-~~-~-~~-~-~~-~-~~-~-~")
					sdkCtx.Logger().Info("removing negative delegation from staking keeper store")
					if err := k.StakingKeeper.RemoveDelegation(sdkCtx, del); err != nil {
						return nil, err
					} else {
						sdkCtx.Logger().Info("removed negative delegation store value successfully!")
					}
				} else {
					// check if we need to patch distribution by manually claiming rewards for any impaced delegations once again...
					hasInfo, err := k.DistrKeeper.HasDelegatorStartingInfo(sdkCtx, sdk.ValAddress(valAddr), sdk.AccAddress(del.GetDelegatorAddr()))
					if err != nil {
						return nil, err
					}
					if !hasInfo {
						sdkCtx.Logger().Info(fmt.Sprintf("delegation does not have starting info: val: %q, del: %v", val.GetOperator(), del.GetDelegatorAddr()))
						continue
					}
					// calculate rewards
					endingPeriod, err := k.DistrKeeper.IncrementValidatorPeriod(sdkCtx, val)
					if err != nil {
						return nil, err
					}
					rewardsRaw, patched := CustomCalculateDelegationRewards(sdkCtx, k, val, del, endingPeriod)

					outstanding, err := k.DistrKeeper.GetValidatorOutstandingRewardsCoins(sdkCtx, sdk.ValAddress(del.GetValidatorAddr()))
					if err != nil {
						return nil, err
					}

					if patched {
						err = V018ManualDelegationRewardsPatch(sdkCtx, rewardsRaw, outstanding, k, val, del, endingPeriod)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}

		/*  ensure no delegations exist without starting info*/
		allVals, err := k.StakingKeeper.GetAllValidators(sdkCtx)
		if err != nil {
			panic(err)
		}
		for _, val := range allVals {
			/* ensure all rewards are patched */
			dels, err := k.StakingKeeper.GetValidatorDelegations(sdkCtx, sdk.ValAddress(val.OperatorAddress))
			if err != nil {
				panic(err)
			}
			for _, del := range dels {
				endingPeriod, err := k.DistrKeeper.IncrementValidatorPeriod(sdkCtx, val)
				if err != nil {
					panic(err)
				}
				/* will error if still broken */
				_, err = k.DistrKeeper.CalculateDelegationRewards(sdkCtx, val, del, endingPeriod)
				if err != nil {
					panic(err)
				}
			}
		}

		// setup vote extension
		consensusParams, err := k.ConsensusParamsKeeper.ParamsStore.Get(sdkCtx)
		if err != nil {
			return nil, err
		}

		// TODO: these are values specific to 1.5s block time. We may need to tune this to bitsongs desired blocktime:
		// Update consensus params in order to safely enable comet pruning
		// consensusParams.Evidence.MaxAgeNumBlocks = 1_000_000
		// consensusParams.Evidence.MaxAgeDuration = time.Second * 1209600
		err = k.ConsensusParamsKeeper.ParamsStore.Set(sdkCtx, cmtproto.ConsensusParams{
			Block:     consensusParams.Block,
			Evidence:  consensusParams.Evidence,
			Validator: consensusParams.Validator,
			Version:   consensusParams.Version,
			Abci: &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: sdkCtx.BlockHeight() + 1,
			},
		})
		if err != nil {
			return nil, err
		}
		// interchain query params (ICQ)
		setICQParams(sdkCtx, k.ICQKeeper)

		// set x/cadance params
		cadanceParams := cadancetypes.DefaultParams()
		cadanceParams.ContractGasLimit = 1000000 // 1mb
		if err := k.CadanceKeeper.SetParams(sdkCtx, cadanceParams); err != nil {
			return nil, err
		}

		// Set the x/smart-account authenticator params in the store
		authenticatorParams := sca.DefaultParams()
		authenticatorParams.CircuitBreakerControllers = append(authenticatorParams.CircuitBreakerControllers, CircuitBreakerController)
		k.SmartAccountKeeper.SetParams(sdkCtx, authenticatorParams)

		// set wasm client as an allowed client.
		// https://github.com/cosmos/ibc-go/blob/main/docs/docs/03-light-clients/04-wasm/03-integration.md
		// ibcCLientParams := ibcclient.DefaultParams()
		params := k.IBCKeeper.ClientKeeper.GetParams(sdkCtx)
		params.AllowedClients = append(params.AllowedClients, wasmlctypes.Wasm)
		k.IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)

		// configure expidited proposals
		govparams, _ := k.GovKeeper.Params.Get(sdkCtx)
		govparams.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin("ubtsg", math.NewInt(10000000000))) // 10K
		newExpeditedVotingPeriod := time.Minute * 60 * 24                                            // 1 DAY
		govparams.ExpeditedVotingPeriod = &newExpeditedVotingPeriod
		govparams.ExpeditedThreshold = "0.75" // 75% voting threshold
		k.GovKeeper.Params.Set(sdkCtx, govparams)

		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}

func setICQParams(ctx sdk.Context, icqKeeper *icqkeeper.Keeper) {
	icqparams := icqtypes.DefaultParams()
	// icqparams.AllowQueries = wasmbinding.GetStargateWhitelistedPaths()
	// Adding SmartContractState query to allowlist
	icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
	//nolint:errcheck
	icqKeeper.SetParams(ctx, icqparams)
}

func V018ManualDelegationRewardsPatch(sdkCtx sdk.Context, rewardsRaw, outstanding sdk.DecCoins, k *keepers.AppKeepers, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) error {

	valAddr := del.GetValidatorAddr()
	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.Equal(rewardsRaw) {
		logger := k.DistrKeeper.Logger(sdkCtx)
		logger.Info(
			"rounding error withdrawing rewards from validator",
			"delegator", del.GetDelegatorAddr(),
			"validator", val.GetOperator(),
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	// truncate reward dec coins, return remainder to community pool
	finalRewards, remainder := rewards.TruncateDecimal()

	// add coins to user account
	if !finalRewards.IsZero() {
		withdrawAddr, err := k.DistrKeeper.GetDelegatorWithdrawAddr(sdkCtx, sdk.AccAddress(del.GetDelegatorAddr()))
		if err != nil {
			return err
		}
		err = k.BankKeeper.SendCoinsFromModuleToAccount(sdkCtx, distrtypes.ModuleName, withdrawAddr, finalRewards)
		if err != nil {
			return err
		}
		sdkCtx.Logger().Info(fmt.Sprintf("Rewards %v manually claimed for: %q", finalRewards, del.GetDelegatorAddr()))
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	k.DistrKeeper.SetValidatorOutstandingRewards(sdkCtx, sdk.ValAddress(valAddr), distrtypes.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	feePool, err := k.DistrKeeper.FeePool.Get(sdkCtx)
	if err != nil {
		return err
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	err = k.DistrKeeper.FeePool.Set(sdkCtx, feePool)
	if err != nil {
		return err
	}

	// decrement reference count of starting period
	startingInfo, err := k.DistrKeeper.GetDelegatorStartingInfo(sdkCtx, sdk.ValAddress(del.GetValidatorAddr()), sdk.AccAddress(del.GetDelegatorAddr()))
	if err != nil {
		return err
	}
	startingPeriod := startingInfo.PreviousPeriod
	customDecrementReferenceCount(sdkCtx, k, sdk.ValAddress(del.GetValidatorAddr()), startingPeriod)

	// remove delegator starting info
	k.DistrKeeper.DeleteDelegatorStartingInfo(sdkCtx, sdk.ValAddress(del.GetValidatorAddr()), sdk.AccAddress(del.GetDelegatorAddr()))

	if finalRewards.IsZero() {
		// Note, we do not call the NewCoins constructor as we do not want the zero
		// coin removed.
		sdkCtx.Logger().Info("~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~~=~=~=~=~")
		sdkCtx.Logger().Info(fmt.Sprintf("No final rewards: %q %v", val.GetOperator(), del.GetDelegatorAddr()))
	}

	// reinitialize the delegation
	// period has already been incremented - we want to store the period ended by this delegation action
	vcr, _ := k.DistrKeeper.GetValidatorCurrentRewards(sdkCtx, sdk.ValAddress(valAddr))
	previousPeriod := vcr.Period - 1
	// increment reference count for the period we're going to track
	incrementReferenceCount(sdkCtx, k, sdk.ValAddress(valAddr), previousPeriod)

	validator, _ := k.StakingKeeper.Validator(sdkCtx, sdk.ValAddress(valAddr))
	delegation, _ := k.StakingKeeper.Delegation(sdkCtx, sdk.AccAddress(del.GetDelegatorAddr()), sdk.ValAddress(valAddr))

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())

	// save new delegator starting info to kv store
	k.DistrKeeper.SetDelegatorStartingInfo(sdkCtx, sdk.ValAddress(valAddr), sdk.AccAddress(del.GetDelegatorAddr()), distrtypes.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))

	return nil
}

func CustomCalculateDelegationRewards(ctx sdk.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, patched bool) {
	patched = false
	// fetch starting info for delegation
	startingInfo, err := k.DistrKeeper.GetDelegatorStartingInfo(ctx, sdk.ValAddress(del.GetValidatorAddr()), sdk.AccAddress(del.GetDelegatorAddr()))
	if err != nil {
		return
	}
	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake
	startingHeight := startingInfo.Height
	endingHeight := uint64(ctx.BlockHeight())
	if endingHeight > startingHeight {
		k.DistrKeeper.IterateValidatorSlashEventsBetween(ctx, sdk.ValAddress(del.GetValidatorAddr()), startingHeight, endingHeight,
			func(height uint64, event distrtypes.ValidatorSlashEvent) (stop bool) {
				endingPeriod := event.ValidatorPeriod
				if endingPeriod > startingPeriod {
					rewards = rewards.Add(customCalculateDelegationRewardsBetween(ctx, k, val, startingPeriod, endingPeriod, stake)...)
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
	fmt.Printf("del: %q", del.GetDelegatorAddr())
	fmt.Printf("val: %q", del.GetValidatorAddr())
	fmt.Printf("stake: %q", stake)
	fmt.Printf("currentStake: %q", currentStake)
	if stake.GT(currentStake) {
		marginOfErr := currentStake.Mul(math.LegacyNewDecWithPrec(50, 3)) // 5.0%
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
			patched = true
		} else {
			panic(fmt.Sprintln("current stake is not delgator from slashed validator, and is more than maximum margin of error"))
		}
	}
	// calculate rewards for final period
	rewards = rewards.Add(customCalculateDelegationRewardsBetween(ctx, k, val, startingPeriod, endingPeriod, stake)...)
	return rewards, patched
}

func customCalculateDelegationRewardsBetween(ctx sdk.Context, k *keepers.AppKeepers, val stakingtypes.ValidatorI,
	startingPeriod, endingPeriod uint64, stake math.LegacyDec,
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
	starting, _ := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, sdk.ValAddress(val.GetOperator()), startingPeriod)
	ending, _ := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, sdk.ValAddress(val.GetOperator()), endingPeriod)
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
	historical, _ := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
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
	historical, _ := k.DistrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// calculate the token worth of provided shares, truncated
func CustommTokensFromSharesTruncated(t math.Int, ds math.LegacyDec, shares math.LegacyDec) math.LegacyDec {
	return (shares.MulInt(t)).QuoTruncate(ds)
}
