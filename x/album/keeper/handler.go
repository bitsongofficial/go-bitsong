package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/album/types"
)

// Handle all "album" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateAlbum:
			return handleMsgCreateAlbum(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized album message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateAlbum handles the creation of a new album
func handleMsgCreateAlbum(ctx sdk.Context, keeper Keeper, msg types.MsgCreateAlbum) sdk.Result {
	album, err := keeper.CreateAlbum(ctx, msg.Title, msg.AlbumType, msg.ReleaseDate, msg.ReleaseDatePrecision, msg.Owner)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	)

	return sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(album.AlbumID),
		Events: ctx.EventManager().Events(),
	}
}
