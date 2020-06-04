package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/content/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Handle all "content" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgContentAdd:
			return handleMsgContentAdd(ctx, keeper, msg)
		case types.MsgContentAction:
			return handleMsgAction(ctx, keeper, msg)
		case types.MsgStoreHls:
			return handleMsgStoreHls(ctx, keeper, &msg)

		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgStoreHls(ctx sdk.Context, k Keeper, msg *types.MsgStoreHls) (*sdk.Result, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	hlsID, err := k.CreateHls(ctx, msg.From, msg.HLSByteCode)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentAdded,
			sdk.NewAttribute(types.AttributeKeyContentUri, fmt.Sprintf("%d", hlsID)), // TODO: fix events
		),
	)

	return &sdk.Result{
		Data:   []byte(fmt.Sprintf("%d", hlsID)),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgContentAdd(ctx sdk.Context, keeper Keeper, msg types.MsgContentAdd) (*sdk.Result, error) {
	content := *types.NewContent(
		msg.Uri,
		msg.Hash,
		msg.Dao,
	)

	uri, err := keeper.Add(ctx, &content)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentAdded,
			sdk.NewAttribute(types.AttributeKeyContentUri, uri),
		),
	)

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(uri),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgAction(ctx sdk.Context, keeper Keeper, msg types.MsgContentAction) (*sdk.Result, error) {
	err := keeper.Action(ctx, msg.Uri, msg.From)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentAction,
			sdk.NewAttribute(types.AttributeKeyContentUri, msg.Uri),
			sdk.NewAttribute(types.AttributeKeyAction, ""),
		),
	)

	return &sdk.Result{
		Data:   nil,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
