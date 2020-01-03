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
	MaxTitleLength       int = 140
	MaxDescriptionLength int = 500
	MaxCopyrightLength   int = 500
)

// TODO: image, cid, duration
type Track struct {
	TrackID     uint64      `json:"id" yaml:"id"`       // Track ID
	Title       string      `json:"title" yaml:"title"` // Track Title
	Description string      `json:"description" yaml:"description"`
	Status      TrackStatus `json:"status" yaml:"status"` // Status of the Track {Nil, Verified, Rejected, Failed}
	Audio       string      `json:"audio" yaml:"audio"`
	Image       string      `json:"image" yaml:"image"`
	Duration    string      `json:"duration" yaml:"duration"`

	Owner        sdk.AccAddress `json:"owner" yaml:"owner"` // Album owner
	TotalPlays   uint64         `json:"total_plays" yaml:"total_plays"`
	TotalRewards sdk.Coins      `json:"total_rewards" yaml:"total_rewards"`

	Hidden    bool   `json:"hidden" yaml:"hidden"`
	Explicit  bool   `json:"explicit" yaml:"explicit"`
	Genre     string `json:"genre" yaml:"genre"`
	Mood      string `json:"mood" yaml:"mood"`
	Artists   string `json:"artists" yaml:"artists"`
	Featuring string `json:"featuring" yaml:"featuring"`
	Producers string `json:"producers" yaml:"producers"`
	Copyright string `json:"copyright" yaml:"copyright"`

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

func NewTrack(id uint64, title, audio, image, duration string, hidden bool, explicit bool, genre, mood, artists, featuring, producers, description, copyright string, owner sdk.AccAddress, submitTime time.Time) Track {
	return Track{
		TrackID:      id,
		Title:        title,
		Description:  description,
		Status:       StatusNil,
		Audio:        audio,
		Image:        image,
		Duration:     duration,
		Genre:        genre,
		Mood:         mood,
		Artists:      artists,
		Featuring:    featuring,
		Producers:    producers,
		Copyright:    copyright,
		Owner:        owner,
		TotalPlays:   0,
		TotalRewards: sdk.NewCoins(),
		TotalDeposit: sdk.NewCoins(),
		Hidden:       hidden,
		Explicit:     explicit,
		SubmitTime:   submitTime,
	}
}

// nolint
func (t Track) String() string {
	return fmt.Sprintf(`TrackID %d:
  Title:    %s
  Description: %s
  Status:  %s
  Audio: %s
  Image: %s
  Duration: %s
  Owner:   %s
  Total Plays: %d
  Total Rewards: %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s`,
		t.TrackID, t.Title, t.Description, t.Audio, t.Image, t.Duration, t.Status.String(), t.Owner.String(), t.TotalPlays, t.TotalRewards.String(), t.SubmitTime, t.DepositEndTime, t.TotalDeposit.String(),
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
