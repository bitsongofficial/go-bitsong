package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/launchpad/types"
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

func (m msgServer) CreateLaunchPad(goCtx context.Context, msg *types.MsgCreateLaunchPad) (*types.MsgCreateLaunchPadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.CreateLaunchPad(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateLaunchPadResponse{}, nil
}

func (m msgServer) UpdateLaunchPad(goCtx context.Context, msg *types.MsgUpdateLaunchPad) (*types.MsgUpdateLaunchPadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.UpdateLaunchPad(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateLaunchPadResponse{}, nil
}

func (m msgServer) CloseLaunchPad(goCtx context.Context, msg *types.MsgCloseLaunchPad) (*types.MsgCloseLaunchPadResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.CloseLaunchPad(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCloseLaunchPadResponse{}, nil
}

func (m msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nftId, err := m.Keeper.MintNFT(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgMintNFTResponse{
		NftId: nftId,
	}, nil
}

func (m msgServer) MintNFTs(goCtx context.Context, msg *types.MsgMintNFTs) (*types.MsgMintNFTsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nftIds, err := m.Keeper.MintNFTs(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgMintNFTsResponse{
		NftIds: nftIds,
	}, nil
}
