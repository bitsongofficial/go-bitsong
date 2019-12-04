package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Play struct {
	TrackID   uint64         `json:"track_id" yaml:"track_id"`
	AccAddr   sdk.AccAddress `json:"acc_addr" yaml:"acc_addr"`
	CreatedAt sdk.Int        `json:"created_at" yaml:"created_at"`
}

func NewPlay(trackID uint64, accAddr sdk.AccAddress, createdAt sdk.Int) Play {
	return Play{
		TrackID:   trackID,
		AccAddr:   accAddr,
		CreatedAt: createdAt,
	}
}

func (p Play) String() string {
	return fmt.Sprintf("play added on trackID %d by %s at height %s", p.TrackID, p.AccAddr.String(), p.CreatedAt.String())
}

func (p Play) Validate() error {
	if p.CreatedAt.Equal(sdk.ZeroInt()) {
		return fmt.Errorf("invalid height")
	}

	if p.AccAddr.Empty() {
		return fmt.Errorf("invalid acc_addr")
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
		out += fmt.Sprintf("\n  %d: %d - %s", play.TrackID, play.CreatedAt, play.AccAddr)
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
