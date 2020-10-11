package keeper

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	"github.com/bitsongofficial/go-bitsong/x/artist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgArtistCreate:
			return handleMsgArtistCreate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized content message type: %T", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgArtistCreate(ctx sdk.Context, keeper Keeper, msg types.MsgArtistCreate) (*sdk.Result, error) {
	artistData := types.NewArtist(btsg.ID(msg.ID), msg.Name, msg.URLs, msg.Genres, msg.MetadataURI, msg.Creator)
	artist, err := keeper.CreateArtist(ctx, artistData)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeArtistCreate,
			sdk.NewAttribute(types.AttributeKeyID, artist.ID.String()),
		),
	)

	return &sdk.Result{
		Data:   keeper.codec.MustMarshalBinaryLengthPrefixed(artist),
		Events: ctx.EventManager().Events(),
	}, nil
}
