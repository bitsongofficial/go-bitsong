package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens []FanToken) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

// ValidateGenesis validates the provided token genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	// validate token
	for _, token := range data.Tokens {
		if err := ValidateToken(token); err != nil {
			return err
		}
	}

	// validate token
	for _, coin := range data.BurnedCoins {
		if err := coin.Validate(); err != nil {
			return err
		}
	}
	return nil
}
