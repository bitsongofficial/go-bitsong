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
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgTrackCreate(ctx sdk.Context, keeper Keeper, msg types.MsgTrackCreate) (*sdk.Result, error) {
	track := types.NewTrack(msg.TrackID, msg.TrackInfo, msg.Creator)

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
			types.EventTypeTrackCreated,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%s", trackID)),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(trackID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgTrackAddShare(ctx sdk.Context, keeper Keeper, msg types.MsgTrackAddShare) (*sdk.Result, error) {
	if err := keeper.AddShares(ctx, msg.TrackID, msg.Shares, msg.Entity); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTrackAddedShare,
			sdk.NewAttribute(types.AttributeKeyTrackID, fmt.Sprintf("%s", msg.TrackID)),
			sdk.NewAttribute("entity", fmt.Sprintf("%s", msg.Entity.String())),
			sdk.NewAttribute("share", fmt.Sprintf("%s", msg.Shares.String())),
		),
	)

	return &sdk.Result{
		//Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(msg.TrackID),
		Events: ctx.EventManager().Events(),
	}, nil
}
