package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

// Content messages types and routes
const (
	TypeMsgContentAdd = "content_add"
)

var _ sdk.Msg = MsgContentAdd{}

type MsgContentAdd struct {
	Uri  string `json:"uri" yaml:"uri"`
	Hash string `json:"hash" yaml:"hash"`
	Dao  Dao    `json:"dao" yaml:"dao"`
}

func NewMsgAddContent(uri, hash string, dao Dao) MsgContentAdd {
	return MsgContentAdd{
		Uri:  uri,
		Hash: hash,
		Dao:  dao,
	}
}

func (msg MsgContentAdd) Route() string { return RouterKey }
func (msg MsgContentAdd) Type() string  { return TypeMsgContentAdd }

func (msg MsgContentAdd) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Uri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be empty"))
	}

	if len(msg.Uri) > MaxUriLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be longer than %d characters", MaxUriLength))
	}

	if len(strings.TrimSpace(msg.Hash)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hash cannot be empty"))
	}

	if len(msg.Hash) > MaxHashLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hash cannot be longer than %d characters", MaxUriLength))
	}

	if err := msg.Dao.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgContentAdd) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgContentAdd) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Dao))
	for i, de := range msg.Dao {
		addrs[i] = de.Address
	}

	return addrs
}

func (msg MsgContentAdd) String() string {
	return fmt.Sprintf(`Msg Content Add
Uri: %s
Hash: %s`,
		msg.Uri, msg.Hash,
	)
}

var _ sdk.Msg = MsgContentAction{}

type MsgContentAction struct {
	Uri  string         `json:"uri" yaml:"uri"`
	From sdk.AccAddress `json:"from" yaml:"from"`
}

func NewMsgContentAction(uri string, from sdk.AccAddress) MsgContentAction {
	return MsgContentAction{
		Uri:  uri,
		From: from,
	}
}

func (msg MsgContentAction) Route() string { return RouterKey }
func (msg MsgContentAction) Type() string  { return TypeMsgContentAdd }

func (msg MsgContentAction) ValidateBasic() error {
	if len(strings.TrimSpace(msg.Uri)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("uri cannot be empty"))
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid from: %s", msg.From.String()))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgContentAction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgContentAction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func (msg MsgContentAction) String() string {
	return fmt.Sprintf(`Msg Stream
Uri: %s
From: %s`,
		msg.Uri, msg.From,
	)
}
