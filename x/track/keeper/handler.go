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
		case types.MsgTrackPublish:
			return handleMsgTrackPublish(ctx, keeper, msg)
		case types.MsgTrackTokenize:
			return handleMsgTrackTokenize(ctx, keeper, msg)
		case types.MsgTokenMint:
			return handleMsgTokenMint(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgTrackPublish(ctx sdk.Context, keeper Keeper, msg types.MsgTrackPublish) (*sdk.Result, error) {
	track := types.NewTrack(msg.TrackInfo, msg.Creator)

	trackID, err := keeper.Add(ctx, track)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackPublished,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", trackID)),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(cid),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgTrackTokenize(ctx sdk.Context, keeper Keeper, msg types.MsgTrackTokenize) (*sdk.Result, error) {
	track, ok := keeper.GetTrack(ctx, msg.TrackID)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trackID not found")
	}

	if !msg.Creator.Equals(track.Creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creator address is invalid")
	}

	if track.TokenInfo.Tokenized {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trackID is already tokenized")
	}

	track.TokenInfo = types.NewTokenInfo(msg.Denom)
	keeper.SetTrack(ctx, &track)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackTokenized,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%d", track.TrackID)),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(cid),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgTokenMint(ctx sdk.Context, keeper Keeper, msg types.MsgTokenMint) (*sdk.Result, error) {
	track, ok := keeper.GetTrack(ctx, msg.TrackID)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trackID not found")
	}

	if !msg.Creator.Equals(track.Creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "creator address is invalid")
	}

	if !track.TokenInfo.Mintable {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trackID is not mintable")
	}

	keeper.Mint(ctx, msg.Amount, msg.Recipient)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenMinted,
			sdk.NewAttribute(types.AttributeKeyTokenAmount, fmt.Sprintf("%s", msg.Amount)),
			sdk.NewAttribute(types.AttributeKeyTokenRecipient, fmt.Sprintf("%s", msg.Recipient.String())),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(cid),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
