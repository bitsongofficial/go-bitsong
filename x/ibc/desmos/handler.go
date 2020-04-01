package desmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns sdk.Handler for IBC token transfer module messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgCreateSongPost:
			return handleMsgCreatePost(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 transfer message type: %T", msg)
		}
	}
}

func handleMsgCreatePost(ctx sdk.Context, k Keeper, msg MsgCreateSongPost) (*sdk.Result, error) {
	if err := k.SendPostCreation(
		ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, msg.SongID, msg.CreationTime, msg.Sender,
	); err != nil {
		return nil, err
	}

	// TODO: Replace with custom event
	//ctx.EventManager().EmitEvent(
	//sdk.NewEvent(
	//	sdk.EventTypeMessage,
	//	sdk.NewAttribute(sdk.AttributeKeyModule, ibcxfer.AttributeValueCategory),
	//	sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
	//),
	//)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
