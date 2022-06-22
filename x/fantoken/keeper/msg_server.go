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

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.Authority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.Authority)
	}

	// at the moment is disabled, will be enabled once some test will be done
	if m.Keeper.blockedAddrs[msg.Minter] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.Minter)
	}

	denom, err := m.Keeper.Issue(ctx, msg.Name, msg.Symbol, msg.URI, msg.MaxSupply, minter, authority)
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

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.DisableMint(ctx, msg.Denom, minter); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventDisableMint{
		Denom: msg.Denom,
	})

	return &types.MsgDisableMintResponse{}, nil
}

func (m msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
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
		recipient = minter
	}

	if m.Keeper.blockedAddrs[recipient.String()] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", recipient)
	}

	if err := m.Keeper.Mint(ctx, recipient, msg.Denom, msg.Amount, minter); err != nil {
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

func (m msgServer) SetAuthority(goCtx context.Context, msg *types.MsgSetAuthority) (*types.MsgSetAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	oldAuthority, err := sdk.AccAddressFromBech32(msg.OldAuthority)
	if err != nil {
		return nil, err
	}

	newAuthority, err := sdk.AccAddressFromBech32(msg.NewAuthority)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.OldAuthority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.OldAuthority)
	}

	if m.Keeper.blockedAddrs[msg.NewAuthority] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.NewAuthority)
	}

	if err := m.Keeper.SetAuthority(ctx, msg.Denom, oldAuthority, newAuthority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSetAuthority{
		Denom:        msg.Denom,
		OldAuthority: msg.OldAuthority,
		NewAuthority: msg.NewAuthority,
	})

	return &types.MsgSetAuthorityResponse{}, nil
}

func (m msgServer) SetMinter(goCtx context.Context, msg *types.MsgSetMinter) (*types.MsgSetMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	oldMinter, err := sdk.AccAddressFromBech32(msg.OldMinter)
	if err != nil {
		return nil, err
	}

	newMinter, err := sdk.AccAddressFromBech32(msg.NewMinter)
	if err != nil {
		return nil, err
	}

	if m.Keeper.blockedAddrs[msg.OldMinter] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.OldMinter)
	}

	if m.Keeper.blockedAddrs[msg.NewMinter] {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is a module account", msg.NewMinter)
	}

	if err := m.Keeper.SetMinter(ctx, msg.Denom, oldMinter, newMinter); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSetMinter{
		Denom:     msg.Denom,
		OldMinter: msg.OldMinter,
		NewMinter: msg.NewMinter,
	})

	return &types.MsgSetMinterResponse{}, nil
}
