package v021

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	v020 "github.com/bitsongofficial/go-bitsong/app/upgrades/v020"
	cadancetypes "github.com/bitsongofficial/go-bitsong/x/cadance/types"
	sca "github.com/bitsongofficial/go-bitsong/x/smart-account/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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
					rewardsRaw, patched := v020.CustomCalculateDelegationRewards(sdkCtx, k, val, del, endingPeriod)

					outstanding, err := k.DistrKeeper.GetValidatorOutstandingRewardsCoins(sdkCtx, sdk.ValAddress(del.GetValidatorAddr()))
					if err != nil {
						return nil, err
					}

					if patched {
						err = v020.V018ManualDelegationRewardsPatch(sdkCtx, rewardsRaw, outstanding, k, val, del, endingPeriod)
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
