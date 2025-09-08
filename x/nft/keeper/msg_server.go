package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return msgServer{Keeper: keeper}
}

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	denom, err := k.Keeper.CreateCollection(ctx, msg.Creator, msg.Minter, msg.Authority, msg.Symbol, msg.Name, msg.Uri)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateCollectionResponse{
		Denom: denom,
	}, nil
}

func (k msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minter, err := k.ac.StringToBytes(msg.Minter)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address: %s", err)
	}

	recipient, err := k.ac.StringToBytes(msg.Recipient)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address: %s", err)
	}

	err = k.Keeper.MintNFT(ctx, minter, recipient, msg.Collection, msg.TokenId, msg.Name, msg.Uri)
	if err != nil {
		return nil, err
	}

	return &types.MsgMintNFTResponse{
		Collection: msg.Collection,
		TokenId:    msg.TokenId,
	}, nil
}

func (k msgServer) SendNFT(goCtx context.Context, msg *types.MsgSendNFT) (*types.MsgSendNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := k.ac.StringToBytes(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	recipient, err := k.ac.StringToBytes(msg.Recipient)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address: %s", err)
	}

	err = k.Keeper.SendNFT(ctx, sender, recipient, msg.Collection, msg.TokenId)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendNFTResponse{}, nil
}

func (k msgServer) PrintEdition(goCtx context.Context, msg *types.MsgPrintEdition) (*types.MsgPrintEditionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	minter, err := k.ac.StringToBytes(msg.Minter)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address: %s", err)
	}

	recipient, err := k.ac.StringToBytes(msg.Recipient)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address: %s", err)
	}

	seq, err := k.Keeper.PrintEdition(ctx, minter, recipient, msg.Collection, msg.TokenId)
	if err != nil {
		return nil, err
	}

	return &types.MsgPrintEditionResponse{
		Collection: msg.Collection,
		TokenId:    msg.TokenId,
		Seq:        seq,
	}, nil
}
