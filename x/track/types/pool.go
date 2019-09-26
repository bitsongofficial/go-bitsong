package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// global pool for distribution
type Pool struct {
	Rewards sdk.Int `json:"rewards"` // pool for track funds yet to be send
}

// zero fee pool
func InitialPool() Pool {
	return Pool{
		Rewards: sdk.Int{},
	}
}

// ValidateGenesis validates the pool for a genesis state
func (f Pool) ValidateGenesis() error {
	if f.Rewards.IsNegative() {
		return fmt.Errorf("negative Rewards in track pool, is %v",
			f.Rewards)
	}

	return nil
}
