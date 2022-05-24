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

func (m msgServer) Create(goCtx context.Context, msg *types.MsgCreate) (*types.MsgCreateResponse, error) {
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check start height > 0
	if msg.StartHeight == 0 {
		msg.StartHeight = ctx.BlockHeight()
	}

	// check end time and start time
	if msg.EndHeight <= msg.StartHeight {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndHeight, "end height must be > start height")
	}

	if msg.EndHeight <= ctx.BlockHeight() {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndHeight, "end height (%d) must be > current block height (%d)", msg.EndHeight, ctx.BlockHeight())
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

	// deduct creation fee
	fee, err := m.DeductCreationFee(ctx, owner)
	if err != nil {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrCreationFee, "creation-fee %s", fee.String())
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
		Id:          mdId,
		MerkleRoot:  msg.MerkleRoot,
		StartHeight: msg.StartHeight,
		EndHeight:   msg.EndHeight,
		Amount:      msg.Coin.Amount,
		Denom:       msg.Coin.Denom,
		Claimed:     sdk.ZeroInt(),
		Owner:       msg.Owner,
		Withdrawn:   false,
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
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrInvalidSender, "sender %s", sender.String())
	}

	// get merkledrop
	merkledrop, err := m.Keeper.GetMerkleDropById(ctx, msg.MerkledropId)
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotExist, "merkledrop: %d does not exist", msg.MerkledropId)
	}

	// merkledrop begun
	if merkledrop.StartHeight > ctx.BlockHeight() {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotBegun, "start-height %d, current-height %d", merkledrop.StartHeight, ctx.BlockHeight())
	}

	// merkledrop not expired
	if merkledrop.EndHeight < ctx.BlockHeight() {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropExpired, "end-height %d, current-height %d", merkledrop.EndHeight, ctx.BlockHeight())
	}

	// remaining funds are withdrawn
	if merkledrop.Withdrawn {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrAlreadyWithdrawn, "claim error")
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
	valid := types.IsValidProof(msg.Index, sender, msg.Amount, merkleRoot, proofs)
	if !valid {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrInvalidMerkleProofs, "invalid proofs")
	}

	// set claimed
	m.Keeper.SetClaimed(ctx, msg.MerkledropId, msg.Index)

	// add claimed amount
	merkledrop.Claimed = merkledrop.Claimed.Add(msg.Amount)
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// TODO: if claimed amount == total amount, then prune the merkledrop from the state

	// send coins
	coin := sdk.NewCoin(merkledrop.Denom, msg.Amount)
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{coin})
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "%s%s", msg.Amount, merkledrop.Denom)
	}

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventClaim{
		MerkledropId: merkledrop.Id,
		Index:        msg.Index,
		Coin:         coin,
	})

	return &types.MsgClaimResponse{
		Id:     0,
		Index:  0,
		Amount: msg.Amount,
	}, nil
}

func (m msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	// unwrap context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// decode owner
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrInvalidOwner, "owner %s", owner.String())
	}

	// get merkledrop
	merkledrop, err := m.Keeper.GetMerkleDropById(ctx, msg.Id)
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotExist, "merkledrop: %d does not exist", msg.Id)
	}

	// remaining funds are withdrawn
	if merkledrop.Withdrawn {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrAlreadyWithdrawn, "withdraw error")
	}

	// check owner
	if merkledrop.Owner != owner.String() {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrInvalidOwner, "unauthorized: %s", msg.Owner)
	}

	// make sure is expired
	if merkledrop.EndHeight > ctx.BlockHeight() {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotExpired, "end-height: %d, current-height: %d", merkledrop.EndHeight, ctx.BlockHeight())
	}

	// check if total amount < claimed amount  (who knows?)
	if merkledrop.Amount.LT(merkledrop.Claimed) {
		panic(fmt.Errorf("merkledrop-id: %d, total_amount (%s) < claimed_amount (%s)", merkledrop.Id, merkledrop.Amount, merkledrop.Claimed))
	}

	// set withdrawn flag
	merkledrop.Withdrawn = true
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// get balance
	balance := merkledrop.Amount.Sub(merkledrop.Claimed)

	// send coins
	coin := sdk.NewCoin(merkledrop.Denom, balance)
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.Coins{coin})
	if err != nil {
		return &types.MsgWithdrawResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "%s", coin)
	}

	// TODO: prune the merkledrop from the state

	// emit event
	ctx.EventManager().EmitTypedEvent(&types.EventWithdraw{
		MerkledropId: merkledrop.Id,
		Coin:         coin,
	})

	return &types.MsgWithdrawResponse{
		Id:   merkledrop.Id,
		Coin: coin,
	}, nil
}
