package types

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
	for _, fantoken := range gs.FanTokens {
		if err := fantoken.ValidateWithDenom(); err != nil {
			return err
		}
	}

	return nil
}
