package types

import (
	"encoding/base64"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgSdpOffer  = "sdp_offer"
	TypeMsgSdpAnswer = "sdp_answer"
)

var _ sdk.Msg = MsgSdp{}

type MsgSdp struct {
	From      sdk.AccAddress `json:"from" yaml:"from"`
	Recipient sdk.AccAddress `json:"recipient" yaml:"recipient"`
	Sdp       string         `json:"sdp" yaml:"sdp"`
	Data      []byte         `json:"data" yaml:"data"`
}

func NewMsgSdp(from, recipient sdk.AccAddress, sdp string, data []byte) MsgSdp {
	return MsgSdp{
		From:      from,
		Recipient: recipient,
		Sdp:       sdp,
		Data:      data,
	}
}

func (msg MsgSdp) Route() string { return RouterKey }
func (msg MsgSdp) Type() string {
	switch msg.Sdp {
	case "answer":
		return TypeMsgSdpAnswer
	default:
		return TypeMsgSdpOffer
	}
}

func (msg MsgSdp) ValidateBasic() error {
	if msg.From.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if msg.Recipient.Empty() {
		return fmt.Errorf(`from address cannot be empty`)
	}

	if len(strings.TrimSpace(msg.Type())) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("type cannot be empty"))
	}

	if msg.Type() != "offer" && msg.Type() != "answer" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("type must be offer or answer"))
	}

	if len(msg.Data) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("data cannot be empty"))
	}

	// max len(msg.Data) >

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSdp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSdp) GetSigners() []sdk.AccAddress {
	switch msg.Sdp {
	case "answer":
		return []sdk.AccAddress{msg.Recipient}
	default:
		return []sdk.AccAddress{msg.From}
	}
}

func (msg MsgSdp) String() string {
	return fmt.Sprintf(`Msg Sdp
From: %s
Recipient: %s
Sdp: %s
Data: %v`,
		msg.From, msg.Recipient, msg.Sdp, base64.StdEncoding.EncodeToString(msg.Data),
	)
}
