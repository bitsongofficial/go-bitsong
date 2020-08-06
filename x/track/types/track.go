package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Track struct {
	TrackID   string         `json:"track_id" yaml:"track_id"`     // the bitsong unique track id ****
	Hash      string         `json:"hash" yaml:"hash"`             // the track hash
	Uri       string         `json:"uri" yaml:"uri"`               // bitsong uri for track e.g: bitsong:track:<track_id> ****
	TrackInfo []byte         `json:"track_info" yaml:"track_info"` // Raw Track Info (see specs)
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`       // creator of the track
}

func NewTrack(id string, info []byte, creator sdk.AccAddress) *Track {
	trackHash := sha256.Sum256(info)
	trackHashStr := hex.EncodeToString(trackHash[:])

	return &Track{
		TrackID:   id,
		TrackInfo: info,
		Hash:      trackHashStr,
		Creator:   creator,
	}
}

func (t *Track) ToCoinDenom() string {
	denomID := strings.Replace(t.TrackID, "-", "", -1)
	// TODO: cosmos-sdk v0.39 accept max 15chars, fix is applied to v0.40
	// TODO: add security checks
	return fmt.Sprintf(`btrack%s`, denomID[0:10])
}

func (t *Track) String() string {
	return fmt.Sprintf(`New Track:
Track ID: %s
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
