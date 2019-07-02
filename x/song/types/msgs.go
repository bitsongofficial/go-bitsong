package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgPublish defines a Publish message
type MsgPublish struct {
	Title                   string         `json:"title"`
	Content                 string         `json:"content"`
	Owner                   sdk.AccAddress `json:"owner"`
	TotalReward             sdk.Int        `json:"total_reward"`
	RedistributionSplitRate string         `json:"redistribution_split_rate"`
}

// MsgPublish defines a Publish message
type MsgPlay struct {
	SongID   string         `json:"song_id"`
	Listener sdk.AccAddress `json:"listener"`
}

// NewMsgPublish is a constructor function for MsgPublish
func NewMsgPublish(title string, owner sdk.AccAddress, content string, redistributionSplitRate string) MsgPublish {
	return MsgPublish{
		Title:                   title,
		Content:                 content,
		Owner:                   owner,
		TotalReward:             sdk.NewInt(0),
		RedistributionSplitRate: redistributionSplitRate,
	}
}

// NewMsgPlay is a constructor function for MsgPublish
func NewMsgPlay(songID string, listener sdk.AccAddress) MsgPlay {
	return MsgPlay{
		SongID:   songID,
		Listener: listener,
	}
}

// ValidateBasic runs stateless checks on the message
func (msg MsgPublish) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Title) == 0 {
		return sdk.ErrUnknownRequest("Title cannot be empty")
	}
	if len(msg.Content) == 0 {
		return sdk.ErrUnknownRequest("Content cannot be empty")
	}
	if len(msg.RedistributionSplitRate) == 0 {
		return sdk.ErrUnknownRequest("Redistribution Split Rate cannot be empty")
	}

	return nil
}

func (msg MsgPlay) ValidateBasic() sdk.Error {
	if msg.Listener.Empty() {
		return sdk.ErrInvalidAddress(msg.Listener.String())
	}
	if len(msg.SongID) == 0 {
		return sdk.ErrUnknownRequest("Id cannot be empty")
	}
	return nil
}

// Route should return the name of the module
func (msg MsgPublish) Route() string { return RouterKey }
func (msg MsgPlay) Route() string    { return RouterKey }

// Type should return the action
func (msg MsgPublish) Type() string { return "publish" }
func (msg MsgPlay) Type() string    { return "play" }

// GetSignBytes encodes the message for signing
func (msg MsgPublish) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgPlay) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgPublish) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgPlay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Listener}
}
