package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
	"time"
)

type Track struct {
	TrackID string `json:"track_id" yaml:"track_id"` // the bitsong unique track id ****
	//Hash           string         `json:"hash" yaml:"hash"`             // the track hash
	//Uri            string         `json:"uri" yaml:"uri"`               // bitsong uri for track e.g: bitsong:track:<track_id> ****
	TrackInfo      []byte         `json:"track_info" yaml:"track_info"` // Raw Track Info (see specs)
	Creator        sdk.AccAddress `json:"creator" yaml:"creator"`       // creator of the track
	Provider       sdk.AccAddress `json:"provider" yaml:"provider"`
	TotalShares    sdk.Coin       `json:"total_shares" yaml:"total_shares"`
	TotalDownloads sdk.Int        `json:"total_downloads" yaml:"total_downloads"`
	TotalStreams   sdk.Int        `json:"total_streams" yaml:"total_streams"`
	StreamPrice    sdk.Coin       `json:"stream_price" yaml:"stream_price"`
	DownloadPrice  sdk.Coin       `json:"download_price" yaml:"download_price"`
	StartTime      time.Time      `json:"start_time" yaml:"start_time"`
	EndTime        time.Time      `json:"end_time" yaml:"end_time"`
}

//func NewTrack(id string, info []byte, creator, provider sdk.AccAddress, streamPrice, downloadPrice sdk.Coin) *Track {
func NewTrack(id string, creator, provider sdk.AccAddress, streamPrice, downloadPrice sdk.Coin) *Track {
	//trackHash := sha256.Sum256(info)
	//trackHashStr := hex.EncodeToString(trackHash[:])

	// TODO: init start and end time

	return &Track{
		TrackID: id,
		//TrackInfo: info,
		//Hash:      trackHashStr,
		Creator:  creator,
		Provider: provider,
		TotalShares: sdk.Coin{
			Denom:  getDenom(id),
			Amount: sdk.ZeroInt(),
		},
		TotalDownloads: sdk.ZeroInt(),
		TotalStreams:   sdk.ZeroInt(),
		StreamPrice:    streamPrice,
		DownloadPrice:  downloadPrice,
		StartTime:      time.Now(),
	}
}

func getDenom(trackID string) string {
	denomID := strings.Replace(trackID, "-", "", -1)
	// TODO: cosmos-sdk v0.39 accept max 15chars, fix is applied to v0.40
	// TODO: add security checks
	return fmt.Sprintf(`btrack%s`, denomID[0:10])
}

func (t *Track) ToCoinDenom() string {
	return getDenom(t.TrackID)
}

func (t *Track) String() string {
	return fmt.Sprintf(`New Track:
Track ID: %s
Creator: %s`, t.TrackID, t.Creator)
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
