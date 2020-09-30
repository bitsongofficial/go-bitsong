package types

type GenesisState struct {
	Profiles []Profile `json:"profiles"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(profiles []Profile) GenesisState {
	return GenesisState{
		Profiles: profiles,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Profiles: []Profile{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Profiles {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
