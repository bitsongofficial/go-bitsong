package types

// GenesisState - all tracks state that must be provided at genesis
type GenesisState struct {
	Tracks []Track `json:"tracks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(tracks []Track) GenesisState {
	return GenesisState{
		Tracks: tracks,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Tracks: []Track{},
	}
}

// ValidateGenesis validates the tracks genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Tracks {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
