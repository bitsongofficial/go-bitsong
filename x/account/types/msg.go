package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAccountCreate = "account_create"
)

var _ sdk.Msg = MsgAccountCreate{}

type MsgAccountCreate struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Handle  string         `json:"handle"`
}

func NewMsgAccountCreate(address sdk.AccAddress, handle string) MsgAccountCreate {
	return MsgAccountCreate{
		Address: address,
		Handle:  handle,
	}
}

func (msg MsgAccountCreate) Route() string { return RouterKey }
func (msg MsgAccountCreate) Type() string  { return TypeMsgAccountCreate }

func (msg MsgAccountCreate) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "account address cannot be empty")
	}

	if msg.Handle == "" || len(msg.Handle) < 3 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "account handle must have a length > 3")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAccountCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAccountCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgAccountCreate) String() string {
	return fmt.Sprintf(`Msg Account Create
Address: %s,
Handle: %s`,
		msg.Address.String(), msg.Handle,
	)
}
