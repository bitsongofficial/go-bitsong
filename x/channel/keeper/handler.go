package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/channel/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgChannelCreate:
			return handleMsgChannelCreate(ctx, keeper, msg)
		case types.MsgChannelEdit:
			return handleMsgChannelEdit(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgChannelCreate(ctx sdk.Context, keeper Keeper, msg types.MsgChannelCreate) (*sdk.Result, error) {
	channel, err := keeper.CreateChannel(ctx, msg.Owner, msg.Handle, msg.MetadataURI)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChannelCreate,
			sdk.NewAttribute(types.AttributeKeyProfileHandle, channel.Handle),
		),
	)

	return &sdk.Result{
		Data:   keeper.codec.MustMarshalBinaryLengthPrefixed(channel),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgChannelEdit(ctx sdk.Context, keeper Keeper, msg types.MsgChannelEdit) (*sdk.Result, error) {
	channel, err := keeper.EditChannel(ctx, msg.Owner, msg.MetadataURI)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChannelEdit,
			sdk.NewAttribute(types.AttributeKeyProfileHandle, channel.Handle),
		),
	)

	return &sdk.Result{
		Data:   keeper.codec.MustMarshalBinaryLengthPrefixed(channel),
		Events: ctx.EventManager().Events(),
	}, nil
}
