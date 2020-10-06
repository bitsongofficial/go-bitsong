package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgChannelCreate = "channel_create"
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

	if len(msg.Handle) > 256 {
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
