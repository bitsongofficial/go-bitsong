package types

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Sudo Message called on the contracts
	EndBlockSudoMessage = `{"cadence_end_block":{}}`
)

// == MsgUpdateParams ==
const (
	TypeMsgRegisterCadenceContract   = "register_cadence_contract"
	TypeMsgUnregisterCadenceContract = "unregister_cadence_contract"
	TypeMsgUnjailCadenceContract     = "unjail_cadence_contract"
	TypeMsgUpdateParams              = "update_cadence_params"
)

var (
	_ sdk.Msg = &MsgRegisterCadenceContract{}
	_ sdk.Msg = &MsgUnregisterCadenceContract{}
	_ sdk.Msg = &MsgUnjailCadenceContract{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// Route returns the name of the module
func (msg MsgRegisterCadenceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterCadenceContract) Type() string { return TypeMsgRegisterCadenceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterCadenceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterCadenceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterCadenceContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnregisterCadenceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnregisterCadenceContract) Type() string { return TypeMsgUnregisterCadenceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnregisterCadenceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnregisterCadenceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnregisterCadenceContract) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// Route returns the name of the module
func (msg MsgUnjailCadenceContract) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUnjailCadenceContract) Type() string { return TypeMsgUnjailCadenceContract }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnjailCadenceContract) ValidateBasic() error {
	return validateAddresses(msg.SenderAddress, msg.ContractAddress)
}

// GetSignBytes encodes the message for signing
func (msg *MsgUnjailCadenceContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnjailCadenceContract) GetSigners() []sdk.AccAddress {
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
