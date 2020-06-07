package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	TrackID   uint64         `json:"track_id" yaml:"track_id"` // the bitsong track id ****
	Uri       string         `json:"uri" yaml:"uri"`           // bitsong uri for track e.g: bitsong:track:the-show-must-go-on ****
	TrackInfo []byte         `json:"track_info" yaml:"track_info"`
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`

	/*Title        string         `json:"title" yaml:"title"`                 // title of the track ****
	Artists      []string       `json:"artists" yaml:"artists"`             // the artists who performed the track ****
	Images       Externals      `json:"images" yaml:"images"`               // the track images
	Sources      Externals      `json:"sources" yaml:"sources"`             // the track sources
	Feat         []string       `json:"feat" yaml:"feat"`                   // the feat who performed the track ****
	Producers    []string       `json:"producers" yaml:"producers"`         // the producers who performed the track ****
	Genre        string         `json:"genre" yaml:"genre"`                 // ****
	Mood         string         `json:"mood" yaml:"mood"`                   // ****
	Tags         []string       `json:"tags" yaml:"tags"`                   // ****
	Explicit     bool           `json:"explicit" yaml:"explicit"`           // parental advisory, explicit content tag, as supplied to bitsong by issuer ****
	Label        string         `json:"label" yaml:"label"`                 // ****
	Credits      string         `json:"credits" yaml:"credits"`             // ****
	Copyright    string         `json:"copyright" yaml:"copyright"`         // ****
	PreviewUrl   string         `json:"preview_url" yaml:"preview_url"`     // a link to a 30s preview (mp3 format), can be nil ****
	ExternalIds  Externals      `json:"external_ids" yaml:"external_ids"`   // Known external IDs for the track. eg. key: isrc|ean|upc -> value...
	ExternalUrls Externals      `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	Number       uint           `json:"number" yaml:"number"`               // the track number (usually 1 unless the album consists of more than one disc).
	Duration     uint           `json:"duration" yaml:"duration"`           // the length of the track in milliseconds
	Dao          Dao            `json:"dao" yaml:"dao"`
	*/
	// album
	// Popularity
	// download
	// subscriptionStreaming
}

func NewTrack(info []byte, creator sdk.AccAddress) *Track {
	return &Track{
		//Uri:       uri,
		TrackInfo: info,
		Creator:   creator,
	}
}

/*func NewTrack(title string, artists, feat, producers, tags []string, genre, mood, label,
	credits, copyright, pUrl string, number, duration uint, explicit bool, images, sources, extIds,
	extUrls Externals, dao Dao, creator sdk.AccAddress) (*Track, error) {
	return &Track{
		Title:        title,
		Artists:      artists,
		Images:       images,
		Sources:      sources,
		Feat:         feat,
		Producers:    producers,
		Genre:        genre,
		Mood:         mood,
		Tags:         tags,
		Explicit:     explicit,
		Label:        label,
		Credits:      credits,
		Copyright:    copyright,
		PreviewUrl:   pUrl,
		Number:       number,
		Duration:     duration,
		ExternalIds:  extIds,
		ExternalUrls: extUrls,
		Dao:          dao,
		Creator:      creator,
	}, nil
}*/

func (t *Track) String() string {
	// TODO
	return fmt.Sprintf("Uri: %s", t.Uri)
}

func (t *Track) Equals(track Track) bool {
	// TODO
	return true
}

func (t *Track) Validate() error {
	// TODO

	if len(strings.TrimSpace(t.Uri)) == 0 {
		return fmt.Errorf("title cannot be empty")
	}

	//if len(c.Uri) > MaxUriLength {
	//	return fmt.Errorf("uri cannot be longer than %d characters", MaxUriLength)
	//}

	//if err := t.Dao.Validate(); err != nil {
	//	return fmt.Errorf("%s", err.Error())
	//}

	return nil
}
