package types

import (
	"encoding/binary"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/************************************
 * Track
 ************************************/

// Constants pertaining to an Track object
const (
	MaxTitleLength int = 140
)

type Track struct {
	TrackID uint64         `json:"id"`     // Track ID
	Title   string         `json:"title"`  // Track Title
	Status  TrackStatus    `json:"status"` // Status of the Track {Nil, Verified, Rejected, Failed}
	Owner   sdk.AccAddress `json:"owner"`  // Album owner
}

// TrackKey gets a specific track from the store
func TrackKey(trackID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, trackID)
	return append(TracksKeyPrefix, bz...)
}

func NewTrack(id uint64, title string, owner sdk.AccAddress) Track {
	return Track{
		TrackID: id,
		Title:   title,
		Status:  StatusNil,
		Owner:   owner,
	}
}

// nolint
func (t Track) String() string {
	return fmt.Sprintf(`TrackID %d:
  Title:    %s
  Status:  %s
  Owner:   %s`,
		t.TrackID, t.Title, t.Status.String(), t.Owner.String(),
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
