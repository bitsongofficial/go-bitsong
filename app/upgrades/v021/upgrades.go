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
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// Run migrations first
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// reapplies v018 patch after removing delegations with 0 power, letting us revert back upstream to cosmos-sdk library
		vals, _ := k.StakingKeeper.GetAllValidators(ctx)
		for _, val := range vals {
			valAddr := sdk.ValAddress(val.OperatorAddress)
			dels, _ := k.StakingKeeper.GetValidatorDelegations(ctx, valAddr)

			for _, del := range dels {
				if del.Shares == math.LegacyZeroDec() {
					// remove delegation from staking store
					if err := k.StakingKeeper.RemoveDelegation(ctx, del); err != nil {
						return nil, err
					}
					// remove reward information from distribution store
					if exists, err := k.DistrKeeper.HasDelegatorStartingInfo(ctx, valAddr, sdk.AccAddress(del.DelegatorAddress)); err != nil || !exists {
						return nil, err
					}
					if err := k.DistrKeeper.DeleteDelegatorStartingInfo(ctx, valAddr, sdk.AccAddress(del.DelegatorAddress)); err != nil {
						return nil, err
					}
				} else {
					// check if we need to patch distribution by manually claiming rewards again
					hasInfo, err := k.DistrKeeper.HasDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(del.GetDelegatorAddr()))
					if !hasInfo {
						return nil, err
					}
					// calculate rewards
					endingPeriod, err := k.DistrKeeper.IncrementValidatorPeriod(ctx, val)
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
