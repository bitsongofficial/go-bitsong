package v021

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	cadancetypes "github.com/bitsongofficial/go-bitsong/x/cadance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v8/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

func CreateV021UpgradeHandler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)

		// resolve store from v018 bug patch
		removeZeroDelegationFromStore(sdkCtx, k)

		// setup vote extension
		consensusParams, err := k.ConsensusParamsKeeper.ParamsStore.Get(sdkCtx)
		if err != nil {
			return nil, err
		}
		// TODO: these are values specific to 1.5s block time. We may need to tune this to bitsongs desired blocktime:
		// Update consensus params in order to safely enable comet pruning
		// consensusParams.Evidence.MaxAgeNumBlocks = 1_000_000
		// consensusParams.Evidence.MaxAgeDuration = time.Second * 1209600
		consensusParams.Abci.VoteExtensionsEnableHeight = 69
		err = k.ConsensusParamsKeeper.ParamsStore.Set(sdkCtx, consensusParams)
		if err != nil {
			return nil, err
		}
		// interchain query params (ICQ)
		setICQParams(sdkCtx, k.ICQKeeper)

		// x/cadance params
		if err := k.CadanceKeeper.SetParams(sdkCtx, cadancetypes.DefaultParams()); err != nil {
			return nil, err
		}
		// x/smart-account
		// Set the authenticator params in the store
		authenticatorParams := k.SmartAccountKeeper.GetParams(sdkCtx)
		authenticatorParams.MaximumUnauthenticatedGas = MaximumUnauthenticatedGas
		authenticatorParams.IsSmartAccountActive = IsSmartAccountActive
		authenticatorParams.CircuitBreakerControllers = append(authenticatorParams.CircuitBreakerControllers, CircuitBreakerController)
		k.SmartAccountKeeper.SetParams(sdkCtx, authenticatorParams)
		// setup wasm client
		// https://github.com/cosmos/ibc-go/blob/main/docs/docs/03-light-clients/04-wasm/03-integration.md
		params := k.IBCKeeper.ClientKeeper.GetParams(sdkCtx)
		params.AllowedClients = append(params.AllowedClients, wasmlctypes.Wasm)
		k.IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)

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

// remove delegations with 0 power to any jailed validators
func removeZeroDelegationFromStore(ctx sdk.Context, k *keepers.AppKeepers) error {
	vals, _ := k.StakingKeeper.GetAllValidators(ctx)
	for _, val := range vals {
		if val.IsJailed() {
			valAddr := sdk.ValAddress(val.OperatorAddress)
			dels, _ := k.StakingKeeper.GetValidatorDelegations(ctx, valAddr)
			for _, del := range dels {
				if del.Shares.IsZero() || val.TokensFromShares(val.DelegatorShares).IsZero() {
					// confirm no rewards remaining
					if rewards, err := k.DistrKeeper.CalculateDelegationRewards(ctx, val, del, uint64(ctx.BlockHeight())); err != nil {
						return err
					} else if !rewards.Empty() {
						return fmt.Errorf("rewards not empty")
					}

					// delete delegation from store
					if err := k.StakingKeeper.RemoveDelegation(ctx, del); err != nil {
						return err
					}
				}
			}

		}
	}
	return nil
}
func setICQParams(ctx sdk.Context, icqKeeper *icqkeeper.Keeper) {
	icqparams := icqtypes.DefaultParams()
	// icqparams.AllowQueries = wasmbinding.GetStargateWhitelistedPaths()
	// Adding SmartContractState query to allowlist
	icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
	//nolint:errcheck
	icqKeeper.SetParams(ctx, icqparams)
}
