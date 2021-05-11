package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/bitsong/x/fantoken/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) IssueFanToken(goCtx context.Context, msg *types.MsgIssueFanToken) (*types.MsgIssueFanTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.Owner] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.Owner)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	fmt.Println("IssueFanToken1")
	// handle fee for token
	if err := m.Keeper.DeductIssueFanTokenFee(ctx, owner, msg.Denom); err != nil {
		fmt.Println("IssueFanToken Error", err)
		return nil, err
	}

	if err := m.Keeper.IssueFanToken(
		ctx, msg.Denom, msg.Name, msg.MaxSupply, msg.Mintable, msg.MetadataUri, owner,
	); err != nil {
		fmt.Println("IssueFanToken2 Error", err)
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeIssueFanToken,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgIssueFanTokenResponse{}, nil
}

func (m msgServer) UpdateFanTokenMintable(goCtx context.Context, msg *types.MsgUpdateFanTokenMintable) (*types.MsgUpdateFanTokenMintableResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.UpdateFanTokenMintable(
		ctx, msg.Denom, msg.Mintable, owner,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateFanTokenMintable,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgUpdateFanTokenMintableResponse{}, nil
}

func (m msgServer) MintFanToken(goCtx context.Context, msg *types.MsgMintFanToken) (*types.MsgMintFanTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	var recipient sdk.AccAddress

	if len(msg.Recipient) != 0 {
		recipient, err = sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			return nil, err
		}
	} else {
		recipient = owner
	}

	if m.Keeper.blockedAddrs[recipient.String()] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", recipient)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.MintFanToken(ctx, recipient, msg.Denom, msg.Amount, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintFanToken,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, recipient.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner),
		),
	})

	return &types.MsgMintFanTokenResponse{}, nil
}

func (m msgServer) BurnFanToken(goCtx context.Context, msg *types.MsgBurnFanToken) (*types.MsgBurnFanTokenResponse, error) {
	owner, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnFanToken(ctx, msg.Denom, msg.Amount, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnFanToken,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnFanTokenResponse{}, nil
}

func (m msgServer) TransferFanTokenOwner(goCtx context.Context, msg *types.MsgTransferFanTokenOwner) (*types.MsgTransferFanTokenOwnerResponse, error) {
	srcOwner, err := sdk.AccAddressFromBech32(msg.SrcOwner)
	if err != nil {
		return nil, err
	}

	dstOwner, err := sdk.AccAddressFromBech32(msg.DstOwner)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.DstOwner] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.DstOwner)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.TransferFanTokenOwner(ctx, msg.Denom, srcOwner, dstOwner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferFanTokenOwner,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.SrcOwner),
			sdk.NewAttribute(types.AttributeKeyDstOwner, msg.DstOwner),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SrcOwner),
		),
	})

	return &types.MsgTransferFanTokenOwnerResponse{}, nil
}
