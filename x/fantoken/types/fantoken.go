package types

import (
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
		Denom:     GetFantokenDenom(height, authority, symbol, name),
		MaxSupply: maxSupply,
		MetaData:  NewMetadata(name, symbol, uri, authority),
		Minter:    minter.String(),
	}
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
