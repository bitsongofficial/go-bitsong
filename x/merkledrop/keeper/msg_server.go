package keeper

import (
	"context"
	"encoding/hex"
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

func (m msgServer) Create(goCtx context.Context, msg *types.MsgCreate) (*types.MsgCreateResponse, error) {
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check end time and start time
	if msg.EndTime.Before(msg.StartTime) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndTime, "end time must be after start time")
	}

	if msg.EndTime.Before(ctx.BlockTime()) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndTime, "end time must be in the future")
	}

	// check coin amount > 0
	if !msg.Coin.Amount.GT(sdk.ZeroInt()) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidCoin, "invalid coin amount, must be greater then zero")
	}

	// decode owner
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidOwner, "owner %s", owner.String())
	}

	// check and decode merkle root
	_, err = hex.DecodeString(msg.MerkleRoot)
	if err != nil {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
	}

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.Coins{msg.Coin})
	if err != nil {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "%s", msg.Coin)
	}

	// increment merkledrop id
	mdId := m.Keeper.GetLastMerkleDropId(ctx) + 1
	m.Keeper.SetLastMerkleDropId(ctx, mdId)

	// set merkledrop
	merkledrop := types.Merkledrop{
		Id:         mdId,
		MerkleRoot: msg.MerkleRoot,
		StartTime:  msg.StartTime,
		EndTime:    msg.EndTime,
		Coin:       msg.Coin,
		Claimed:    sdk.Coin{Amount: sdk.ZeroInt(), Denom: msg.Coin.Denom},
		Owner:      msg.Owner,
	}
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventCreate{
		Owner:        msg.Owner,
		MerkledropId: mdId,
	})

	return &types.MsgCreateResponse{
		Owner: msg.Owner,
		Id:    mdId,
	}, nil
}

func (m msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// decode sender
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return &types.MsgClaimResponse{}, err
	}

	// get merkledrop
	merkledrop, err := m.Keeper.GetMerkleDropById(ctx, msg.MerkledropId)
	if err != nil {
		return &types.MsgClaimResponse{}, err
	}

	// merkledrop begun
	if merkledrop.StartTime.After(ctx.BlockTime()) {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotBegun, "start-time %s", merkledrop.StartTime.String())
	}

	// merkledrop not expired
	if merkledrop.EndTime.Before(ctx.BlockTime()) {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropExpired, "end-time %s", merkledrop.EndTime.String())
	}

	// check if is claimed
	isClaimed := m.Keeper.IsClaimed(ctx, msg.MerkledropId, msg.Index)
	if isClaimed {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrAlreadyClaimed, "merkledrop_id (%d)", msg.MerkledropId)
	}

	// check and decode merkle root
	merkleRoot, err := hex.DecodeString(merkledrop.GetMerkleRoot())
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
	}

	// verify proofs
	proofs := types.ConvertProofs(msg.Proofs)
	valid := types.IsValidProof(msg.Index, sender, msg.Coin.Amount, merkleRoot, proofs)
	if !valid {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleProofs, "invalid proofs")
	}

	// set claimed
	m.Keeper.SetClaimed(ctx, msg.MerkledropId, msg.Index)

	// add claimed amount
	merkledrop.Claimed = merkledrop.Claimed.Add(msg.Coin)
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// send coins
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{msg.Coin})
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "%s", msg.Coin)
	}

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventClaim{
		MerkledropId: merkledrop.Id,
		Index:        msg.Index,
		Coin:         msg.Coin,
	})

	return &types.MsgClaimResponse{}, nil
}
