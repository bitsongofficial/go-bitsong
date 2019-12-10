package types

import (
	"encoding/binary"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/************************************
 * Artist
 ************************************/

// Constants pertaining to an Artist object
const (
	MaxNameLength int = 140
)

type Artist struct {
	ArtistID uint64         `json:"id"`   // Artist ID
	Name     string         `json:"name"` // Artist Name
	Image    ArtistImage    `json:"image" // Artist Image`
	Status   ArtistStatus   `json:"status"` // Status of the Artist {Nil, Verified, Rejected, Failed}
	Owner    sdk.AccAddress `json:"owner"`  // Artist Address
}

// ArtistKey gets a specific artist from the store
func ArtistKey(artistID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, artistID)
	return append(ArtistsKeyPrefix, bz...)
}

// Split keys function; used for iterators
// SplitArtistKey split the artist key and returns the artist id
func SplitArtistKey(key []byte) (artistID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	}

	return binary.LittleEndian.Uint64(key[1:])
}

func NewArtist(id uint64, name string, owner sdk.AccAddress) Artist {
	return Artist{
		ArtistID: id,
		Name:     name,
		Status:   StatusNil,
		Owner:    owner,
	}
}

// nolint
func (a Artist) String() string {
	return fmt.Sprintf(`ArtistID %d:
  Name:    %s
  Status:  %s
  Address:   %s`,
		a.ArtistID, a.Name, a.Status.String(), a.Owner.String(),
	)
}

/************************************
 * Artists
 ************************************/

// Artists is an array of artist
type Artists []Artist

// nolint
func (a Artists) String() string {
	out := "ID - (Status) Name\n"
	for _, artist := range a {
		out += fmt.Sprintf("%d - (%s) %s\n",
			artist.ArtistID, artist.Status, artist.Name)
	}
	return strings.TrimSpace(out)
}
