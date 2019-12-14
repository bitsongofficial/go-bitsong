package types

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

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
	AlbumID     uint64         `json:"id"`           // Album ID
	Title       string         `json:"title"`        // Album Title
	MetadataURI string         `json:"metadata_uri"` // Metadata uri example: ipfs:QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u
	AlbumType   AlbumType      `json:"album_type"`   // The type of the album: one of 'album', 'single', or 'compilation'.
	Status      AlbumStatus    `json:"status"`       // Status of the Album {Nil, Verified, Rejected, Failed}
	Owner       sdk.AccAddress `json:"owner"`        // Album owner

	SubmitTime     time.Time `json:"submit_time" yaml:"submit_time"`
	TotalDeposit   sdk.Coins `json:"total_deposit" yaml:"total_deposit"`
	DepositEndTime time.Time `json:"deposit_end_time" yaml:"deposit_end_time"`
	VerifiedTime   time.Time `json:"verified_time" yaml:"verified_time"`
}

// AlbumKey gets a specific artist from the store
func AlbumKey(albumID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, albumID)
	return append(AlbumsKeyPrefix, bz...)
}

func NewAlbum(id uint64, title string, albumType AlbumType, metadataUri string, owner sdk.AccAddress, submitTime time.Time) Album {
	return Album{
		AlbumID:      id,
		AlbumType:    albumType,
		Title:        title,
		MetadataURI:  metadataUri,
		Status:       StatusNil,
		Owner:        owner,
		TotalDeposit: sdk.NewCoins(),
		SubmitTime:   submitTime,
	}
}

// nolint
func (a Album) String() string {
	return fmt.Sprintf(`AlbumID %d:
  Type:    %s
  Title:    %s
  Status:  %s
  Address:   %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s`,
		a.AlbumID, a.AlbumType.String(), a.Title, a.Status.String(), a.Owner.String(), a.SubmitTime.String(), a.DepositEndTime.String(), a.TotalDeposit.String(),
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
