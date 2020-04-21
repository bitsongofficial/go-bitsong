package reward

import (
	"fmt"

	"github.com/bitsongofficial/go-bitsong/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/bitsongofficial/go-bitsong/x/reward/keeper"
)

type GenesisState struct {
	RewardPool types.RewardPool `json:"reward_pool" yaml:"reward_pool"`
	RewardTax  sdk.Dec          `json:"reward_tax" yaml:"reward_tax"`
	Rewards    types.Rewards    `json:"rewards" yaml:"rewards"`
}

func NewGenesisState(rewardPool types.RewardPool, rewardTax sdk.Dec) GenesisState {
	return GenesisState{
		RewardPool: rewardPool,
		RewardTax:  rewardTax,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		RewardPool: types.InitialRewardPool(),
		RewardTax:  sdk.NewDecWithPrec(3, 2), // 3%
	}
}

func ValidateGenesis(data GenesisState) error {
	if data.RewardTax.IsNegative() || data.RewardTax.GT(sdk.OneDec()) {
		return fmt.Errorf("mint parameter RewardTax should non-negative and less than one, is %s", data.RewardTax.String())
	}

	return data.RewardPool.ValidateGenesis()
}

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, data GenesisState) {
	var moduleHoldings = sdk.NewDecCoins()

	keeper.SetRewardPool(ctx, data.RewardPool)
	keeper.SetRewardTax(ctx, data.RewardTax)

	for _, reward := range data.Rewards {
		keeper.SetReward(ctx, reward.AccAddr, reward)
	}

	moduleHoldings = moduleHoldings.Add(data.RewardPool.Amount...)
	moduleHoldingsInt, _ := moduleHoldings.TruncateDecimal()

	// check if the module account exists
	moduleAcc := keeper.GetRewardModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	coins := bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if coins.IsZero() {
		if err := bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), moduleHoldingsInt); err != nil {
			panic(err)
		}
		accountKeeper.SetModuleAccount(ctx, moduleAcc)
	}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	rewardPool := keeper.GetRewardPool(ctx)
	rewardTax := keeper.GetRewardTax(ctx)

	return NewGenesisState(
		rewardPool,
		rewardTax,
	)
}
