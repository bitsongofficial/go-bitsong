package types

import (
	"fmt"
	"strings"

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
	if len(strings.TrimSpace(msg.Meta.Name)) == 0 {
		return ErrInvalidArtistMeta(DefaultCodespace, "artist name cannot be blank")
	}

	if len(msg.Meta.Name) > MaxNameLength {
		return ErrInvalidArtistMeta(DefaultCodespace, fmt.Sprintf("artist name is longer than max length of %d", MaxNameLength))
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	// TODO: to remove
	// return msg.ValidateBasic()
	return nil
}

// Implements Msg.
func (msg MsgCreateArtist) String() string {
	return fmt.Sprintf(`Create Artist Message:
  Name:         %s
  Owner: %s
`, msg.Meta.Name, msg.Owner.String())
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
