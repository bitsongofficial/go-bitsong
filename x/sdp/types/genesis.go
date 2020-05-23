package types

// GenesisState - all content state that must be provided at genesis
type GenesisState struct {
	Sdps []Sdp `json:"sdps"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(sdps []Sdp) GenesisState {
	return GenesisState{
		Sdps: sdps,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Sdps: []Sdp{},
	}
}

// ValidateGenesis validates the contents genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Sdps {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
