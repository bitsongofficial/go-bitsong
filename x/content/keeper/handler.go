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
		case types.MsgMintContent:
			return handleMsgMintContent(ctx, keeper, msg)
		case types.MsgBurnContent:
			return handleMsgBurnContent(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgAddContent(ctx sdk.Context, keeper Keeper, msg types.MsgAddContent) (*sdk.Result, error) {
	content := types.NewContent(
		msg.Name,
		msg.Uri,
		msg.MetaUri,
		msg.ContentUri,
		msg.Denom,
		msg.Creator,
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

func handleMsgMintContent(ctx sdk.Context, keeper Keeper, msg types.MsgMintContent) (*sdk.Result, error) {
	err := keeper.Mint(ctx, msg.Uri, msg.Amount, msg.From)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentMinted,
			sdk.NewAttribute(types.AttributeKeyContentUri, msg.Uri),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &sdk.Result{
		Data:   nil,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgBurnContent(ctx sdk.Context, keeper Keeper, msg types.MsgBurnContent) (*sdk.Result, error) {
	err := keeper.Burn(ctx, msg.Uri, msg.Amount, msg.From)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContentBurned,
			sdk.NewAttribute(types.AttributeKeyContentUri, msg.Uri),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &sdk.Result{
		Data:   nil,
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
