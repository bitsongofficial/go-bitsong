package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// global pool for distribution
type Pool struct {
	SongsPool sdk.DecCoins `json:"songs_pool"` // pool for community funds yet to be spent
}

// zero fee pool
func InitialPool() Pool {
	return Pool{
		SongsPool: sdk.DecCoins{},
	}
}

// ValidateGenesis validates the pool for a genesis state
func (f Pool) ValidateGenesis() error {
	if f.SongsPool.IsAnyNegative() {
		return fmt.Errorf("negative SongsPool in distribution pool, is %v",
			f.SongsPool)
	}

	return nil
}