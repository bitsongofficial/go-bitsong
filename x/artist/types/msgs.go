package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/****
 * Artist Msg
 ***/

// Artist messages types and routes
const (
	TypeMsgCreateArtist = "create_artist"
	TypeMsgDeposit      = "deposit"
)

/****************************************
 * MsgCreateArtist
 ****************************************/

var _ sdk.Msg = MsgCreateArtist{}

// MsgCreateArtist defines CreateArtist message
type MsgCreateArtist struct {
	Name        string         `json:"name"` // Artist name
	MetadataURI string         `json:"metadata_uri"`
	Owner       sdk.AccAddress `json:"owner"` // Artist owner
}

func NewMsgCreateArtist(name string, metadataUri string, owner sdk.AccAddress) MsgCreateArtist {
	return MsgCreateArtist{
		Name:        name,
		MetadataURI: metadataUri,
		Owner:       owner,
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

	// TODO:
	// - Add more check for CID (Metadata uri ipfs:)
	if len(strings.TrimSpace(msg.MetadataURI)) == 0 {
		return ErrInvalidArtistMetadataURI(DefaultCodespace, "artist metadata uri cannot be blank")
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
  MetadataURI: %s
  Address: %s
`, msg.Name, msg.MetadataURI, msg.Owner.String())
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
 * MsgDeposit
 ****************************************/

var _ sdk.Msg = MsgDeposit{}

type MsgDeposit struct {
	ArtistID  uint64         `json:"artist_id" yaml:"artist_id"` // ID of the artist
	Depositor sdk.AccAddress `json:"depositor" yaml:"depositor"` // Address of the depositor
	Amount    sdk.Coins      `json:"amount" yaml:"amount"`       // Coins to add to the proposal's deposit
}

func NewMsgDeposit(depositor sdk.AccAddress, artistID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{artistID, depositor, amount}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return RouterKey }
func (msg MsgDeposit) Type() string  { return TypeMsgDeposit }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}

	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf(`Deposit Message:
  Depositer:   %s
  Artist ID: %d
  Amount:      %s
`, msg.Depositor, msg.ArtistID, msg.Amount)
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}
