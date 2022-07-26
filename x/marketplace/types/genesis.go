package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params Params,
) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	return nil
}
