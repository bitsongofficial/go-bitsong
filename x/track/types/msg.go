package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

/****
 * Track Msg
 ***/

// Track messages types and routes
const (
	TypeMsgCreate = "create"
)

/****************************************
 * MsgCreate
 ****************************************/

var _ sdk.Msg = MsgCreate{}

// MsgCreateTrack defines Create message
type MsgCreate struct {
	Title         string         `json:"title" yaml:"title"`
	Attributes    Attributes     `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Media         TrackMedia     `json:"media,omitempty" yaml:"media,omitempty"`
	Rewards       TrackRewards   `json:"rewards,omitempty" yaml:"rewards,omitempty"`
	RightsHolders RightsHolders  `json:"rights_holders" yaml:"rights_holders"`
	Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewMsgCreate(title string, attrs Attributes, media TrackMedia, rewards TrackRewards, rightsHolders RightsHolders, owner sdk.AccAddress) MsgCreate {
	return MsgCreate{
		Title:         title,
		Attributes:    attrs,
		Media:         media,
		Rewards:       rewards,
		RightsHolders: rightsHolders,
		Owner:         owner,
	}
}

//nolint
func (msg MsgCreate) Route() string { return RouterKey }
func (msg MsgCreate) Type() string  { return TypeMsgCreate }

// ValidateBasic
func (msg MsgCreate) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Title)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "title cannot be blank")
	}

	if len(msg.Title) > MaxTitleLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("track title is longer than max length of %d", MaxTitleLength))
	}

	if len(msg.Attributes) > MaxAttributesLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("attributes data cannot be longer than %d fields", MaxAttributesLength))
	}

	for key, value := range msg.Attributes {
		if len(value) > MaxAttributesLength {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("attributes value data cannot be longer than %d. %s exceeds the limit", MaxAttributesValueLength, key))
		}
	}

	if err := msg.Media.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := msg.Rewards.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := msg.RightsHolders.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid owner: %s", msg.Owner.String()))
	}

	return nil
}

// String MsgCreate
func (msg MsgCreate) String() string {
	return fmt.Sprintf(`Create Message:
  Title: %s
  Attributes
  %s
  Media: %s
  Rights Holders
  %s
  Owner: %s
`, msg.Title, msg.Attributes.String(), msg.Media.String(), msg.RightsHolders.String(), msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
