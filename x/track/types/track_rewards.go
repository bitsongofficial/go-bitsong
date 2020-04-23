package types

import "fmt"

type TrackRewards struct {
	Users     uint `json:"users" yaml:"users"`
	Playlists uint `json:"playlists" yaml:"playlists"`
}

func (tr TrackRewards) Sum() uint {
	return tr.Users + tr.Playlists
}

func (tr TrackRewards) Validate() error {
	if tr.Sum() > 100 {
		return fmt.Errorf("track rewards cannot be more than 100 percent")
	}

	return nil
}

func (tr TrackRewards) Equals(trackRewards TrackRewards) bool {
	return tr.Users == trackRewards.Users && tr.Playlists == trackRewards.Playlists
}

func (tr TrackRewards) String() string {
	return fmt.Sprintf(`Users: %d, Playlists: %d`, tr.Users, tr.Playlists)
}
