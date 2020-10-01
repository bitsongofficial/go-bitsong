package keeper

import (
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/release/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgReleaseCreate:
			return handleMsgReleaseCreate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgReleaseCreate(ctx sdk.Context, keeper Keeper, msg types.MsgReleaseCreate) (*sdk.Result, error) {
	release, err := keeper.CreateRelease(ctx, msg.Creator, msg.ReleaseID, msg.MetadataURI)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeReleaseCreate,
			sdk.NewAttribute(types.AttributeKeyReleaseID, release.ReleaseID),
		),
	)

	return &sdk.Result{
		Data:   keeper.codec.MustMarshalBinaryLengthPrefixed(release),
		Events: ctx.EventManager().Events(),
	}, nil
}
