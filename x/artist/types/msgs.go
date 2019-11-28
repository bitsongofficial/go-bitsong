package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Artist message types and routes
const (
	TypeMsgCreateArtist = "create_artist"
)

var _ sdk.Msg = MsgCreateArtist{}

// MsgCreateArtist
type MsgCreateArtist struct {
	Meta  Meta           `json:"meta" yaml:"meta"`   // Meta information about artist
	Owner sdk.AccAddress `json:"owner" yaml:"owner"` //  Address of the owner
}

func NewMsgCreateArtist(meta Meta, owner sdk.AccAddress) MsgCreateArtist {
	return MsgCreateArtist{meta, owner}
}

//nolint
func (msg MsgCreateArtist) Route() string { return RouterKey }
func (msg MsgCreateArtist) Type() string  { return TypeMsgCreateArtist }

// ValidateBasic
func (msg MsgCreateArtist) ValidateBasic() sdk.Error {
	if msg.Meta == nil {
		return ErrInvalidArtistMeta(DefaultCodespace, "missing meta")
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return msg.Meta.ValidateBasic()
}

// Implements Msg.
func (msg MsgCreateArtist) String() string {
	return fmt.Sprintf(`Create Artist Message:
  Meta:         %s
  Owner: %s
`, msg.Meta.String(), msg.Owner.String())
}

// Implements Msg.
func (msg MsgCreateArtist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateArtist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
