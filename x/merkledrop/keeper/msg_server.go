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

	startHeight := sdk.NewInt(msg.StartHeight)
	endHeight := sdk.NewInt(msg.EndHeight)

	// check end height and start height
	if startHeight.IsNegative() {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidStartHeight, "start height must be greater then zero")
	}

	// check start height > current height
	if startHeight.LT(sdk.NewInt(ctx.BlockHeight())) {
		msg.StartHeight = ctx.BlockHeight()
	}

	// check end height and start height
	if endHeight.LTE(startHeight) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndHeight, "end height must be > start height")
	}

	if endHeight.LTE(sdk.NewInt(ctx.BlockHeight())) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndHeight, "end height (%d) must be > current block height (%d)", msg.EndHeight, ctx.BlockHeight())
	}

	// add check startheight
	// - max-start-height = blockheight + 100_000
	maxStartHeight := ctx.BlockHeight() + int64(100_000)

	// - max-end-height = msg.StartHeight + 5_000_000
	maxEndHeight := msg.StartHeight + int64(5_000_000)

	// start-height > max-start-height: return error
	if startHeight.GT(sdk.NewInt(maxStartHeight)) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidStartHeight, "start height is > block-height + 100000")
	}

	// end-height > max-end-height: return error
	if endHeight.GT(sdk.NewInt(maxEndHeight)) {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidEndHeight, "end height is > msg.StartHeight + 5000000")
	}

	// validate coin
	if err := msg.Coin.Validate(); err != nil {
		return &types.MsgCreateResponse{}, err
	}

	// check coin amount > 0
	if msg.Coin.Amount.LTE(sdk.ZeroInt()) {
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

	// deduct creation fee
	if err = m.DeductCreationFee(ctx, owner); err != nil {
		return nil, err
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
	}
	if err := m.Keeper.SetMerkleDrop(ctx, merkledrop); err != nil {
		return &types.MsgCreateResponse{}, sdkerrors.Wrapf(types.ErrInvalidSender, "sender %s", owner.String())
	}

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
	merkledrop, err := m.Keeper.getMerkleDropById(ctx, msg.MerkledropId)
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotExist, "merkledrop: %d does not exist", msg.MerkledropId)
	}

	startHeight := sdk.NewInt(merkledrop.StartHeight)
	endHeight := sdk.NewInt(merkledrop.EndHeight)

	// merkledrop begun
	if startHeight.GT(sdk.NewInt(ctx.BlockHeight())) {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropNotBegun, "start-height %d, current-height %d", merkledrop.StartHeight, ctx.BlockHeight())
	}

	// merkledrop not expired, last block is included
	if endHeight.LTE(sdk.NewInt(ctx.BlockHeight())) {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrMerkledropExpired, "end-height %d, current-height %d", merkledrop.EndHeight, ctx.BlockHeight())
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

	amtAvailable := merkledrop.Amount.Sub(merkledrop.Claimed)
	if amtAvailable.LT(msg.Amount) {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "something went wrong")
	}

	// send coins
	coin := sdk.NewCoin(merkledrop.Denom, msg.Amount)
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{coin})
	if err != nil {
		return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrTransferCoins, "%s%s", msg.Amount, merkledrop.Denom)
	}

	// set claimed
	m.Keeper.SetClaimed(ctx, msg.MerkledropId, msg.Index)

	// add claimed amount
	merkledrop.Claimed = merkledrop.Claimed.Add(msg.Amount)
	m.Keeper.SetMerkleDrop(ctx, merkledrop)

	// if claimed amount == total amount, then prune the merkledrop from the state
	if merkledrop.Claimed.Equal(merkledrop.Amount) {
		err := m.Keeper.DeleteMerkledropByID(ctx, merkledrop.Id)
		if err != nil {
			return &types.MsgClaimResponse{}, sdkerrors.Wrapf(types.ErrDeleteMerkledrop, err.Error())
		}
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
