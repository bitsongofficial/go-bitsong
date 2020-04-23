package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		genesis GenesisState
		expErr  error
	}{
		{
			name:    "Zero last track id and tracks nil",
			genesis: NewGenesisState(0, nil),
			expErr:  fmt.Errorf("starting track id must be > 0"),
		},
		{
			name:    "One last track id and tracks nil",
			genesis: NewGenesisState(1, nil),
			expErr:  nil,
		},
		{
			name:    "10000 last track id and tracks nil",
			genesis: NewGenesisState(10000, nil),
			expErr:  nil,
		},
		{
			name:    "One last track id and tracks empty",
			genesis: NewGenesisState(1, Tracks{}),
			expErr:  nil,
		},
		{
			name:    "One last track id and mockTrack",
			genesis: NewGenesisState(1, Tracks{mockTrack}),
			expErr:  nil,
		},
		{
			name:    "One last track id and mockTrackOwnerNil",
			genesis: NewGenesisState(1, Tracks{mockTrackOwnerNil}),
			expErr:  fmt.Errorf("invalid track owner: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expErr, ValidateGenesis(test.genesis))
		})
	}

}

func TestDefaultGenesisState(t *testing.T) {
	genesisState := GenesisState{
		LastTrackID: 1,
		Tracks:      Tracks{},
	}

	require.Equal(t, genesisState, DefaultGenesisState())
}

func TestNewGenesisState(t *testing.T) {
	lastTrackID := uint64(1)
	tracks := Tracks{}

	genesisState := GenesisState{
		LastTrackID: lastTrackID,
		Tracks:      tracks,
	}

	require.Equal(t, genesisState, NewGenesisState(lastTrackID, tracks))
}
