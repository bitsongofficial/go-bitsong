package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetTitle defines a SetName message
type MsgSetTitle struct {
	Title  string         `json:"title"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetTitle(title string, value string, owner sdk.AccAddress) MsgSetTitle {
	return MsgSetTitle{
		Title: title,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgSetTitle) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetTitle) Type() string { return "set_title" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetTitle) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	//if len(msg.Name) == 0 || len(msg.Value) == 0 {
	if len(msg.Title) == 0 {
		return sdk.ErrUnknownRequest("Title cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetTitle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetTitle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
