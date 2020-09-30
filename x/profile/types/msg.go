package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgProfileCreate = "profile_create"
)

var _ sdk.Msg = MsgProfileCreate{}

type MsgProfileCreate struct {
	Address     sdk.AccAddress `json:"address" yaml:"address"`
	Handle      string         `json:"handle"`
	MetadataURI string         `json:"metadata_uri"`
}

func NewMsgProfileCreate(address sdk.AccAddress, handle, metadataUri string) MsgProfileCreate {
	return MsgProfileCreate{
		Address:     address,
		Handle:      handle,
		MetadataURI: metadataUri,
	}
}

func (msg MsgProfileCreate) Route() string { return RouterKey }
func (msg MsgProfileCreate) Type() string  { return TypeMsgProfileCreate }

func (msg MsgProfileCreate) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile address cannot be empty")
	}

	if msg.Handle == "" || len(msg.Handle) < 3 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile handle must have a length > 3")
	}

	if len(msg.Handle) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "profile metadataURI cannot be more than 256 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgProfileCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgProfileCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgProfileCreate) String() string {
	return fmt.Sprintf(`Msg Profile Create
Address: %s,
Handle: %s
MetadataURI: %s`,
		msg.Address.String(), msg.Handle, msg.MetadataURI,
	)
}
