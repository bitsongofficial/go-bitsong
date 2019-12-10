package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	TypeMsgCreateDistributor = "create_distributor"
)

var _ sdk.Msg = MsgCreateDistributor{}

// MsgCreateDistributor
type MsgCreateDistributor struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
}

func NewMsgCreateDistributor(name string, address sdk.AccAddress) MsgCreateDistributor {
	return MsgCreateDistributor{
		Name:    name,
		Address: address,
	}
}

//nolint
func (msg MsgCreateDistributor) Route() string { return RouterKey }
func (msg MsgCreateDistributor) Type() string  { return TypeMsgCreateDistributor }

// ValidateBasic
func (msg MsgCreateDistributor) ValidateBasic() sdk.Error {
	if len(strings.TrimSpace(msg.Name)) == 0 {
		return ErrInvalidDistributorName(DefaultCodespace, "distributor name cannot be blank")
	}

	if len(msg.Name) > MaxNameLength {
		return ErrInvalidDistributorName(DefaultCodespace, fmt.Sprintf("distributor name is longer than max length of %d", MaxNameLength))
	}

	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}

	return nil
}

// Implements Msg.
func (msg MsgCreateDistributor) String() string {
	return fmt.Sprintf(`Create Distributor Message:
  Name: %s
  Address: %s
`, msg.Name, msg.Address.String())
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateDistributor) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateDistributor) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
