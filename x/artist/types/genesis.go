package types

type GenesisState struct {
	Artists []Artist `json:"artists"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(artists []Artist) GenesisState {
	return GenesisState{
		Artists: artists,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Artists: []Artist{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Artists {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
