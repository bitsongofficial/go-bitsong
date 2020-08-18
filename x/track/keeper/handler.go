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
		case types.MsgTrackCreate:
			return handleMsgTrackCreate(ctx, keeper, msg)
		case types.MsgTrackAddShare:
			return handleMsgTrackAddShare(ctx, keeper, msg)
		case types.MsgTrackRemoveShare:
			return handleMsgTrackRemoveShare(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgTrackCreate(ctx sdk.Context, keeper Keeper, msg types.MsgTrackCreate) (*sdk.Result, error) {
	//track := types.NewTrack(msg.TrackID, msg.TrackInfo, msg.Creator, msg.Provider, msg.StreamPrice, msg.DownloadPrice)
	track := types.NewTrack(msg.TrackID, msg.Creator, msg.Provider, msg.StreamPrice, msg.DownloadPrice)

	// store track
	trackID, err := keeper.Add(ctx, track)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	for _, entity := range msg.Entities {
		// mint nft
		// TODO: convert with standard nft module
		coin := sdk.Coin{
			Denom:  track.ToCoinDenom(),
			Amount: entity.Shares, // TODO: entity shares must be > 0
		}
		if err := keeper.MintAndSend(ctx, coin, entity.Address); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackCreate,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%s", trackID)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTrackAddShare(ctx sdk.Context, keeper Keeper, msg types.MsgTrackAddShare) (*sdk.Result, error) {
	if err := keeper.AddShare(ctx, msg.TrackID, msg.Entity, msg.Share); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackAddShare,
			sdk.NewAttribute(types.AttributeKeyTrackID, msg.TrackID),
			sdk.NewAttribute(types.AttributeKeyEntity, msg.Entity.String()),
			sdk.NewAttribute(types.AttributeKeyShare, msg.Share.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgTrackRemoveShare(ctx sdk.Context, keeper Keeper, msg types.MsgTrackRemoveShare) (*sdk.Result, error) {
	if err := keeper.RemoveShare(ctx, msg.TrackID, msg.Entity, msg.Share); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackAddShare,
			sdk.NewAttribute(types.AttributeKeyTrackID, msg.TrackID),
			sdk.NewAttribute(types.AttributeKeyEntity, msg.Entity.String()),
			sdk.NewAttribute(types.AttributeKeyShare, msg.Share.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
