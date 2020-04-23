package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/track/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Handle all "track" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgCreate:
			return handleMsgCreate(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized track message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreate handles the creation of a new track
func handleMsgCreate(ctx sdk.Context, keeper Keeper, msg types.MsgCreate) (*sdk.Result, error) {
	track := types.NewTrack(msg.Title, msg.Media, msg.Attributes, msg.Rewards, msg.RightsHolders, ctx.BlockHeader().Time, msg.Owner)
	trackAddr := keeper.Create(ctx, track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackCreated,
			sdk.NewAttribute(types.AttributeKeyTrackAddr, trackAddr.String()),
		),
	)

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(trackAddr),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
