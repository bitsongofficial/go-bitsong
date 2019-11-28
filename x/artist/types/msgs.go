package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/****
 * Artist Msg
 ***/

// Artist message types and routes
const (
	TypeMsgCreateArtist = "create_artist"
)

var _ sdk.Msg = MsgCreateArtist{}

// MsgCreateArtist defines CreateArtist message
type MsgCreateArtist struct {
	Name  string         `json:"name"`  // Artist name
	Owner sdk.AccAddress `json:"owner"` // Artist owner
}

func NewMsgCreateArtist(name string, owner sdk.AccAddress) MsgCreateArtist {
	return MsgCreateArtist{
		Name:  name,
		Owner: owner,
	}
}

//nolint
func (msg MsgCreateArtist) Route() string { return RouterKey }
func (msg MsgCreateArtist) Type() string  { return TypeMsgCreateArtist }

// ValidateBasic
func (msg MsgCreateArtist) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(msg.Name)) == 0 {
		return ErrInvalidArtistName(DefaultCodespace, "artist name cannot be blank")
	}

	if len(msg.Name) > MaxNameLength {
		return ErrInvalidArtistName(DefaultCodespace, fmt.Sprintf("artist name is longer than max length of %d", MaxNameLength))
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgCreateArtist) String() string {
	return fmt.Sprintf(`Create Artist Message:
  Name:         %s
  Owner: %s
`, msg.Name, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateArtist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateArtist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
