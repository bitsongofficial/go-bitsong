package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/****
 * Track Msg
 ***/

// Track messages types and routes
const (
	TypeMsgCreateTrack = "create_track"
	TypeMsgPlayTrack   = "play_track"
)

/****************************************
 * MsgCreateTrack
 ****************************************/

var _ sdk.Msg = MsgCreateTrack{}

// MsgCreateTrack defines CreateTrack message
type MsgCreateTrack struct {
	Title string         `json:"title"` // Track title
	Owner sdk.AccAddress `json:"owner"` // Track owner
}

func NewMsgCreateTrack(title string, owner sdk.AccAddress) MsgCreateTrack {
	return MsgCreateTrack{
		Title: title,
		Owner: owner,
	}
}

//nolint
func (msg MsgCreateTrack) Route() string { return RouterKey }
func (msg MsgCreateTrack) Type() string  { return TypeMsgCreateTrack }

// ValidateBasic
func (msg MsgCreateTrack) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(msg.Title)) == 0 {
		return ErrInvalidTrackTitle(DefaultCodespace, "track title cannot be blank")
	}

	if len(msg.Title) > MaxTitleLength {
		return ErrInvalidTrackTitle(DefaultCodespace, fmt.Sprintf("track title is longer than max length of %d", MaxTitleLength))
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgCreateTrack) String() string {
	return fmt.Sprintf(`Create Track Message:
  Title: %s
  Owner: %s
`, msg.Title, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateTrack) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateTrack) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

/****************************************
 * MsgPlay
 ****************************************/

var _ sdk.Msg = MsgPlay{}

// MsgPlay defines PlayTrack message
type MsgPlay struct {
	TrackID uint64         "json:`track_id` yaml:`track_id`"
	AccAddr sdk.AccAddress `json:"acc_addr"`
}

func NewMsgPlay(trackID uint64, accAddr sdk.AccAddress) MsgPlay {
	return MsgPlay{
		TrackID: trackID,
		AccAddr: accAddr,
	}
}

//nolint
func (msg MsgPlay) Route() string { return RouterKey }
func (msg MsgPlay) Type() string  { return TypeMsgPlayTrack }

// ValidateBasic
func (msg MsgPlay) ValidateBasic() sdk.Error {
	// TODO:
	// - improve check

	if msg.TrackID == 0 {
		return ErrUnknownTrack(DefaultCodespace, "album-id cannot be blank")
	}

	if msg.AccAddr.Empty() {
		return sdk.ErrInvalidAddress(msg.AccAddr.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgPlay) String() string {
	return fmt.Sprintf(`Play Track Message:
  TrackID: %d
  AccAddr: %s
`, msg.TrackID, msg.AccAddr.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgPlay) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgPlay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.AccAddr}
}
