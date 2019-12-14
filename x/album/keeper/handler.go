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
		case types.MsgAddTrackAlbum:
			return handleMsgAddTrackAlbum(ctx, keeper, msg)
		case types.MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized album message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateAlbum handles the creation of a new album
func handleMsgCreateAlbum(ctx sdk.Context, keeper Keeper, msg types.MsgCreateAlbum) sdk.Result {
	album, err := keeper.CreateAlbum(ctx, msg.Title, msg.AlbumType, msg.MetadataURI, msg.Owner)
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

// handleMsgAddTrackAlbum
func handleMsgAddTrackAlbum(ctx sdk.Context, keeper Keeper, msg types.MsgAddTrackAlbum) sdk.Result {
	err := keeper.AddTrack(ctx, msg.AlbumID, msg.TrackID, 0)
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
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDeposit) sdk.Result {
	err, verified := keeper.AddDeposit(ctx, msg.AlbumID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	if verified {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDepositAlbum,
				sdk.NewAttribute(types.AttributeKeyAlbumID, fmt.Sprintf("%d", msg.AlbumID)),
			),
		)
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
