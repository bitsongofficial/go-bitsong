package types

import "fmt"

// DefaultGenesisState returns the default genesis state for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		FanTokens: []FanToken{},
	}
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, fantokens []FanToken) GenesisState {
	return GenesisState{
		Params:    params,
		FanTokens: fantokens,
	}
}

// Validate validates the provided token genesis state to ensure the
// expected invariants holds.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// validate fantoken
	for _, token := range gs.FanTokens {
		if err := ValidateFanTokenWithDenom(&token); err != nil {
			return err
		}

		if !token.Mintable && token.Authority != "" {
			return fmt.Errorf("invalid authority")
		}

		if token.Mintable && token.Authority == "" {
			return fmt.Errorf("invalid authority")
		}

		if token.MetaData.Symbol == "" {
			return fmt.Errorf("invalid symbol")
		}

		if token.MetaData.Name == "" {
			return fmt.Errorf("invalid name")
		}
	}

	// validate token
	for _, coin := range gs.BurnedCoins {
		if err := coin.Validate(); err != nil {
			return err
		}
	}
	return nil
}
