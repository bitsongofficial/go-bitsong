package track

import (
	"fmt"

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
	/*song, err := k.Publish(ctx, msg.Title, msg.Owner, msg.Content, msg.RedistributionSplitRate)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, song.Owner.String()),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}*/
	return sdk.Result{}

}

// Handle a message to play track
func hanleMsgPlay(ctx sdk.Context, k Keeper, msg MsgPlay) sdk.Result {
	/*err := k.Play(ctx, msg.SongID, msg.Listener)

	if err != nil {
		return err.Result()
	}*/
	return sdk.Result{}
}
