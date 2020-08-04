package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
)

// Content messages types and routes
const (
	TypeMsgTrackCreate = "track_create"
)

var _ sdk.Msg = MsgTrackCreate{}

type MsgTrackCreate struct {
	TrackID   string
	TrackInfo []byte         `json:"track_info" yaml:"track_info"`
	Entities  []Entity       `json:"entities" yaml:"entities"`
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgTrackCreate(info []byte, creator sdk.AccAddress, entities []Entity) MsgTrackCreate {
	return MsgTrackCreate{
		TrackID:   uuid.New().String(),
		TrackInfo: info,
		Creator:   creator,
		Entities:  entities,
	}
}

func (msg MsgTrackCreate) Route() string { return RouterKey }
func (msg MsgTrackCreate) Type() string  { return TypeMsgTrackCreate }

func (msg MsgTrackCreate) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track creator address cannot be empty")
	}

	if len(msg.Entities) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track entities cannot be empty")
	}

	if len(msg.TrackInfo) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track info cannot be empty")
	}

	if len(msg.TrackInfo) > MaxTrackInfoLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track info too large")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTrackCreate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTrackCreate) GetSigners() []sdk.AccAddress {
	// TODO: all party must sign
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgTrackCreate) String() string {
	return fmt.Sprintf(`Msg Track Create
Track ID: %s,
TrackInfo: %s,
Creator: %s`,
		msg.TrackID, msg.TrackInfo, msg.Creator,
	)
}
