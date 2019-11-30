package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/artist/types"
)

// Handle all "artist" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateArtist:
			return handleMsgCreateArtist(ctx, keeper, msg)
		case types.MsgSetArtistImage:
			return handleMsgSetArtistImage(ctx, keeper, msg)
		case types.MsgSetArtistStatus:
			return handleMsgSetArtistStatus(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized artist message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgCreateArtist handles the creation of a new artist
func handleMsgCreateArtist(ctx sdk.Context, keeper Keeper, msg types.MsgCreateArtist) sdk.Result {
	artist, err := keeper.CreateArtist(ctx, msg.Name, msg.Owner)
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
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(artist.ArtistID),
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgSetArtistImage handles the image of an artist
func handleMsgSetArtistImage(ctx sdk.Context, keeper Keeper, msg types.MsgSetArtistImage) sdk.Result {
	image := types.NewArtistImage(msg.Height, msg.Width, msg.CID)

	err := keeper.SetArtistImage(ctx, msg.ArtistID, image, msg.Owner)
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

// handleMsgSetArtistStatus handles the status of an artist
func handleMsgSetArtistStatus(ctx sdk.Context, keeper Keeper, msg types.MsgSetArtistStatus) sdk.Result {
	err := keeper.SetArtistStatus(ctx, msg.ArtistID, msg.Status, msg.Owner)
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
