package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

type msgServer struct {
	*Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the token MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) Issue(goCtx context.Context, msg *types.MsgIssue) (*types.MsgIssueResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.Authority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.Authority)
	}

	denom, err := m.Keeper.Issue(ctx, msg.Name, msg.Symbol, msg.URI, msg.MaxSupply, owner)
	if err != nil {
		return nil, err
	}
	m.Logger(ctx).Info(fmt.Sprintf("minted a new fantoken denom: %s", denom))
	ctx.EventManager().EmitTypedEvent(&types.EventIssue{
		Denom: denom,
	})

	return &types.MsgIssueResponse{}, nil
}

func (m msgServer) DisableMint(goCtx context.Context, msg *types.MsgDisableMint) (*types.MsgDisableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.DisableMint(ctx, msg.Denom, authority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventDisableMint{
		Denom: msg.Denom,
	})

	return &types.MsgDisableMintResponse{}, nil
}

func (m msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	var recipient sdk.AccAddress

	if msg.Recipient != "" {
		recipient, err = sdk.AccAddressFromBech32(msg.Recipient)
		if err != nil {
			return nil, err
		}
	} else {
		recipient = authority
	}

	if m.Keeper.blockedAddrs[recipient.String()] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", recipient)
	}

	if err := m.Keeper.Mint(ctx, recipient, msg.Denom, msg.Amount, authority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventMint{
		Denom:     msg.Denom,
		Amount:    msg.Amount.String(),
		Recipient: recipient.String(),
	})

	return &types.MsgMintResponse{}, nil
}

func (m msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.Sender] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.Sender)
	}

	if err := m.Keeper.Burn(ctx, msg.Denom, msg.Amount, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventBurn{
		Denom:  msg.Denom,
		Amount: msg.Amount.String(),
	})

	return &types.MsgBurnResponse{}, nil
}

func (m msgServer) TransferAuthority(goCtx context.Context, msg *types.MsgTransferAuthority) (*types.MsgTransferAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	srcAuthority, err := sdk.AccAddressFromBech32(msg.SrcAuthority)
	if err != nil {
		return nil, err
	}

	dstAuthority, err := sdk.AccAddressFromBech32(msg.DstAuthority)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.SrcAuthority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.SrcAuthority)
	}

	if m.Keeper.blockedAddrs[msg.DstAuthority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.DstAuthority)
	}

	if err := m.Keeper.TransferAuthority(ctx, msg.Denom, srcAuthority, dstAuthority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventTransferAuthority{
		Denom:        msg.Denom,
		SrcAuthority: msg.SrcAuthority,
		DstAuthority: msg.DstAuthority,
	})

	return &types.MsgTransferAuthorityResponse{}, nil
}
