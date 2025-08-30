package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

var _ types.MsgServer = &msgServer{}

// msgServer is a wrapper of Keeper.
type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the x/cadence MsgServer interface.
func NewMsgServerImpl(k *Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

// RegisterCadenceContract handles incoming transactions to register cadence contract s.
func (k msgServer) RegisterCadenceContract(goCtx context.Context, req *types.MsgRegisterCadenceContract) (*types.MsgRegisterCadenceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgRegisterCadenceContractResponse{}, k.RegisterContract(ctx, req.SenderAddress, req.ContractAddress)
}

// UnregisterCadenceContract handles incoming transactions to unregister cadence contract s.
func (k msgServer) UnregisterCadenceContract(goCtx context.Context, req *types.MsgUnregisterCadenceContract) (*types.MsgUnregisterCadenceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgUnregisterCadenceContractResponse{}, k.UnregisterContract(ctx, req.SenderAddress, req.ContractAddress)
}

// UnjailCadenceContract handles incoming transactions to unjail cadence contract s.
func (k msgServer) UnjailCadenceContract(goCtx context.Context, req *types.MsgUnjailCadenceContract) (*types.MsgUnjailCadenceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgUnjailCadenceContractResponse{}, k.SetJailStatusBySender(ctx, req.SenderAddress, req.ContractAddress, false)
}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
