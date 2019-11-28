package artist

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "artist" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateArtist:
			return handleMsgCreateArtist(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized artist message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateArtist(ctx sdk.Context, keeper Keeper, msg types.MsgCreateArtist) sdk.Result {
	artist, err := keeper.CreateArtist(ctx, msg.Meta, msg.Owner)
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
