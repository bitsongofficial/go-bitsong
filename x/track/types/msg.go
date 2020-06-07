package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Content messages types and routes
const (
	TypeMsgTrackAdd = "track_add"
)

var _ sdk.Msg = MsgTrackAdd{}

type MsgTrackAdd struct {
	TrackInfo []byte         `json:"track_info" yaml:"track_info"`
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`

	/*Title        string         `json:"title" yaml:"title"`     // title of the track
	Artists      []string       `json:"artists" yaml:"artists"` // artists of the track
	Images       Externals      `json:"images" yaml:"images"`
	Sources      Externals      `json:"sources" yaml:"sources"`
	Feat         []string       `json:"feat,omitempty" yaml:"feat,omitempty"`           // track feat
	Producers    []string       `json:"producers,omitempty" yaml:"producers,omitempty"` // track producers
	Genre        string         `json:"genre" yaml:"genre"`
	Mood         string         `json:"mood" yaml:"mood"`
	Tags         []string       `json:"tags,omitempty" yaml:"tags,omitempty"`
	Label        string         `json:"label" yaml:"label"`
	Credits      string         `json:"credits" yaml:"credits"`
	Copyright    string         `json:"copyright" yaml:"copyright"`
	PreviewUrl   string         `json:"preview_url" yaml:"preview_url"`                         // a link to a 30s preview (mp3 format), can be nil
	Number       uint           `json:"number" yaml:"number"`                                   // the track number (usually 1 unless the album consists of more than one disc).
	Duration     uint           `json:"duration" yaml:"duration"`                               // the length of the track in milliseconds
	Explicit     bool           `json:"explicit" yaml:"explicit"`                               // parental advisory, explicit content tag, as supplied to bitsong by issuer
	ExternalIds  Externals      `json:"external_ids,omitempty" yaml:"external_ids,omitempty"`   // Known external IDs for the track. eg. key: isrc|ean|upc -> value...
	ExternalUrls Externals      `json:"external_urls,omitempty" yaml:"external_urls,omitempty"` // known external URLs for this artist eg. key: spotify|youtube|soundcloud -> value...
	Dao          Dao            `json:"dao" yaml:"dao"`
	*/
}

func NewMsgTrackAdd(info []byte, creator sdk.AccAddress) MsgTrackAdd {
	return MsgTrackAdd{
		TrackInfo: info,
		Creator:   creator,
	}
}

/*func NewMsgTrackAdd(title string, artists, feat, producers, tags []string, genre, mood,
	label, credits, copyright, pUrl string, number, duration uint, explicit bool,
	images, sources, extIds, extUrls Externals, dao Dao, creator sdk.AccAddress) MsgTrackAdd {
	return MsgTrackAdd{
		Title:        title,
		Artists:      artists,
		Images:       images,
		Sources:      sources,
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
		Creator:      creator,
	}
}*/

func (msg MsgTrackAdd) Route() string { return RouterKey }
func (msg MsgTrackAdd) Type() string  { return TypeMsgTrackAdd }

func (msg MsgTrackAdd) ValidateBasic() error {
	// TODO:

	/*if len(strings.TrimSpace(msg.Title)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("title cannot be empty"))
	}

	if err := msg.Dao.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	*/
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTrackAdd) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTrackAdd) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
	/*addrs := make([]sdk.AccAddress, len(msg.Dao))
	for i, de := range msg.Dao {
		addrs[i] = de.Address
	}

	return addrs*/
}

func (msg MsgTrackAdd) String() string {
	// TODO
	return fmt.Sprintf(`Msg Track Add
Title: %s`,
		msg.Creator,
	)
}
