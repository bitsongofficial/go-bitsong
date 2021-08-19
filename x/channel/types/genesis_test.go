package types

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {

	testChannels := []Channel{

		{
			Owner:       bytes.Repeat([]byte{1}, sdk.AddrLen),
			Handle:      "Test_Handle",
			MetadataURI: "Test_Metadata",
		},

		{
			Owner:       bytes.Repeat([]byte{2}, sdk.AddrLen),
			Handle:      "THM",
			MetadataURI: "Test_Metadata",
		},

		{
			Owner:       bytes.Repeat([]byte{3}, sdk.AddrLen),
			Handle:      "Test_Handle",
			MetadataURI: strings.Repeat("Metadata", 32),
		},

		{
			Handle:      "Test_Handle",
			MetadataURI: "Test_Metadata",
		},

		{
			Owner:       bytes.Repeat([]byte{4}, sdk.AddrLen),
			MetadataURI: "Test_Metadata",
		},

		{
			Owner:       bytes.Repeat([]byte{5}, sdk.AddrLen),
			Handle:      "TH",
			MetadataURI: "Test_Metadata",
		},

		{
			Owner:       bytes.Repeat([]byte{6}, sdk.AddrLen),
			Handle:      "Test_Handle",
			MetadataURI: strings.Repeat("Metadata", 33),
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
			desc:        "Correct Channel data",
			genesis:     NewGenesisState([]Channel{testChannels[0]}),
			shouldError: false,
		},
		{
			desc:        "Correct Channel data",
			genesis:     NewGenesisState([]Channel{testChannels[1]}),
			shouldError: false,
		},
		{
			desc:        "Correct Channel data",
			genesis:     NewGenesisState([]Channel{testChannels[2]}),
			shouldError: false,
		},
		{
			desc:        "Owner is missing",
			genesis:     NewGenesisState([]Channel{testChannels[3]}),
			shouldError: true,
		},
		{
			desc:        "Handle is missing",
			genesis:     NewGenesisState([]Channel{testChannels[4]}),
			shouldError: true,
		},
		{
			desc:        "Handle length < 3",
			genesis:     NewGenesisState([]Channel{testChannels[5]}),
			shouldError: true,
		},
		{
			desc:        "MetadataURI exceeds 256 characters",
			genesis:     NewGenesisState([]Channel{testChannels[6]}),
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
