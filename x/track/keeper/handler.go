package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Handle all "content" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgTrackAdd:
			return handleMsgTrackAdd(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgTrackAdd(ctx sdk.Context, keeper Keeper, msg types.MsgTrackAdd) (*sdk.Result, error) {
	track, err := types.NewTrack(
		msg.Title, msg.Artists, msg.Feat, msg.Producers, msg.Number, msg.Duration, msg.Explicit, msg.ExternalIds, msg.ExternalUrls, msg.PreviewUrl, msg.Dao,
	)

	if err != nil {
		return nil, err
	}

	trackID, err := keeper.Add(ctx, track)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackAdded,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", trackID)),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(cid),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
