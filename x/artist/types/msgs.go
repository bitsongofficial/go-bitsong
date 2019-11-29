package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	c "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
)

/****
 * Artist Msg
 ***/

// Artist messages types and routes
const (
	TypeMsgCreateArtist    = "create_artist"
	TypeMsgSetArtistImage  = "set_artist_image"
	TypeMsgSetArtistStatus = "set_artist_status"
)

/****************************************
 * MsgCreateArtist
 ****************************************/

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

/****************************************
 * MsgSetArtistImage
 ****************************************/

var _ sdk.Msg = MsgSetArtistImage{}

// MsgCreateArtist defines CreateArtist message
type MsgSetArtistImage struct {
	ArtistID uint64         `json:"artist_id"` // Artist ID
	Height   uint64         `json:"height"`    // Image height
	Width    uint64         `json:"width"`     // Image width
	CID      string         `json:"cid"`       // Image cid
	Owner    sdk.AccAddress `json:"owner"`     // Artist Owner
}

func NewMsgSetArtistImage(artistID uint64, height uint64, width uint64, cid string, owner sdk.AccAddress) MsgSetArtistImage {
	return MsgSetArtistImage{
		ArtistID: artistID,
		Height:   height,
		Width:    width,
		CID:      cid,
		Owner:    owner,
	}
}

//nolint
func (msg MsgSetArtistImage) Route() string { return RouterKey }
func (msg MsgSetArtistImage) Type() string  { return TypeMsgSetArtistImage }

// ValidateBasic
func (msg MsgSetArtistImage) ValidateBasic() sdk.Error {
	if msg.ArtistID == 0 {
		return ErrUnknownArtist(DefaultCodespace, "unknown artist")
	}

	if msg.Height == 0 {
		return ErrInvalidArtistImageHeight(DefaultCodespace, "image height cannot be blank")
	}

	if msg.Width == 0 {
		return ErrInvalidArtistImageWidth(DefaultCodespace, "image width cannot be blank")
	}

	if len(strings.TrimSpace(msg.CID)) == 0 {
		return ErrInvalidArtistImageCid(DefaultCodespace, "image cid cannot be blank")
	}

	_, err := c.Decode(msg.CID)
	if err != nil {
		return ErrInvalidArtistImageCid(DefaultCodespace, fmt.Sprintf("invalid artist image cid: %s", msg.CID))
	}

	return nil
}

// Implements Msg.
func (msg MsgSetArtistImage) String() string {
	return fmt.Sprintf(`Set Artist Image Message:
  ArtistID:         %d
  Height: %d
  Width: %d
  Cid: %s
  Owner: %s
`, msg.ArtistID, msg.Height, msg.Width, msg.CID, msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgSetArtistImage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetArtistImage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

/****************************************
 * MsgSetArtistStatus
 ****************************************/

var _ sdk.Msg = MsgSetArtistStatus{}

// MsgSetArtistStatus defines SetStatus message
type MsgSetArtistStatus struct {
	ArtistID uint64         `json:"artist_id"` // Artist ID
	Status   ArtistStatus   `json:"status"`    // Status of the Artist {Nil, Verified, Rejected, Failed}
	Owner    sdk.AccAddress `json:"owner"`     // Artist Owner
}

func NewMsgSetArtistStatus(artistID uint64, status ArtistStatus, owner sdk.AccAddress) MsgSetArtistStatus {
	return MsgSetArtistStatus{
		ArtistID: artistID,
		Status:   status,
		Owner:    owner,
	}
}

//nolint
func (msg MsgSetArtistStatus) Route() string { return RouterKey }
func (msg MsgSetArtistStatus) Type() string  { return TypeMsgSetArtistStatus }

// ValidateBasic
func (msg MsgSetArtistStatus) ValidateBasic() sdk.Error {
	if msg.ArtistID == 0 {
		return ErrUnknownArtist(DefaultCodespace, "unknown artist")
	}

	if !msg.Status.Valid() {
		return ErrInvalidArtistStatus(DefaultCodespace, "invalid artist status")
	}

	return nil
}

// Implements Msg.
func (msg MsgSetArtistStatus) String() string {
	return fmt.Sprintf(`Set Artist Status Message:
  ArtistID:         %d
  Status: %s
  Owner: %s
`, msg.ArtistID, msg.Status.String(), msg.Owner.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgSetArtistStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetArtistStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
