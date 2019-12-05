package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type Play struct {
	TrackID   uint64         `json:"track_id" yaml:"track_id"`
	AccAddr   sdk.AccAddress `json:"acc_addr" yaml:"acc_addr"`
	Shares    sdk.Dec        `json:"shares" yaml:"shares"`
	Streams   uint64         `json:"streams" yaml:"streams"`
	CreatedAt time.Time      `json:"created_at" yaml:"created_at"`
}

func NewPlay(trackID uint64, accAddr sdk.AccAddress, shares sdk.Dec, streams uint64, createdAt time.Time) Play {
	return Play{
		TrackID:   trackID,
		AccAddr:   accAddr,
		Shares:    shares,
		Streams:   streams,
		CreatedAt: createdAt,
	}
}

func (p Play) String() string {
	return fmt.Sprintf("play added on trackID %d by %s at height %s", p.TrackID, p.AccAddr.String(), p.CreatedAt.String())
}

func (p Play) Validate() error {
	if p.AccAddr.Empty() {
		return fmt.Errorf("invalid acc_addr")
	}

	if p.Shares.Equal(sdk.Dec{}) {
		return fmt.Errorf("invalid shares")
	}

	if p.Streams == 0 {
		return fmt.Errorf("invalid streams")
	}

	if p.CreatedAt.Equal(time.Time{}) {
		return fmt.Errorf("invalid time")
	}

	return nil
}

func (p Play) Equal(comp Play) bool {
	return p.CreatedAt == comp.CreatedAt && p.AccAddr.Equals(comp.AccAddr) && p.CreatedAt.Equal(comp.CreatedAt)
}

// Plays is a collection of Play objects
type Plays []Play

func (p Plays) String() string {
	if len(p) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Plays on Track %d:", p[0].TrackID)
	for _, play := range p {
		out += fmt.Sprintf("\n  %d: %v - %s", play.TrackID, play.CreatedAt, play.AccAddr)
	}
	return out
}

func (p Play) Equals(comp Play) bool {
	return p.TrackID == comp.TrackID && p.AccAddr.String() == comp.AccAddr.String() && p.CreatedAt == comp.CreatedAt
}

// Empty returns whether a track is empty.
func (p Play) Empty() bool {
	return p.Equals(Play{})
}
