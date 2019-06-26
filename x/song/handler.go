package song

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "song" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPublish:
			return hanleMsgPublish(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized song Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to publish song
func hanleMsgPublish(ctx sdk.Context, k Keeper, msg MsgPublish) sdk.Result {
	song, err := k.Publish(ctx, msg.Title, msg.Owner)
	if err != nil {
		return err.Result()
	}
	resTags := sdk.NewTags(
		Category, TxCategory,
		SongId, fmt.Sprintf("%d", song.SongId),
		Sender, song.Owner.String(),
	)
	return sdk.Result{
		Tags: resTags,
	}
}

// Handle a message to set name
/*func handleMsgSetTitle(ctx sdk.Context, keeper Keeper, msg MsgSetTitle) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Title)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	keeper.SetTitle(ctx, msg.Title) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}*/