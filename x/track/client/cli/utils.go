package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// TrackVerifyProposalJSON defines a TrackVerifyProposal with a deposit
	TrackVerifyProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		TrackID     uint64    `json:"id" yaml:"id"`
		Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`
	}
)

// ParseArtistVerifyProposalJSON reads and parses a ArtistVerifyProposalJSON from a file.
func ParseTrackVerifyProposalJSON(cdc *codec.Codec, proposalFile string) (TrackVerifyProposalJSON, error) {
	proposal := TrackVerifyProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
