package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all content state that must be provided at genesis
type GenesisState struct {
	RewardPool
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(amt sdk.Coins) GenesisState {
	return GenesisState{
		RewardPool: RewardPool{
			Amount: amt,
		},
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		RewardPool: RewardPool{
			Amount: sdk.Coins{},
		},
	}
}

// ValidateGenesis validates the contents genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.RewardPool.ValidateGenesis(); err != nil {
		return err
	}

	return nil
}
