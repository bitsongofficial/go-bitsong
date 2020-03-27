package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/track/types"
)

// Handle all "track" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateTrack:
			return handleMsgCreateTrack(ctx, keeper, msg)
		case types.MsgPlay:
			return handleMsgPlay(ctx, keeper, msg)
		case types.MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized track message type: %T", msg)
			return nil, sdkerr.Wrap(sdkerr.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateTrack handles the creation of a new track
func handleMsgCreateTrack(ctx sdk.Context, keeper Keeper, msg types.MsgCreateTrack) (*sdk.Result, error) {
	track, err := keeper.CreateTrack(
		ctx,
		msg.Title,
		msg.Audio,
		msg.Image,
		msg.Duration,
		msg.Hidden,
		msg.Explicit,
		msg.Genre,
		msg.Mood,
		msg.Artists,
		msg.Featuring,
		msg.Producers,
		msg.Description,
		msg.Copyright,
		msg.Owner,
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	)

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(track.TrackID),
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgPlay
func handleMsgPlay(ctx sdk.Context, keeper Keeper, msg types.MsgPlay) (*sdk.Result, error) {
	err := keeper.Play(ctx, msg.TrackID, msg.AccAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AccAddr.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDeposit) (*sdk.Result, error) {
	err, verified := keeper.AddDeposit(ctx, msg.TrackID, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, err
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
				types.EventTypeDepositTrack,
				sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", msg.TrackID)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
