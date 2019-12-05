package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RewardPool struct {
	Amount sdk.DecCoins `json:"reward_pool" yaml:"reward_pool"`
}

func InitialRewardPool() RewardPool {
	return RewardPool{
		Amount: sdk.DecCoins{},
	}
}

func (p RewardPool) ValidateGenesis() error {
	if p.Amount.IsAnyNegative() {
		return fmt.Errorf("negative RewardPool, is %v", p.Amount)
	}

	return nil
}
