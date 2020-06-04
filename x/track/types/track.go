package types

import (
	"fmt"
	"strings"
)

type TrackType uint8

const (
	// TrackAudio is used when the track type is audio
	TrackAudio TrackType = iota + 1
	// TrackVideo is used when the track type is video
	TrackVideo
)

// TrackTypeMap is used to decode the track type flag value
var TrackTypeMap = map[string]TrackType{
	"audio": TrackAudio,
	"video": TrackVideo,
}

type Track struct {
	Cid   string `json:"cid" yaml:"cid"`     // cid of the track
	Title string `json:"title" yaml:"title"` // title of the track
	// album
	Artists      []Artist  `json:"artists" yaml:"artists"`             // the artists who performed the track
	Number       uint      `json:"number" yaml:"number"`               // the track number (usually 1 unless the album consists of more than one disc).
	Duration     uint      `json:"duration" yaml:"duration"`           // the length of the track in milliseconds
	Explicit     bool      `json:"explicit" yaml:"explicit"`           // parental advisory, explicit content tag, as supplied to bitsong by issuer
	ExternalIds  Externals `json:"external_ids" yaml:"external_ids"`   // Known external IDs for the track. eg. key: isrc|ean|upc -> value...
	ExternalUrls Externals `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	// Popularity
	PreviewUrl string `json:"preview_url" yaml:"preview_url"` // a link to a 30s preview (mp3 format), can be nil
	// Uri string `json:"uri" yaml:"uri"` // the bitsong uri for the artist e.g.: bitsong:artist:zmsdksd394394
	// download
	// subscriptionStreaming
	Dao Dao `json:"dao" yaml:"dao"`
}

func NewTrack(title string, artists []Artist, number, duration uint, explicit bool, extIds, extUrls Externals, pUrl string, dao Dao) (*Track, error) {
	return &Track{
		Artists:      artists,
		Number:       number,
		Duration:     duration,
		Explicit:     explicit,
		ExternalIds:  extIds,
		ExternalUrls: extUrls,
		Title:        title,
		PreviewUrl:   pUrl,
		Dao:          dao,
	}, nil
}

func (t *Track) String() string {
	// TODO
	return fmt.Sprintf("Title: %s", t.Title)
}

func (t *Track) Equals(track Track) bool {
	// TODO
	return true
}

func (t *Track) Validate() error {
	// TODO

	if len(strings.TrimSpace(t.Title)) == 0 {
		return fmt.Errorf("title cannot be empty")
	}

	//if len(c.Uri) > MaxUriLength {
	//	return fmt.Errorf("uri cannot be longer than %d characters", MaxUriLength)
	//}

	if err := t.Dao.Validate(); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
