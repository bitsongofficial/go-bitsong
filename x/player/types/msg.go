package types

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterPlayer = "register_player"
)

var _ sdk.Msg = MsgRegisterPlayer{}

type MsgRegisterPlayer struct {
	Moniker    string         `json:"moniker" yaml:"moniker"`
	PlayerAddr sdk.AccAddress `json:"player_addr" yaml:"player_addr"`
	Validator  sdk.ValAddress `json:"validator" yaml:"validator"`
}

func NewMsgRegisterPlayer(moniker string, plAddr sdk.AccAddress, from sdk.ValAddress) MsgRegisterPlayer {
	return MsgRegisterPlayer{
		Moniker:    moniker,
		PlayerAddr: plAddr,
		Validator:  from,
	}
}

func (msg MsgRegisterPlayer) Route() string { return RouterKey }
func (msg MsgRegisterPlayer) Type() string  { return TypeMsgRegisterPlayer }

func (msg MsgRegisterPlayer) ValidateBasic() error {
	if msg.Moniker == "" || len(msg.Moniker) < MinMonikerLength || len(msg.Moniker) > MaxMonikerLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid moniker"))
	}

	if msg.PlayerAddr == nil || msg.PlayerAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid player address"))
	}

	if msg.Validator == nil || msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid validator"))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterPlayer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as player's, then the validator must
// sign the msg as well.
func (msg MsgRegisterPlayer) GetSigners() []sdk.AccAddress {
	// player is first signer so player pays fees
	addrs := []sdk.AccAddress{msg.PlayerAddr}

	if !bytes.Equal(msg.PlayerAddr.Bytes(), msg.Validator.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(msg.Validator))
	}

	return addrs
}

func (msg MsgRegisterPlayer) String() string {
	return fmt.Sprintf(`Msg Register Player
Moniker: %s
Player: %s,
Validator:  %s`,
		msg.Moniker, msg.PlayerAddr, msg.Validator,
	)
}
