package track

import (
	"fmt"

	"github.com/BitSongOfficial/go-bitsong/x/track/types"
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
	song, err := k.Publish(ctx, msg.Title, msg.Owner, msg.Content, msg.RedistributionSplitRate)
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
	}

	// TODO: remove
	/*resTags := sdk.NewTags(
		Category, TxCategory,
		SongID, fmt.Sprintf("%d", track.SongID),
		Owner, track.Owner.String(),
		Content, track.Content,
		TotalReward, fmt.Sprintf("%d", track.TotalReward),
		RedistributionSplitRate, track.RedistributionSplitRate,
	)
	return sdk.Result{
		Tags: resTags,
	}*/
}

// Handle a message to play track
func hanleMsgPlay(ctx sdk.Context, k Keeper, msg MsgPlay) sdk.Result {
	err := k.Play(ctx, msg.SongID, msg.Listener)

	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle a message to set name
/*func handleMsgSetTitle(ctx sdk.Context, keeper Keeper, msg MsgSetTitle) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Title)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	keeper.SetTitle(ctx, msg.Title) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}*/
