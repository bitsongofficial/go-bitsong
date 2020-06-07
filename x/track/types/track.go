package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenInfo struct {
	Denom     string `json:"denom,omitempty" yaml:"denom,omitempty"`
	Tokenized bool   `json:"tokenized" yaml:"tokenized"`
	Mintable  bool   `json:"mintable" yaml:"mintable"`
}

func NewTokenInfo(denom string) TokenInfo {
	return TokenInfo{
		Denom:     denom,
		Tokenized: true,
		Mintable:  true,
	}
}

type Track struct {
	TrackID   uint64         `json:"track_id" yaml:"track_id"`     // the bitsong track id ****
	Hash      string         `json:"hash" yaml:"hash"`             // the track hash
	Uri       string         `json:"uri" yaml:"uri"`               // bitsong uri for track e.g: bitsong:track:the-show-must-go-on ****
	TrackInfo []byte         `json:"track_info" yaml:"track_info"` // Raw Track Info (see specs)
	TokenInfo TokenInfo      `json:"token_info" yaml:"token_info"` // track token info
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`       // creator of the track
}

func NewTrack(info []byte, creator sdk.AccAddress) *Track {
	trackHash := sha256.Sum256(info)
	trackHashStr := hex.EncodeToString(trackHash[:])

	return &Track{
		TrackInfo: info,
		Hash:      trackHashStr,
		Creator:   creator,
	}
}

func (t *Track) String() string {
	return fmt.Sprintf(`New Track:
Track ID: %d
Uri: %s
Creator: %s`, t.TrackID, t.Uri, t.Creator)
}

func (t *Track) Equals(track Track) bool {
	// TODO
	return true
}

// TODO
func (t *Track) Validate() error {
	if len(t.TrackInfo) == 0 {
		return fmt.Errorf("track info cannot be empty")
	}

	if len(t.TrackInfo) > MaxTrackInfoLength {
		return fmt.Errorf("track info cannot be longer than %d bytes", MaxTrackInfoLength)
	}

	if t.Creator.Empty() {
		return fmt.Errorf("track creator address cannot be empty")
	}

	return nil
}
