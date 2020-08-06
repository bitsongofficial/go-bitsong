package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Content messages types and routes
const (
	TypeMsgTrackCreate      = "track_create"
	TypeMsgTrackAddShare    = "track_add_share"
	TypeMsgTrackRemoveShare = "track_remove_share"
)

var _ sdk.Msg = MsgTrackCreate{}

type MsgTrackCreate struct {
	TrackID   string         `json:"track_id" yaml:"track_id"`
	TrackInfo []byte         `json:"track_info" yaml:"track_info"`
	Entities  []Entity       `json:"entities" yaml:"entities"`
	Creator   sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgTrackCreate(trackID string, info []byte, creator sdk.AccAddress, entities []Entity) MsgTrackCreate {
	return MsgTrackCreate{
		//TrackID:   uuid.New().String(),
		TrackID:   trackID,
		TrackInfo: info,
		Creator:   creator,
		Entities:  entities,
	}
}

func (msg MsgTrackCreate) Route() string { return RouterKey }
func (msg MsgTrackCreate) Type() string  { return TypeMsgTrackCreate }

func (msg MsgTrackCreate) ValidateBasic() error {
	for _, entity := range msg.Entities {
		if entity.Address.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "entity address cannot be empty")
		}
	}

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

var _ sdk.Msg = MsgTrackAddShare{}

type MsgTrackAddShare struct {
	TrackID string         `json:"track_id" yaml:"track_id"`
	Entity  sdk.AccAddress `json:"entity" yaml:"entity"`
	Share   sdk.Coin       `json:"share" yaml:"share"`
}

func NewMsgTrackAddShare(trackID string, share sdk.Coin, entity sdk.AccAddress) MsgTrackAddShare {
	return MsgTrackAddShare{
		TrackID: trackID,
		Entity:  entity,
		Share:   share,
	}
}

func (msg MsgTrackAddShare) Route() string { return RouterKey }
func (msg MsgTrackAddShare) Type() string  { return TypeMsgTrackAddShare }

func (msg MsgTrackAddShare) ValidateBasic() error {
	// TODO: add security checks

	if err := sdk.VerifyAddressFormat(msg.Entity); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "entity address cannot be empty")
	}

	if msg.TrackID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track id cannot be empty")
	}

	if !msg.Share.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "share amount must be positive")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTrackAddShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTrackAddShare) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Entity}
}

func (msg MsgTrackAddShare) String() string {
	return fmt.Sprintf(`Msg Track Add Share
Track ID: %s,
Entity: %s,
Share: %s`,
		msg.TrackID, msg.Entity, msg.Share,
	)
}

var _ sdk.Msg = MsgTrackRemoveShare{}

type MsgTrackRemoveShare struct {
	TrackID string         `json:"track_id" yaml:"track_id"`
	Entity  sdk.AccAddress `json:"entity" yaml:"entity"`
	Share   sdk.Coin       `json:"share" yaml:"share"`
}

func NewMsgTrackRemoveShare(trackID string, share sdk.Coin, entity sdk.AccAddress) MsgTrackRemoveShare {
	return MsgTrackRemoveShare{
		TrackID: trackID,
		Entity:  entity,
		Share:   share,
	}
}

func (msg MsgTrackRemoveShare) Route() string { return RouterKey }
func (msg MsgTrackRemoveShare) Type() string  { return TypeMsgTrackAddShare }

func (msg MsgTrackRemoveShare) ValidateBasic() error {
	// TODO: add security checks

	if err := sdk.VerifyAddressFormat(msg.Entity); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "entity address cannot be empty")
	}

	if msg.TrackID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "track id cannot be empty")
	}

	if !msg.Share.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "share amount must be positive")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTrackRemoveShare) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgTrackRemoveShare) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Entity}
}

func (msg MsgTrackRemoveShare) String() string {
	return fmt.Sprintf(`Msg Track Remove Share
Track ID: %s,
Entity: %s,
Share: %s`,
		msg.TrackID, msg.Entity, msg.Share,
	)
}
