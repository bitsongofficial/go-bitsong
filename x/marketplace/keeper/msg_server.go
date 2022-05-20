package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	auctionId, err := m.Keeper.CreateAuction(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateAuctionResponse{Id: auctionId}, nil
}

func (m msgServer) StartAuction(goCtx context.Context, msg *types.MsgStartAuction) (*types.MsgStartAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.StartAuction(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgStartAuctionResponse{}, nil
}

func (m msgServer) SetAuctionAuthority(goCtx context.Context, msg *types.MsgSetAuctionAuthority) (*types.MsgSetAuctionAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.SetAuctionAuthority(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSetAuctionAuthorityResponse{}, nil
}

func (m msgServer) EndAuction(goCtx context.Context, msg *types.MsgEndAuction) (*types.MsgEndAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.EndAuction(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgEndAuctionResponse{}, nil
}

func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.PlaceBid(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgPlaceBidResponse{}, nil
}

func (m msgServer) CancelBid(goCtx context.Context, msg *types.MsgCancelBid) (*types.MsgCancelBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.CancelBid(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCancelBidResponse{}, nil
}

func (m msgServer) ClaimBid(goCtx context.Context, msg *types.MsgClaimBid) (*types.MsgClaimBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.ClaimBid(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgClaimBidResponse{}, nil
}
