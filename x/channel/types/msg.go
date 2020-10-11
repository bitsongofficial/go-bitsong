package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgChannelCreate = "channel_create"
	TypeMsgChannelEdit   = "channel_edit"
)

var _ sdk.Msg = MsgChannelCreate{}

type MsgChannelCreate struct {
	Owner       sdk.AccAddress `json:"owner" yaml:"owner"`
	Handle      string         `json:"handle"`
	MetadataURI string         `json:"metadata_uri"`
}

func NewMsgChannelCreate(owner sdk.AccAddress, handle, metadataUri string) MsgChannelCreate {
	return MsgChannelCreate{
		Owner:       owner,
		Handle:      handle,
		MetadataURI: metadataUri,
	}
}

func (msg MsgChannelCreate) Route() string { return RouterKey }
func (msg MsgChannelCreate) Type() string  { return TypeMsgChannelCreate }

func (msg MsgChannelCreate) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel owner cannot be empty")
	}

	if msg.Handle == "" || len(msg.Handle) < 3 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel handle must have a length > 3")
	}

	if len(msg.Handle) > 64 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel Handle cannot be more than 256 characters")
	}

	if len(msg.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgChannelCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgChannelCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgChannelCreate) String() string {
	return fmt.Sprintf(`Msg Channel Create
Owner: %s,
Handle: %s
MetadataURI: %s`,
		msg.Owner.String(), msg.Handle, msg.MetadataURI,
	)
}

var _ sdk.Msg = MsgChannelEdit{}

type MsgChannelEdit struct {
	Owner       sdk.AccAddress `json:"owner" yaml:"owner"`
	MetadataURI string         `json:"metadata_uri"`
}

func NewMsgChannelEdit(owner sdk.AccAddress, metadataUri string) MsgChannelEdit {
	return MsgChannelEdit{
		Owner:       owner,
		MetadataURI: metadataUri,
	}
}

func (msg MsgChannelEdit) Route() string { return RouterKey }
func (msg MsgChannelEdit) Type() string  { return TypeMsgChannelEdit }

func (msg MsgChannelEdit) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel owner cannot be empty")
	}

	if len(msg.MetadataURI) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgChannelEdit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgChannelEdit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgChannelEdit) String() string {
	return fmt.Sprintf(`Msg Channel Edit
Owner: %s,
MetadataURI: %s`,
		msg.Owner.String(), msg.MetadataURI,
	)
}
