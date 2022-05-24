package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMerkledropNotExist   = sdkerrors.Register(ModuleName, 1, "merkledrop does not exist")
	ErrInvalidMerkleRoot    = sdkerrors.Register(ModuleName, 2, "invalid merkle root")
	ErrInvalidCoin          = sdkerrors.Register(ModuleName, 3, "invalid coin")
	ErrAlreadyClaimed       = sdkerrors.Register(ModuleName, 4, "merkledrop already claimed")
	ErrInvalidMerkleProofs  = sdkerrors.Register(ModuleName, 5, "invalid merkle proofs")
	ErrTransferCoins        = sdkerrors.Register(ModuleName, 6, "error transfer coins")
	ErrInvalidOwner         = sdkerrors.Register(ModuleName, 7, "invalid owner")
	ErrInvalidSender        = sdkerrors.Register(ModuleName, 8, "invalid sender")
	ErrInvalidStartHeight   = sdkerrors.Register(ModuleName, 9, "invalid start height")
	ErrInvalidEndHeight     = sdkerrors.Register(ModuleName, 10, "invalid end height")
	ErrMerkledropNotBegun   = sdkerrors.Register(ModuleName, 11, "merkledrop not begun")
	ErrMerkledropExpired    = sdkerrors.Register(ModuleName, 12, "merkledrop expired")
	ErrMerkledropNotExpired = sdkerrors.Register(ModuleName, 13, "merkledrop not expired")
	ErrAlreadyWithdrawn     = sdkerrors.Register(ModuleName, 14, "funds have been already withdrawn")
	ErrCreationFee          = sdkerrors.Register(ModuleName, 15, "cannot deduct creation fee")
	ErrDeleteMerkledrop     = sdkerrors.Register(ModuleName, 16, "failed delete merkledrop")
)
