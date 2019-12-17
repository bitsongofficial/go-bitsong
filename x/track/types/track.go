package types

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/************************************
 * Track
 ************************************/

// Constants pertaining to an Track object
const (
	MaxTitleLength int = 140
)

// TODO: image, cid, duration
type Track struct {
	TrackID     uint64         `json:"id" yaml:"id"`         // Track ID
	Title       string         `json:"title" yaml:"title"`   // Track Title
	Status      TrackStatus    `json:"status" yaml:"status"` // Status of the Track {Nil, Verified, Rejected, Failed}
	MetadataURI string         `json:"metadata_uri"`         // Metadata uri example: ipfs:QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u
	Owner       sdk.AccAddress `json:"owner" yaml:"owner"`   // Album owner
	TotalPlays  uint64         `json:"total_plays" yaml:"total_plays"`

	SubmitTime     time.Time `json:"submit_time" yaml:"submit_time"`
	TotalDeposit   sdk.Coins `json:"total_deposit" yaml:"total_deposit"`
	DepositEndTime time.Time `json:"deposit_end_time" yaml:"deposit_end_time"`
	VerifiedTime   time.Time `json:"verified_time" yaml:"verified_time"`
}

// TrackKey gets a specific track from the store
func TrackKey(trackID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackID)
	return append(TracksKeyPrefix, bz...)
}

func NewTrack(id uint64, title string, uri string, owner sdk.AccAddress, submitTime time.Time) Track {
	return Track{
		TrackID:      id,
		Title:        title,
		MetadataURI:  uri,
		Status:       StatusNil,
		Owner:        owner,
		TotalPlays:   0,
		TotalDeposit: sdk.NewCoins(),
		SubmitTime:   submitTime,
	}
}

// nolint
func (t Track) String() string {
	return fmt.Sprintf(`TrackID %d:
  Title:    %s
  Metadata: %s
  Status:  %s
  Owner:   %s
  Total Plays: %d
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s`,
		t.TrackID, t.Title, t.MetadataURI, t.Status.String(), t.Owner.String(), t.TotalPlays, t.SubmitTime, t.DepositEndTime, t.TotalDeposit.String(),
	)
}

/************************************
 * Tracks
 ************************************/

// Tracks is an array of track
type Tracks []Track

// nolint
func (t Tracks) String() string {
	out := "ID - (Status) Title\n"
	for _, track := range t {
		out += fmt.Sprintf("%d - (%s) %s\n",
			track.TrackID, track.Status.String(), track.Title)
	}
	return strings.TrimSpace(out)
}
