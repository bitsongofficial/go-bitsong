package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

const (
	TypeMsgReleaseCreate = "release_create"
)

var _ sdk.Msg = MsgReleaseCreate{}

type MsgReleaseCreate struct {
	ReleaseID   string         `json:"release_id"`
	MetadataURI string         `json:"metadata_uri"`
	Creator     sdk.AccAddress `json:"creator"`
}

func NewMsgReleseCreate(releaseID, metadataURI string, creator sdk.AccAddress) MsgReleaseCreate {
	return MsgReleaseCreate{
		ReleaseID:   releaseID,
		MetadataURI: metadataURI,
		Creator:     creator,
	}
}

func (msg MsgReleaseCreate) Route() string { return RouterKey }
func (msg MsgReleaseCreate) Type() string  { return TypeMsgReleaseCreate }

func (msg MsgReleaseCreate) ValidateBasic() error {
	release := NewRelease(msg.ReleaseID, msg.MetadataURI, msg.Creator, time.Time{})
	return release.Validate()
}

// GetSignBytes encodes the message for signing
func (msg MsgReleaseCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgReleaseCreate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgReleaseCreate) String() string {
	return fmt.Sprintf(`Msg Release Create
  ReleaseID: %s,
  MetadataURI: %s,
  Creator: %s`,
		msg.ReleaseID, msg.MetadataURI, msg.Creator.String(),
	)
}
