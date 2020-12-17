package types

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {

	testReleases := []Release{

		{
			ReleaseID:   "Test-ID",
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{1}, sdk.AddrLen),
		},

		{
			ReleaseID:   "Test-ID",
			MetadataURI: strings.Repeat("Metadata", 32),
			Creator:     bytes.Repeat([]byte{2}, sdk.AddrLen),
		},
		{
			ReleaseID: "Test-ID",
			Creator:   bytes.Repeat([]byte{2}, sdk.AddrLen),
		},

		{
			MetadataURI: "Test_Metadata",
			Creator:     bytes.Repeat([]byte{2}, sdk.AddrLen),
		},

		{
			ReleaseID:   "Test-ID",
			MetadataURI: strings.Repeat("Metadata", 33),
			Creator:     bytes.Repeat([]byte{2}, sdk.AddrLen),
		},

		{
			ReleaseID:   "Test-ID",
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
			desc:        "Correct Release data",
			genesis:     NewGenesisState([]Release{testReleases[0]}),
			shouldError: false,
		},
		{
			desc:        "Correct Release data",
			genesis:     NewGenesisState([]Release{testReleases[1]}),
			shouldError: false,
		},
		{
			desc:        "Correct Release data",
			genesis:     NewGenesisState([]Release{testReleases[2]}),
			shouldError: false,
		},
		{
			desc:        "Release ID is missing",
			genesis:     NewGenesisState([]Release{testReleases[3]}),
			shouldError: true,
		},
		{
			desc:        "MetadataURI exceeds 256 characters",
			genesis:     NewGenesisState([]Release{testReleases[4]}),
			shouldError: true,
		},
		{
			desc:        "Creator is missing",
			genesis:     NewGenesisState([]Release{testReleases[5]}),
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
