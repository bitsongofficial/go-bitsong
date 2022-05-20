package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateMerkledrop = "create_merkledrop"
)

var _ sdk.Msg = &MsgCreateMerkledrop{}

func NewMsgCreateMerkledrop(owner sdk.AccAddress, merkleRoot string, totalAmount uint64) *MsgCreateMerkledrop {
	return &MsgCreateMerkledrop{
		Owner:       owner.String(),
		MerkleRoot:  merkleRoot,
		TotalAmount: totalAmount,
	}
}

func (msg MsgCreateMerkledrop) Route() string { return RouterKey }

func (msg MsgCreateMerkledrop) Type() string { return TypeMsgCreateMerkledrop }

func (msg MsgCreateMerkledrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateMerkledrop) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateMerkledrop) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}
