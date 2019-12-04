package types

import (
	"fmt"
)

// Track
type Track struct {
	AlbumID  uint64 `json:"album_id" yaml:"album_id"`
	TrackID  uint64 `json:"track_id" yaml:"track_id"`
	Position uint64 `json:"position" yaml:"position"`
}

// NewTrack create a new Track instance
func NewTrack(albumID uint64, trackID uint64, position uint64) Track {
	return Track{
		AlbumID:  albumID,
		TrackID:  trackID,
		Position: position,
	}
}

func (t Track) String() string {
	return fmt.Sprintf("trackID %d added on album %d on position %d", t.TrackID, t.AlbumID, t.Position)
}

// Tracks is a collection of Track objects
type Tracks []Track

func (t Tracks) String() string {
	if len(t) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Tracks on Album %d:", t[0].AlbumID)
	for _, track := range t {
		out += fmt.Sprintf("\n  %d: %d", track.Position, track.TrackID)
	}
	return out
}

// Equals returns whether two tracks are equal.
func (t Track) Equals(comp Track) bool {
	// TODO: compare position?
	return t.TrackID == comp.TrackID && t.AlbumID == comp.AlbumID
}

// Empty returns whether a track is empty.
func (t Track) Empty() bool {
	return t.Equals(Track{})
}
