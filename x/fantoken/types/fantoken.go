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
func NewFanToken(name, symbol, uri string, maxSupply sdk.Int, authority sdk.AccAddress, height int64) *FanToken {
	return &FanToken{
		Denom:     GetFantokenDenom(height, authority, symbol, name),
		MaxSupply: maxSupply,
		Mintable:  true,
		Authority: authority.String(),
		MetaData:  NewMetadata(name, symbol, uri),
	}
}

// NewMetadata constructs a new FanToken Metadata instance
func NewMetadata(name, symbol, uri string) Metadata {
	return Metadata{
		Name:   name,
		Symbol: symbol,
		URI:    uri,
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
	return ft.Mintable
}

// GetAuthority implements exported.FanTokenI
func (ft FanToken) GetAuthority() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(ft.Authority)
	return owner
}

// GetUri implements exported.FanTokenI
func (ft FanToken) GetUri() string {
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
