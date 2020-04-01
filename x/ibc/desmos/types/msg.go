package types

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
)

type MsgCreateSongPost struct {
	SourcePort    string `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	SourceChannel string `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	DestHeight    uint64 `json:"dest_height" yaml:"dest_height"`       // the current height of the destination chain

	SongID       string         `json:"song_id" yaml:"song_id"`             // Song ID
	CreationTime time.Time      `json:"creation_time" yaml:"creation_time"` // Post creation time
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`               // Post creator
}

// NewMsgCreateSongPost creates a new MsgCreateSongPost instance
func NewMsgCreateSongPost(
	sourcePort, sourceChannel string, destHeight uint64,
	songID string, creationTime time.Time, sender sdk.AccAddress,
) MsgCreateSongPost {
	return MsgCreateSongPost{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		DestHeight:    destHeight,

		SongID:       songID,
		CreationTime: creationTime,
		Sender:       sender,
	}
}

// Route implements sdk.Msg
func (MsgCreateSongPost) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgCreateSongPost) Type() string {
	return "create-post"
}

// ValidateBasic implements sdk.Msg
func (msg MsgCreateSongPost) ValidateBasic() error {
	if err := host.DefaultPortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.DefaultChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if len(strings.TrimSpace(msg.SongID)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "song id cannot be empty")
	}
	if msg.CreationTime.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creation time cannot be empty")
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgCreateSongPost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgCreateSongPost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
