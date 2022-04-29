package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateNFT                 = "create_nft"
	TypeMsgTransferNFT               = "transfer_nft"
	TypeMsgSignMetadata              = "sign_metadata"
	TypeMsgUpdateMetadata            = "update_metadata"
	TypeMsgUpdateMetadataAuthority   = "update_metadata_authority"
	TypeMsgCreateCollection          = "create_collection"
	TypeMsgVerifyCollection          = "verify_collection"
	TypeMsgUnverifyCollection        = "unverify_collection"
	TypeMsgUpdateCollectionAuthority = "update_collection_authority"
)

var _ sdk.Msg = &MsgCreateNFT{}

func NewMsgCreateNFT(sender sdk.AccAddress, updateAuthority string, data Data, presaleHappened, isMutable bool) *MsgCreateNFT {
	return &MsgCreateNFT{
		Sender: sender.String(),
		Metadata: Metadata{
			UpdateAuthority:     updateAuthority,
			Mint:                sender.String(),
			Data:                &data,
			PrimarySaleHappened: presaleHappened,
			IsMutable:           isMutable,
		},
	}
}

func (msg MsgCreateNFT) Route() string { return RouterKey }

func (msg MsgCreateNFT) Type() string { return TypeMsgCreateNFT }

func (msg MsgCreateNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateNFT) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgTransferNFT{}

func NewMsgTransferNFT() *MsgTransferNFT {
	return &MsgTransferNFT{}
}

func (msg MsgTransferNFT) Route() string { return RouterKey }

func (msg MsgTransferNFT) Type() string { return TypeMsgTransferNFT }

func (msg MsgTransferNFT) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferNFT) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgSignMetadata{}

func NewMsgSignMetadata() *MsgSignMetadata {
	return &MsgSignMetadata{}
}

func (msg MsgSignMetadata) Route() string { return RouterKey }

func (msg MsgSignMetadata) Type() string { return TypeMsgSignMetadata }

func (msg MsgSignMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSignMetadata) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSignMetadata) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateMetadata{}

func NewMsgUpdateMetadata() *MsgUpdateMetadata {
	return &MsgUpdateMetadata{}
}

func (msg MsgUpdateMetadata) Route() string { return RouterKey }

func (msg MsgUpdateMetadata) Type() string { return TypeMsgUpdateMetadata }

func (msg MsgUpdateMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateMetadata) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateMetadata) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateMetadataAuthority{}

func NewMsgUpdateMetadataAuthority() *MsgUpdateMetadataAuthority {
	return &MsgUpdateMetadataAuthority{}
}

func (msg MsgUpdateMetadataAuthority) Route() string { return RouterKey }

func (msg MsgUpdateMetadataAuthority) Type() string { return TypeMsgUpdateMetadataAuthority }

func (msg MsgUpdateMetadataAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateMetadataAuthority) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateMetadataAuthority) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgCreateCollection{}

func NewMsgCreateCollection() *MsgCreateCollection {
	return &MsgCreateCollection{}
}

func (msg MsgCreateCollection) Route() string { return RouterKey }

func (msg MsgCreateCollection) Type() string { return TypeMsgCreateCollection }

func (msg MsgCreateCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateCollection) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateCollection) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgVerifyCollection{}

func NewMsgVerifyCollection() *MsgVerifyCollection {
	return &MsgVerifyCollection{}
}

func (msg MsgVerifyCollection) Route() string { return RouterKey }

func (msg MsgVerifyCollection) Type() string { return TypeMsgVerifyCollection }

func (msg MsgVerifyCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgVerifyCollection) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgVerifyCollection) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUnverifyCollection{}

func NewMsgUnverifyCollection() *MsgUnverifyCollection {
	return &MsgUnverifyCollection{}
}

func (msg MsgUnverifyCollection) Route() string { return RouterKey }

func (msg MsgUnverifyCollection) Type() string { return TypeMsgUnverifyCollection }

func (msg MsgUnverifyCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUnverifyCollection) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUnverifyCollection) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgUpdateCollectionAuthority{}

func NewMsgUpdateCollectionAuthority() *MsgUpdateCollectionAuthority {
	return &MsgUpdateCollectionAuthority{}
}

func (msg MsgUpdateCollectionAuthority) Route() string { return RouterKey }

func (msg MsgUpdateCollectionAuthority) Type() string { return TypeMsgUpdateCollectionAuthority }

func (msg MsgUpdateCollectionAuthority) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUpdateCollectionAuthority) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgUpdateCollectionAuthority) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
