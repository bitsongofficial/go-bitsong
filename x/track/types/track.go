package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
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
	TrackID   uint64         `json:"track_id" yaml:"track_id"` // the bitsong track id ****
	Uri       string         `json:"uri" yaml:"uri"`           // bitsong uri for track e.g: bitsong:track:the-show-must-go-on ****
	TrackInfo []byte         `json:"track_info" yaml:"track_info"`
	TokenInfo TokenInfo      `json:"token_info" yaml:"token_info"`
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
		TrackInfo: info,
		Creator:   creator,
	}
}

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
