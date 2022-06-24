package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreate   = "create"
	TypeMsgClaim    = "claim"
	TypeMsgWithdraw = "withdraw"
)

var _ sdk.Msg = &MsgCreate{}

func NewMsgCreate(owner sdk.AccAddress, merkleRoot string, startHeight, endHeight int64, coin sdk.Coin) *MsgCreate {
	return &MsgCreate{
		Owner:       owner.String(),
		MerkleRoot:  merkleRoot,
		StartHeight: startHeight,
		EndHeight:   endHeight,
		Coin:        coin,
	}
}

func (msg MsgCreate) Route() string { return RouterKey }

func (msg MsgCreate) Type() string { return TypeMsgCreate }

func (msg MsgCreate) ValidateBasic() error {
	if msg.EndHeight <= msg.StartHeight {
		return sdkerrors.Wrapf(ErrInvalidEndHeight, "end height must be > start height")
	}

	if err := msg.Coin.Validate(); err != nil {
		return err
	}

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
func (msg MsgCreate) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreate) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

var _ sdk.Msg = &MsgClaim{}

func NewMsgClaim(index, mdId uint64, amount sdk.Int, proofs []string, sender sdk.AccAddress) *MsgClaim {
	return &MsgClaim{
		Index:        index,
		MerkledropId: mdId,
		Amount:       amount,
		Proofs:       proofs,
		Sender:       sender.String(),
	}
}

func (msg MsgClaim) Route() string { return RouterKey }

func (msg MsgClaim) Type() string { return TypeMsgClaim }

func (msg MsgClaim) ValidateBasic() error {
	for _, p := range msg.Proofs {
		_, err := hex.DecodeString(p)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidMerkleRoot, "invalid merkle proof (%s)", err)
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgClaim) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgClaim) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgWithdraw{}

func NewMsgWithdraw(owner sdk.AccAddress, merkledropID uint64) *MsgWithdraw {
	return &MsgWithdraw{
		Owner: owner.String(),
		Id:    merkledropID,
	}
}

func (msg MsgWithdraw) Route() string { return RouterKey }

func (msg MsgWithdraw) Type() string { return TypeMsgWithdraw }

func (msg MsgWithdraw) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdraw) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}
