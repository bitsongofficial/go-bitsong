package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/nft/keeper"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

// NewHandler handles all messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateNFT:
			res, err := msgServer.CreateNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgTransferNFT:
			res, err := msgServer.TransferNFT(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSignMetadata:
			res, err := msgServer.SignMetadata(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateMetadata:
			res, err := msgServer.UpdateMetadata(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateMetadataAuthority:
			res, err := msgServer.UpdateMetadataAuthority(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateCollection:
			res, err := msgServer.CreateCollection(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateCollectionAuthority:
			res, err := msgServer.UpdateCollectionAuthority(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized nft message type: %T", msg)
		}
	}
}
