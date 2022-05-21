package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, err
	}

	_, err = hex.DecodeString(msg.MerkleRoot)
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
	}

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.Coins{msg.TotalAmount})
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, err
	}

	// create merkledrop
	mdId := m.Keeper.GetLastMerkleDropId(ctx) + 1
	m.Keeper.SetLastMerkleDropId(ctx, mdId)

	mRoot, err := hex.DecodeString(msg.MerkleRoot)
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, err
	}

	merkledrop := types.Merkledrop{
		Id:          mdId,
		MerkleRoot:  mRoot,
		TotalAmount: msg.TotalAmount,
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

func (m msgServer) ClaimMerkledrop(goCtx context.Context, msg *types.MsgClaimMerkledrop) (*types.MsgClaimMerkledropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	isClaimed := m.Keeper.IsClaimed(ctx, msg.MerkledropId, msg.Index)
	if isClaimed {
		return &types.MsgClaimMerkledropResponse{}, fmt.Errorf("merkledrop already claimed")
	}

	merkledrop, err := m.Keeper.GetMerkleDropById(ctx, msg.MerkledropId)
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	proofs := types.ConvertProofs(msg.Proofs)
	valid := types.IsValidProof(msg.Index, sender, msg.Amount.Amount, merkledrop.GetMerkleRoot(), proofs)

	if !valid {
		return &types.MsgClaimMerkledropResponse{}, fmt.Errorf("invalid proofs")
	}

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{msg.Amount})
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	// set claimed
	m.Keeper.SetClaimed(ctx, msg.MerkledropId, msg.Index)

	return &types.MsgClaimMerkledropResponse{}, nil
}
