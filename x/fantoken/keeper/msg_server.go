package keeper

import (
	"context"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	if err := m.Keeper.Mint(ctx, minter, recipient, msg.Coin); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventMint{
		Recipient: recipient.String(),
		Coin:      msg.Coin.String(),
	})

	return &types.MsgMintResponse{}, nil
}

func (m msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.Burn(ctx, msg.Coin, owner); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventBurn{
		Sender: msg.Sender,
		Coin:   msg.Coin.String(),
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

func (m msgServer) DisableMint(goCtx context.Context, msg *types.MsgDisableMint) (*types.MsgDisableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.SetMinter(ctx, msg.Denom, minter, sdk.AccAddress{}); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventDisableMint{
		Denom: msg.Denom,
	})

	return &types.MsgDisableMintResponse{}, nil
}

func (m msgServer) SetUri(goCtx context.Context, msg *types.MsgSetUri) (*types.MsgSetUriResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	if err := m.Keeper.SetUri(ctx, msg.Denom, msg.URI, authority); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSetUri{
		Denom: msg.Denom,
	})

	return &types.MsgSetUriResponse{}, nil
}
