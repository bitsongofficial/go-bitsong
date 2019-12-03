package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	// CreateAlbumJSON defines a CreateAlbum msg
	CreateAlbumJSON struct {
		AlbumType            types.AlbumType `json:"album_type" yaml:"album_type"`
		Title                string          `json:"title" yaml:"title"`
		ReleaseDate          string          `json:"release_date" yaml:"release_date"`
		ReleaseDatePrecision string          `json:"release_date_precision" yaml:"release_date_precision"`
	}
)

// ParseCreateAlbumJSON reads and parses a CreateAlbumJSON from a file.
func ParseCreateAlbumJSON(cdc *codec.Codec, albumFile string) (CreateAlbumJSON, error) {
	album := CreateAlbumJSON{}

	payload, err := ioutil.ReadFile(albumFile)
	if err != nil {
		return album, err
	}

	if err := cdc.UnmarshalJSON(payload, &album); err != nil {
		return album, err
	}

	return album, nil
}

type (
	// AlbumVerifyProposalJSON defines a ArtistVerifyProposal with a deposit
	AlbumVerifyProposalJSON struct {
		Title       string    `json:"title" yaml:"title"`
		Description string    `json:"description" yaml:"description"`
		AlbumID     uint64    `json:"id" yaml:"id"`
		Deposit     sdk.Coins `json:"deposit" yaml:"deposit"`
	}
)

// ParseAlbumVerifyProposalJSON reads and parses a ArtistVerifyProposalJSON from a file.
func ParseAlbumVerifyProposalJSON(cdc *codec.Codec, proposalFile string) (AlbumVerifyProposalJSON, error) {
	proposal := AlbumVerifyProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
