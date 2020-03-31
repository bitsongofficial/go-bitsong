package keeper

import (
	"github.com/bitsongofficial/go-bitsong/x/reward"
	cmint "github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/supply"

	rewardTypes "github.com/bitsongofficial/go-bitsong/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	supplyKeeper supply.Keeper
	rewardKeeper reward.Keeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(rk reward.Keeper, sk supply.Keeper) Keeper {
	return Keeper{
		rewardKeeper: rk,
		supplyKeeper: sk,
	}
}

func (k Keeper) AddToRewardPool(ctx sdk.Context, coins sdk.Coins) error {
	rewardPool := k.rewardKeeper.GetRewardPool(ctx)
	rewardPool.Amount = rewardPool.Amount.Add(sdk.NewDecCoinsFromCoins(coins...)...)
	k.rewardKeeper.SetRewardPool(ctx, rewardPool)

	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, cmint.ModuleName, rewardTypes.ModuleName, coins)
}

func (k Keeper) GetRewardPoolSupply(ctx sdk.Context) sdk.Coins {
	return k.rewardKeeper.GetRewardPoolSupply(ctx)
}
