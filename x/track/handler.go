package track

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "track" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPublish:
			return hanleMsgPublish(ctx, keeper, msg)
		case MsgPlay:
			return hanleMsgPlay(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized track Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to publish track
func hanleMsgPublish(ctx sdk.Context, k Keeper, msg MsgPublish) sdk.Result {
	track, err := k.PublishTrack(ctx, msg.Title, msg.Owner, msg.Content, msg.RedistributionSplitRate)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, track.Owner.String()),
			sdk.NewAttribute(AttributeKeyTrackId, strconv.FormatUint(track.TrackID, 10)),
			sdk.NewAttribute(AttributeKeyTitle, track.Title),
			sdk.NewAttribute(AttributeKeyContent, track.Content),
			sdk.NewAttribute(AttributeKeyRedistributionSplitRate, track.RedistributionSplitRate.String()),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle a message to play track
func hanleMsgPlay(ctx sdk.Context, k Keeper, msg MsgPlay) sdk.Result {
	_, ok := k.SavePlay(ctx, msg.Listener, msg.TrackID)
	if !ok {
		// TODO: fix with error
		return sdk.Result{}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Listener.String()),
			sdk.NewAttribute(AttributeKeyTrackId, strconv.FormatUint(msg.TrackID, 10)),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
