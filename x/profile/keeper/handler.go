package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/profile/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgProfileCreate:
			return handleMsgProfileCreate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgProfileCreate(ctx sdk.Context, keeper Keeper, msg types.MsgProfileCreate) (*sdk.Result, error) {
	profile, err := keeper.CreateProfile(ctx, msg.Address, msg.Handle, msg.MetadataURI)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProfileCreate,
			sdk.NewAttribute(types.AttributeKeyProfileHandle, profile.Handle),
		),
	)

	return &sdk.Result{
		Data:   keeper.codec.MustMarshalBinaryLengthPrefixed(profile),
		Events: ctx.EventManager().Events(),
	}, nil
}
