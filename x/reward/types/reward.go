package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Reward struct {
	AccAddr      sdk.AccAddress `json:"acc_addr" yaml:"acc_addr"`
	TotalRewards sdk.Coins      `json:"total_rewards" yaml:"total_rewards"`
}

func NewReward(accAddr sdk.AccAddress) Reward {
	return Reward{
		AccAddr:      accAddr,
		TotalRewards: sdk.Coins{sdk.NewCoin("ubtsg", sdk.NewInt(0))},
	}
}

type Rewards []Reward

func (r Rewards) String() string {
	if len(r) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Rewards")
	for _, reward := range r {
		out += fmt.Sprintf("\n  %s: %s", reward.AccAddr.String(), reward.TotalRewards.String())
	}
	return out
}
