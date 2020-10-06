package types

type GenesisState struct {
	Channels []Channel `json:"channels"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(channels []Channel) GenesisState {
	return GenesisState{
		Channels: channels,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Channels: []Channel{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Channels {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
