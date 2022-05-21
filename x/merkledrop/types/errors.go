package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMerkledropNotExist  = sdkerrors.Register(ModuleName, 1, "merkledrop does not exist")
	ErrInvalidMerkleRoot   = sdkerrors.Register(ModuleName, 2, "invalid merkle root")
	ErrInvalidCoin         = sdkerrors.Register(ModuleName, 3, "invalid coin")
	ErrAlreadyClaimed      = sdkerrors.Register(ModuleName, 4, "merkledrop already claimed")
	ErrInvalidMerkleProofs = sdkerrors.Register(ModuleName, 5, "invalid merkle proofs")
	ErrTransferCoins       = sdkerrors.Register(ModuleName, 6, "error transfer coins")
	ErrInvalidOwner        = sdkerrors.Register(ModuleName, 7, "invalid owner")
	ErrInvalidEndTime      = sdkerrors.Register(ModuleName, 8, "invalid end time")
	ErrMerkledropNotBegun  = sdkerrors.Register(ModuleName, 9, "merkledrop not begun")
	ErrMerkledropExpired   = sdkerrors.Register(ModuleName, 10, "merkledrop expired")
)
