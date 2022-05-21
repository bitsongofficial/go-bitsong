package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateMerkledrop = "create_merkledrop"
	TypeMsgClaimMerkledrop  = "claim_merkledrop"
)

var _ sdk.Msg = &MsgCreateMerkledrop{}

func NewMsgCreateMerkledrop(owner sdk.AccAddress, merkleRoot string, coin sdk.Coin) *MsgCreateMerkledrop {
	return &MsgCreateMerkledrop{
		Owner:      owner.String(),
		MerkleRoot: merkleRoot,
		Coin:       coin,
	}
}

func (msg MsgCreateMerkledrop) Route() string { return RouterKey }

func (msg MsgCreateMerkledrop) Type() string { return TypeMsgCreateMerkledrop }

func (msg MsgCreateMerkledrop) ValidateBasic() error {
	if msg.Coin.Amount.LTE(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidCoin, "invalid coin amount, less then zero")
	}

	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	_, err = hex.DecodeString(msg.MerkleRoot)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidMerkleRoot, "invalid merkle root (%s)", err)
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

var _ sdk.Msg = &MsgClaimMerkledrop{}

func NewMsgClaimMerkledrop(index, mdId uint64, coin sdk.Coin, proofs []string, sender sdk.AccAddress) *MsgClaimMerkledrop {
	return &MsgClaimMerkledrop{
		Index:        index,
		MerkledropId: mdId,
		Coin:         coin,
		Proofs:       proofs,
		Sender:       sender.String(),
	}
}

func (msg MsgClaimMerkledrop) Route() string { return RouterKey }

func (msg MsgClaimMerkledrop) Type() string { return TypeMsgClaimMerkledrop }

func (msg MsgClaimMerkledrop) ValidateBasic() error {
	for _, p := range msg.Proofs {
		_, err := hex.DecodeString(p)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidMerkleRoot, "invalid merkle proof (%s)", err)
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClaimMerkledrop) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClaimMerkledrop) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
