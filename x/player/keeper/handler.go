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
	count := keeper.GetPlayersCount(ctx)
	player := types.Player{
		ID:      types.NewPlayerID(count),
		Moniker: msg.Moniker,
		Deposit: msg.Deposit,
		Owner:   msg.Owner,
	}

	keeper.AddDeposit(ctx, msg.Owner, msg.Deposit)
	keeper.SetPlayer(ctx, player)
	keeper.SetPlayersCount(ctx, count+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePlayerRegistered,
			sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", count)),
		),
	)

	return &sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(count),
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
