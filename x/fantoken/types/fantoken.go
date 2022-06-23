package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const FanTokenDecimal = 6

var (
	_ proto.Message = &FanToken{}
)

// NewFanToken constructs a new FanToken instance
func NewFanToken(name, symbol, uri string, maxSupply sdk.Int, minter, authority sdk.AccAddress, height int64) *FanToken {
	return &FanToken{
		Denom:     GetFantokenDenom(height, minter, symbol, name),
		MaxSupply: maxSupply,
		MetaData:  NewMetadata(name, symbol, uri, authority),
		Minter:    minter.String(),
	}
}

// GetSymbol implements exported.FanTokenI
func (ft FanToken) GetSymbol() string {
	return ft.MetaData.Symbol
}

// GetDenom implements exported.FanTokenI
func (ft FanToken) GetDenom() string {
	return ft.Denom
}

// GetName implements exported.FanTokenI
func (ft FanToken) GetName() string {
	return ft.MetaData.Name
}

// GetMaxSupply implements exported.FanTokenI
func (ft FanToken) GetMaxSupply() sdk.Int {
	return ft.MaxSupply
}

// GetMintable implements exported.FanTokenI
func (ft FanToken) GetMintable() bool {
	return ft.Minter != ""
}

// GetAuthority implements exported.FanTokenI
func (ft FanToken) GetAuthority() sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(ft.MetaData.Authority)
	return authority
}

// GetMinter implements exported.FanTokenI
func (ft FanToken) GetMinter() sdk.AccAddress {
	minter, _ := sdk.AccAddressFromBech32(ft.Minter)
	return minter
}

// GetURI implements exported.FanTokenI
func (ft FanToken) GetURI() string {
	return ft.MetaData.URI
}

// GetMetaData returns metadata of the fantoken
func (ft FanToken) GetMetaData() Metadata {
	return ft.MetaData
}

func (ft FanToken) String() string {
	bz, _ := yaml.Marshal(ft)
	return string(bz)
}

func (ft FanToken) Validate() error {
	if len(ft.Minter) > 0 {
		if _, err := sdk.AccAddressFromBech32(ft.Minter); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
		}
	}

	return ft.MetaData.Validate()
}

func (ft FanToken) ValidateWithDenom() error {
	if err := ft.Validate(); err != nil {
		return err
	}

	if err := ValidateDenom(ft.GetDenom()); err != nil {
		return err
	}
	return nil
}

// NewMetadata constructs a new FanToken Metadata instance
func NewMetadata(name, symbol, uri string, authority sdk.AccAddress) Metadata {
	return Metadata{
		Name:      name,
		Symbol:    symbol,
		URI:       uri,
		Authority: authority.String(),
	}
}

func (m Metadata) Validate() error {
	if len(m.Authority) > 0 {
		if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
		}
	}

	if err := ValidateName(m.Name); err != nil {
		return err
	}
	if err := ValidateSymbol(m.Symbol); err != nil {
		return err
	}
	if err := ValidateUri(m.URI); err != nil {
		return err
	}

	return nil
}
