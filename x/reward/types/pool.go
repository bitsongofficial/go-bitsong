package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RewardPool struct {
	Amount sdk.Coins `json:"reward_pool" yaml:"reward_pool"`
}

func (p RewardPool) String() string {
	return fmt.Sprintf(`Reward Pool: %s`, p.Amount)
}

func (p RewardPool) ValidateGenesis() error {
	if p.Amount.IsAnyNegative() {
		return fmt.Errorf("negative RewardPool, is %v", p.Amount)
	}

	return nil
}
