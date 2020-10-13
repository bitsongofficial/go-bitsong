package types

type GenesisState struct {
	Accounts []BitSongAccount `json:"accounts" yaml:"accounts"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(accs []BitSongAccount) GenesisState {
	return GenesisState{
		Accounts: accs,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: []BitSongAccount{},
	}
}

func ValidateGenesis(data GenesisState) error {
	/*for _, item := range data.Accounts {

	}*/

	return nil
}
