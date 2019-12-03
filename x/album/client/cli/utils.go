package cli

import (
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
