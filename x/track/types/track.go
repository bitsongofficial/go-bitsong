package types

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
	"time"
)

/************************************
 * Track
 ************************************/

type Track struct {
	Title         string         `json:"title" yaml:"title"`
	Address       crypto.Address `json:"address" yaml:"address"`
	Attributes    Attributes     `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Media         TrackMedia     `json:"media" yaml:"media"`
	Rewards       TrackRewards   `json:"rewards" yaml:"rewards"`
	RightsHolders RightsHolders  `json:"rights_holders" yaml:"rights_holders"`
	Totals        TrackTotals    `json:"totals" yaml:"totals"`
	SubmitTime    time.Time      `json:"submit_time" yaml:"submit_time"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewTrack(title string, media TrackMedia, attrs Attributes, rewards TrackRewards,
	rightsHolders RightsHolders, submitTime time.Time, owner sdk.AccAddress) Track {
	return Track{
		Title:         title,
		Rewards:       rewards,
		RightsHolders: rightsHolders,
		Media:         media,
		Attributes:    attrs,
		Owner:         owner,
		SubmitTime:    submitTime,
		Totals: TrackTotals{
			Streams:  0,
			Rewards:  sdk.NewCoin(types.BondDenom, sdk.ZeroInt()),
			Accounts: 0,
		},
	}
}

func (t Track) Validate() error {
	if len(strings.TrimSpace(t.Title)) == 0 {
		return fmt.Errorf("track title cannot be empty")
	}

	if len(t.Title) > MaxTitleLength {
		return fmt.Errorf("track title cannot be longer than %d characters", MaxTitleLength)
	}

	if err := t.Rewards.Validate(); err != nil {
		return err
	}

	if err := t.Media.Validate(); err != nil {
		return err
	}

	if err := t.RightsHolders.Validate(); err != nil {
		return err
	}

	if t.Owner == nil {
		return fmt.Errorf("invalid track owner: %s", t.Owner)
	}

	return nil
}

// nolint
func (t Track) String() string {
	return fmt.Sprintf(`Address: %s
Title: %s
%s
Rewards - %s
Rights Holders
%s
Submit Time: %s
Owner: %s
Attributes
%s
Totals
%s`,
		t.Address.String(), t.Title, t.Media.String(), t.Rewards.String(), t.RightsHolders,
		t.SubmitTime, t.Owner.String(), t.Attributes.String(), t.Totals.String(),
	)
}

func (t Track) Equals(track Track) bool {
	return t.Address.String() == track.Address.String() &&
		t.Title == track.Title &&
		t.Media.Equals(track.Media) &&
		t.Rewards.Equals(track.Rewards) &&
		t.RightsHolders.Equals(track.RightsHolders) &&
		t.Owner.Equals(track.Owner)
}

/************************************
 * Tracks
 ************************************/

// Tracks is an array of track
type Tracks []Track

// nolint
func (t Tracks) String() string {
	out := "Address - Title\n"
	for _, track := range t {
		out += fmt.Sprintf("%s - %s\n",
			track.Address, track.Title)
	}
	return strings.TrimSpace(out)
}
