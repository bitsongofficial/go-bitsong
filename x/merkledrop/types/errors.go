package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrMerkleDropNotExist = sdkerrors.Register(ModuleName, 1, "merkledrop does not exist")
	ErrInvalidMerkleRoot  = sdkerrors.Register(ModuleName, 2, "invalid merkle root")
)
