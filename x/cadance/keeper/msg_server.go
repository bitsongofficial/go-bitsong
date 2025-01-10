package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

var _ types.MsgServer = &msgServer{}

// msgServer is a wrapper of Keeper.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the x/cadance MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

// RegisterCadanceContract handles incoming transactions to register cadance contract s.
func (k msgServer) RegisterCadanceContract(goCtx context.Context, req *types.MsgRegisterCadanceContract) (*types.MsgRegisterCadanceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgRegisterCadanceContractResponse{}, k.RegisterContract(ctx, req.SenderAddress, req.ContractAddress)
}

// UnregisterCadanceContract handles incoming transactions to unregister cadance contract s.
func (k msgServer) UnregisterCadanceContract(goCtx context.Context, req *types.MsgUnregisterCadanceContract) (*types.MsgUnregisterCadanceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgUnregisterCadanceContractResponse{}, k.UnregisterContract(ctx, req.SenderAddress, req.ContractAddress)
}

// UnjailCadanceContract handles incoming transactions to unjail cadance contract s.
func (k msgServer) UnjailCadanceContract(goCtx context.Context, req *types.MsgUnjailCadanceContract) (*types.MsgUnjailCadanceContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	return &types.MsgUnjailCadanceContractResponse{}, k.SetJailStatusBySender(ctx, req.SenderAddress, req.ContractAddress, false)
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
