package types

import (
	"encoding/binary"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/************************************
 * Album
 ************************************/

// Constants pertaining to an Album object
const (
	MaxTitleLength int = 140
)

type Album struct {
	AlbumID              uint64         `json:"id"`                     // Album ID
	AlbumType            AlbumType      `json:"album_type"`             // The type of the album: one of 'album', 'single', or 'compilation'.
	Title                string         `json:"title"`                  // Album Title
	ReleaseDate          string         `json:"release_date"`           // The date the album was first released, for example '1981-12-15'. Depending on the precision, it might be shown as '1981' or '1981-12'.
	ReleaseDatePrecision string         `json:"release_date_precision"` // The precision with which release_date value is known: 'year', 'month', or 'day'.
	Status               AlbumStatus    `json:"status"`                 // Status of the Album {Nil, Verified, Rejected, Failed}
	Owner                sdk.AccAddress `json:"owner"`                  // Album owner
}

// AlbumKey gets a specific artist from the store
func AlbumKey(artistID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, artistID)
	return append(AlbumsKeyPrefix, bz...)
}

func NewAlbum(id uint64, title string, albumType AlbumType, releaseDate string, releasePrecision string, owner sdk.AccAddress) Album {
	return Album{
		AlbumID:              id,
		AlbumType:            albumType,
		Title:                title,
		ReleaseDate:          releaseDate,
		ReleaseDatePrecision: releasePrecision,
		Status:               StatusNil,
		Owner:                owner,
	}
}

// nolint
func (a Album) String() string {
	return fmt.Sprintf(`AlbumID %d:
  Type:    %s
  Title:    %s
  Release Date:    %s
  Release Date Precision:    %s
  Status:  %s
  Owner:   %s`,
		a.AlbumID, a.AlbumType.String(), a.Title, a.ReleaseDate, a.ReleaseDatePrecision, a.Status.String(), a.Owner.String(),
	)
}

/************************************
 * Albums
 ************************************/

// Albums is an array of album
type Albums []Album

// nolint
func (a Albums) String() string {
	out := "ID - (Status) Title\n"
	for _, album := range a {
		out += fmt.Sprintf("%d - (%s) %s\n",
			album.AlbumID, album.Status, album.Title)
	}
	return strings.TrimSpace(out)
}
