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
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check coin amount > 0
	if !msg.Coin.Amount.GT(sdk.ZeroInt()) {
		return &types.MsgCreateMerkledropResponse{}, sdkerrors.Wrapf(types.ErrInvalidCoin, "invalid coin amount, must be greater then zero")
	}

	// decode owner
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, err
	}

	// check and decode merkle root
	_, err = hex.DecodeString(msg.MerkleRoot)
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
	}

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.Coins{msg.Coin})
	if err != nil {
		return &types.MsgCreateMerkledropResponse{}, err
	}

	// increment merkledrop id
	mdId := m.Keeper.GetLastMerkleDropId(ctx) + 1
	m.Keeper.SetLastMerkleDropId(ctx, mdId)

	// set merkledrop
	merkledrop := types.Merkledrop{
		Id:         mdId,
		MerkleRoot: msg.MerkleRoot,
		Coin:       msg.Coin,
		Claimed:    sdk.Coin{Amount: sdk.ZeroInt(), Denom: msg.Coin.Denom},
		Owner:      msg.Owner,
	}
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// emit event
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
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// decode sender
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	// TODO: merkledrop begun
	// TODO: merkledrop not expired

	// check if is claimed
	isClaimed := m.Keeper.IsClaimed(ctx, msg.MerkledropId, msg.Index)
	if isClaimed {
		return &types.MsgClaimMerkledropResponse{}, fmt.Errorf("merkledrop already claimed")
	}

	// get merkledrop
	merkledrop, err := m.Keeper.GetMerkleDropById(ctx, msg.MerkledropId)
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	// check and decode merkle root
	merkleRoot, err := hex.DecodeString(merkledrop.GetMerkleRoot())
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
	}

	// verify proofs
	proofs := types.ConvertProofs(msg.Proofs)
	valid := types.IsValidProof(msg.Index, sender, msg.Coin.Amount, merkleRoot, proofs)
	if !valid {
		return &types.MsgClaimMerkledropResponse{}, fmt.Errorf("invalid proofs")
	}

	// set claimed
	m.Keeper.SetClaimed(ctx, msg.MerkledropId, msg.Index)

	// add claimed amount
	merkledrop.Claimed = merkledrop.Claimed.Add(msg.Coin)
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{msg.Coin})
	if err != nil {
		return &types.MsgClaimMerkledropResponse{}, err
	}

	return &types.MsgClaimMerkledropResponse{}, nil
}
