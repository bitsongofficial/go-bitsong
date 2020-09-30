package types

type GenesisState struct {
	Accounts []Account `json:"accounts"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(accounts []Account) GenesisState {
	return GenesisState{
		Accounts: accounts,
	}
}

// DefaultGenesisState - default GenesisState
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: []Account{},
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, item := range data.Accounts {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
