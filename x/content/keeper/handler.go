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
		case types.MsgAddContent:
			return handleMsgAddContent(ctx, keeper, msg)
		case types.MsgStream:
			return handleMsgStream(ctx, keeper, msg)
		case types.MsgDownload:
			return handleMsgDownload(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgAddContent(ctx sdk.Context, keeper Keeper, msg types.MsgAddContent) (*sdk.Result, error) {
	streamPrice, err := sdk.ParseCoin(msg.StreamPrice)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	downloadPrice, err := sdk.ParseCoin(msg.DownloadPrice)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	content := types.NewContent(
		msg.Name,
		msg.Uri,
		msg.MetaUri,
		msg.ContentUri,
		msg.Denom,
		streamPrice,
		downloadPrice,
		msg.Creator,
		msg.RightsHolders,
	)

	uri, err := keeper.Add(ctx, content)
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

func handleMsgStream(ctx sdk.Context, keeper Keeper, msg types.MsgStream) (*sdk.Result, error) {
	err := keeper.Stream(ctx, msg.Uri, msg.From)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentStreamed,
			sdk.NewAttribute(types.AttributeKeyContentUri, msg.Uri),
		),
	)

	return &sdk.Result{
		Data:   nil,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgDownload(ctx sdk.Context, keeper Keeper, msg types.MsgDownload) (*sdk.Result, error) {
	err := keeper.Download(ctx, msg.Uri, msg.From)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentDownloaded,
			sdk.NewAttribute(types.AttributeKeyContentUri, msg.Uri),
		),
	)

	return &sdk.Result{
		Data:   nil,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
