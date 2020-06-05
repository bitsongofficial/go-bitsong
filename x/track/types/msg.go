package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

// Content messages types and routes
const (
	TypeMsgTrackAdd = "track_add"
)

var _ sdk.Msg = MsgTrackAdd{}

type MsgTrackAdd struct {
	Title      string   `json:"title" yaml:"title"`         // title of the track
	Artists    []string `json:"artists" yaml:"artists"`     // artists of the track
	Feat       []string `json:"feat" yaml:"feat"`           // track feat
	Producers  []string `json:"producers" yaml:"producers"` // track producers
	Genre      string   `json:"genre" yaml:"genre"`
	Mood       string   `json:"mood" yaml:"mood"`
	Tags       []string `json:"tags" yaml:"tags"`
	Label      string   `json:"label" yaml:"label"`
	Credits    string   `json:"credits" yaml:"credits"`
	Copyright  string   `json:"copyright" yaml:"copyright"`
	PreviewUrl string   `json:"preview_url" yaml:"preview_url"` // a link to a 30s preview (mp3 format), can be nil

	Number       uint              `json:"number" yaml:"number"`               // the track number (usually 1 unless the album consists of more than one disc).
	Duration     uint              `json:"duration" yaml:"duration"`           // the length of the track in milliseconds
	Explicit     bool              `json:"explicit" yaml:"explicit"`           // parental advisory, explicit content tag, as supplied to bitsong by issuer
	ExternalIds  map[string]string `json:"external_ids" yaml:"external_ids"`   // Known external IDs for the track. eg. key: isrc|ean|upc -> value...
	ExternalUrls map[string]string `json:"external_urls" yaml:"external_urls"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...

	Dao Dao `json:"dao" yaml:"dao"`
}

func NewMsgTrackAdd(title string, artists, feat, producers, tags []string, genre, mood,
	label, credits, copyright, pUrl string, number, duration uint, explicit bool,
	extIds, extUrls map[string]string, dao Dao) MsgTrackAdd {
	return MsgTrackAdd{
		Title:        title,
		Artists:      artists,
		Producers:    producers,
		Feat:         feat,
		Genre:        genre,
		Mood:         mood,
		Tags:         tags,
		Label:        label,
		Credits:      credits,
		Copyright:    copyright,
		PreviewUrl:   pUrl,
		Number:       number,
		Duration:     duration,
		Explicit:     explicit,
		ExternalIds:  extIds,
		ExternalUrls: extUrls,
		Dao:          dao,
	}
}

func (msg MsgTrackAdd) Route() string { return RouterKey }
func (msg MsgTrackAdd) Type() string  { return TypeMsgTrackAdd }

func (msg MsgTrackAdd) ValidateBasic() error {
	// TODO:

	if len(strings.TrimSpace(msg.Title)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("title cannot be empty"))
	}

	if err := msg.Dao.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTrackAdd) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTrackAdd) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Dao))
	for i, de := range msg.Dao {
		addrs[i] = de.Address
	}

	return addrs
}

func (msg MsgTrackAdd) String() string {
	// TODO
	return fmt.Sprintf(`Msg Track Add
Title: %s`,
		msg.Title,
	)
}
