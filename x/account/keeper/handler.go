package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgRegisterHandle:
			return handleMsgRegisterHandle(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgRegisterHandle(ctx sdk.Context, keeper Keeper, msg types.MsgRegisterHandle) (*sdk.Result, error) {
	bacc, err := keeper.RegisterHandle(ctx, msg.From, msg.Handle)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterHandle,
			sdk.NewAttribute(types.AttributeKeyHandle, bacc.Handle),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
