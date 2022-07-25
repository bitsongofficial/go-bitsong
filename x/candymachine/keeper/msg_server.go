package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
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

func (m msgServer) CreateCandyMachine(goCtx context.Context, msg *types.MsgCreateCandyMachine) (*types.MsgCreateCandyMachineResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.CreateCandyMachine(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateCandyMachineResponse{}, nil
}

func (m msgServer) UpdateCandyMachine(goCtx context.Context, msg *types.MsgUpdateCandyMachine) (*types.MsgUpdateCandyMachineResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.UpdateCandyMachine(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateCandyMachineResponse{}, nil
}

func (m msgServer) CloseCandyMachine(goCtx context.Context, msg *types.MsgCloseCandyMachine) (*types.MsgCloseCandyMachineResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.CloseCandyMachine(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCloseCandyMachineResponse{}, nil
}

func (m msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.MintNFT(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgMintNFTResponse{}, nil
}
