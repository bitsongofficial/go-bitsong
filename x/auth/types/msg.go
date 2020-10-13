package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgRegisterHandle = "register_handle"
)

var _ sdk.Msg = MsgRegisterHandle{}

type MsgRegisterHandle struct {
	From   sdk.AccAddress `json:"from" yaml:"from"`
	Handle string         `json:"handle" yaml:"handle"`
}

func NewMsgRegisterHandle(from sdk.AccAddress, handle string) MsgRegisterHandle {
	return MsgRegisterHandle{
		From:   from,
		Handle: handle,
	}
}

func (msg MsgRegisterHandle) Route() string { return RouterKey }
func (msg MsgRegisterHandle) Type() string  { return TypeMsgRegisterHandle }

func (msg MsgRegisterHandle) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "from cannot be empty")
	}

	if strings.TrimSpace(msg.Handle) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "handle cannot be empty")
	}

	if len(strings.TrimSpace(msg.Handle)) > 256 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "channel metadataURI cannot be more than 256 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterHandle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterHandle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgRegisterHandle) String() string {
	return fmt.Sprintf(`Msg Register Handle
From: %s,
Handle: %s`,
		msg.From.String(), msg.Handle,
	)
}
