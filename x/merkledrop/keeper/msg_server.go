package keeper

import (
	"context"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper}
}

func (m msgServer) CreateMerkledrop(goCtx context.Context, msg *types.MsgCreateMerkledrop) (*types.MsgCreateMerkledropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// create merkledrop
	mdId := m.Keeper.GetLastMerkleDropId(ctx) + 1
	m.Keeper.SetLastMerkleDropId(ctx, mdId)

	tAmt := sdk.NewIntFromUint64(msg.TotalAmount)

	merkledrop := types.Merkledrop{
		Id:          mdId,
		MerkleRoot:  msg.MerkleRoot,
		TotalAmount: tAmt,
		Owner:       msg.Owner,
	}
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	ctx.EventManager().EmitTypedEvent(&types.EventMerkledropCreate{
		Owner:        msg.Owner,
		MerkledropId: mdId,
	})

	return &types.MsgCreateMerkledropResponse{
		Owner: msg.Owner,
		Id:    mdId,
	}, nil
}
