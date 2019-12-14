package types

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

/************************************
 * Artist
 ************************************/

// Constants pertaining to an Artist object
const (
	MaxNameLength int = 140
)

type Artist struct {
	ArtistID    uint64         `json:"id"`           // Artist ID
	Name        string         `json:"name"`         // Artist Name
	MetadataURI string         `json:"metadata_uri"` // Metadata uri example: ipfs:QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u
	Status      ArtistStatus   `json:"status"`       // Status of the Artist {Nil, Verified, Rejected, Failed}
	Owner       sdk.AccAddress `json:"owner"`        // Artist Address

	SubmitTime     time.Time `json:"submit_time" yaml:"submit_time"`
	TotalDeposit   sdk.Coins `json:"total_deposit" yaml:"total_deposit"`
	DepositEndTime time.Time `json:"deposit_end_time" yaml:"deposit_end_time"`
	VerifiedTime   time.Time `json:"verified_time" yaml:"verified_time"`
}

// ArtistKey gets a specific artist from the store
func ArtistKey(artistID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, artistID)
	return append(ArtistsKeyPrefix, bz...)
}

func NewArtist(id uint64, name string, uri string, owner sdk.AccAddress, submitTime time.Time) Artist {
	// TODO: first status NIL, then when addDesposit change status to StatusDepositPeriod
	return Artist{
		ArtistID:     id,
		Name:         name,
		MetadataURI:  uri,
		Status:       StatusNil,
		Owner:        owner,
		TotalDeposit: sdk.NewCoins(),
		SubmitTime:   submitTime,
	}
}

// nolint
func (a Artist) String() string {
	return fmt.Sprintf(`ArtistID %d:
  Name:    %s
  Metadata: %s
  Status:  %s
  Address:   %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s`,
		a.ArtistID, a.Name, a.MetadataURI, a.Status.String(), a.Owner.String(), a.SubmitTime, a.DepositEndTime, a.TotalDeposit,
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
