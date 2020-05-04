package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/player/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgRegisterPlayer:
			return handleMsgRegisterPlayer(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgRegisterPlayer(ctx sdk.Context, keeper Keeper, msg types.MsgRegisterPlayer) (*sdk.Result, error) {
	if err := keeper.RegisterPlayer(ctx, msg.Moniker, msg.PlayerAddr, msg.Validator); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePlayerRegistered,
			sdk.NewAttribute(types.AttributeKeyValidator, fmt.Sprintf("%s", msg.Validator)),
		),
	)

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(msg.Validator),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
