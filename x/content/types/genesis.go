package types

// GenesisState - all content state that must be provided at genesis
type GenesisState struct {
	Contents []Content `json:"contents"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(contents []Content) GenesisState {
	return GenesisState{
		Contents: contents,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Contents: []Content{},
	}
}

// ValidateGenesis validates the contents genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Contents {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
