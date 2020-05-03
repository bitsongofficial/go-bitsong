package types

// GenesisState - all content state that must be provided at genesis
type GenesisState struct {
	Players []Player `json:"players"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(players []Player) GenesisState {
	return GenesisState{
		Players: players,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Players: []Player{},
	}
}

// ValidateGenesis validates the contents genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Players {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
