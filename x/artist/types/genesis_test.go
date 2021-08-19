package types

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {

	testArtists := []Artist{

		{
			ID:          "Test-ID",
			Name:        "Test_Artist",
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{1}, sdk.AddrLen),
		},

		{
			ID:          "Test-ID",
			Name:        strings.Repeat("Test", 64),
			MetadataURI: strings.Repeat("Metadata", 32),
			Creator:     bytes.Repeat([]byte{1}, sdk.AddrLen),
		},

		{
			Name:        "Test_Artist",
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{2}, sdk.AddrLen),
		},

		{
			ID:          "Test-ID",
			Name:        strings.Repeat("Test", 65),
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{3}, sdk.AddrLen),
		},

		{
			ID:          "Test-ID",
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{4}, sdk.AddrLen),
		},

		{
			ID:          "Test-ID",
			Name:        "Test_Artist",
			MetadataURI: strings.Repeat("Metadata", 33),
			Creator:     bytes.Repeat([]byte{5}, sdk.AddrLen),
		},

		{
			ID:          "Test-ID",
			Name:        "Test_Artist",
			MetadataURI: "Test_Metadata",
		},
	}
	tests := []struct {
		desc        string
		genesis     GenesisState
		shouldError bool
	}{
		{desc: "Default Genesis",
			genesis:     DefaultGenesisState(),
			shouldError: false,
		},
		{
			desc:        "Correct Artist data",
			genesis:     NewGenesisState([]Artist{testArtists[0]}),
			shouldError: false,
		},
		{
			desc:        "Correct Artist data",
			genesis:     NewGenesisState([]Artist{testArtists[1]}),
			shouldError: false,
		},
		{
			desc:        "Artist does not have an ID",
			genesis:     NewGenesisState([]Artist{testArtists[2]}),
			shouldError: true,
		},
		{
			desc:        "Artist Name exceeds 256 character limit",
			genesis:     NewGenesisState([]Artist{testArtists[3]}),
			shouldError: true,
		},
		{
			desc:        "Artist doesn't have a name",
			genesis:     NewGenesisState([]Artist{testArtists[4]}),
			shouldError: true,
		},
		{
			desc:        "Metadata exceeds 256 character limit",
			genesis:     NewGenesisState([]Artist{testArtists[5]}),
			shouldError: true,
		},
		{
			desc:        "Artist doesn't have a Creator Address",
			genesis:     NewGenesisState([]Artist{testArtists[6]}),
			shouldError: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			if test.shouldError {
				require.Error(t, ValidateGenesis(test.genesis))
			} else {
				require.NoError(t, ValidateGenesis(test.genesis))
			}
		})
	}
}
