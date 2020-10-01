package types

type GenesisState struct {
	Releases []Release `json:"releases"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(releases []Release) GenesisState {
	return GenesisState{
		Releases: releases,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Releases: []Release{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Releases {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
