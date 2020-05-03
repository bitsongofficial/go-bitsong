package types

import (
	"fmt"
	btsg "github.com/bitsongofficial/go-bitsong/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterPlayer = "register_player"
)

var _ sdk.Msg = MsgRegisterPlayer{}

type MsgRegisterPlayer struct {
	Moniker string         `json:"moniker" yaml:"moniker"`
	Deposit sdk.Coin       `json:"deposit" yaml:"deposit"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
}

func NewMsgRegisterPlayer(moniker string, deposit sdk.Coin, from sdk.AccAddress) MsgRegisterPlayer {
	return MsgRegisterPlayer{
		Moniker: moniker,
		Deposit: deposit,
		Owner:   from,
	}
}

func (msg MsgRegisterPlayer) Route() string { return RouterKey }
func (msg MsgRegisterPlayer) Type() string  { return TypeMsgRegisterPlayer }

func (msg MsgRegisterPlayer) ValidateBasic() error {
	if msg.Deposit.Denom != btsg.BondDenom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invelid deposit"))
	}

	if msg.Moniker == "" || len(msg.Moniker) < MinMonikerLength || len(msg.Moniker) > MaxMonikerLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid moniker"))
	}

	if msg.Owner == nil || msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid owner"))
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterPlayer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterPlayer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func (msg MsgRegisterPlayer) String() string {
	return fmt.Sprintf(`Msg Register Player
Moniker: %s
Deposit: %s
Owner:  %s`,
		msg.Moniker, msg.Deposit, msg.Owner,
	)
}
