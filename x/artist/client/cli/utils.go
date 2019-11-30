package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// ArtistVerifyProposalJSON defines a ArtistVerifyProposal with a deposit
	ArtistVerifyProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		ArtistID    uint64    `json:"id" yaml:"id"`
		Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`
	}
)

// ParseArtistVerifyProposalJSON reads and parses a ArtistVerifyProposalJSON from a file.
func ParseArtistVerifyProposalJSON(cdc *codec.Codec, proposalFile string) (ArtistVerifyProposalJSON, error) {
	proposal := ArtistVerifyProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
