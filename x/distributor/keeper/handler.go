package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/distributor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateDistributor:
			return handleMsgCreateDistributor(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized track message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateDistributor(ctx sdk.Context, keeper Keeper, msg types.MsgCreateDistributor) sdk.Result {
	distributor, err := keeper.CreateDistributor(ctx, msg.Name, msg.Address)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	)

	return sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(distributor.Address),
		Events: ctx.EventManager().Events(),
	}
}
