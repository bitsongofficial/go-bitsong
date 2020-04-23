package types

import "fmt"

// GenesisState - all track state that must be provided at genesis
type GenesisState struct {
	LastTrackID uint64 `json:"last_track_id"`
	Tracks      Tracks `json:"tracks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(lastTrackID uint64, tracks Tracks) GenesisState {
	return GenesisState{
		LastTrackID: lastTrackID,
		Tracks:      tracks,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		LastTrackID: 1,
		Tracks:      Tracks{},
	}
}

// ValidateGenesis validates the track genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.LastTrackID < 1 {
		return fmt.Errorf("starting track id must be > 0")
	}

	for _, item := range data.Tracks {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}
