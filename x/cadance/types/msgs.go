package types

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Sudo Message called on the contracts
	EndBlockSudoMessage = `{"cadance_end_block":{}}`
)

// == MsgUpdateParams ==
const (
	TypeMsgRegisterCadanceContract   = "register_cadance_contract"
	TypeMsgUnregisterCadanceContract = "unregister_cadance_contract"
	TypeMsgUnjailCadanceContract     = "unjail_cadance_contract"
	TypeMsgUpdateParams              = "update_cadance_params"
)

var (
	_ sdk.Msg = &MsgRegisterCadanceContract{}
	_ sdk.Msg = &MsgUnregisterCadanceContract{}
	_ sdk.Msg = &MsgUnjailCadanceContract{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// Route returns the name of the module
func (msg MsgRegisterCadanceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterCadanceContract) Type() string { return TypeMsgRegisterCadanceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterCadanceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterCadanceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterCadanceContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnregisterCadanceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnregisterCadanceContract) Type() string { return TypeMsgUnregisterCadanceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnregisterCadanceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnregisterCadanceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnregisterCadanceContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnjailCadanceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnjailCadanceContract) Type() string { return TypeMsgUnjailCadanceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnjailCadanceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnjailCadanceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnjailCadanceContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// NewMsgUpdateParams creates new instance of MsgUpdateParams
func NewMsgUpdateParams(
	sender sdk.Address,
	contractGasLimit uint64,
) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: sender.String(),
		Params:    NewParams(contractGasLimit),
	}
}

// Route returns the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return msg.Params.Validate()
}

// ValidateAddresses validates the provided addresses
func validateAddresses(addresses ...string) error {
	for _, address := range addresses {
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			return errors.Wrapf(err, "invalid address: %s", address)
		}
	}

	return nil
}
