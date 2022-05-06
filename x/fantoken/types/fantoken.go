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

// FanTokenI defines an interface for FanToken
type FanTokenI interface {
	GetName() string
	GetSymbol() string
	GetDenom() string
	GetUri() string
	GetMaxSupply() sdk.Int
	GetMintable() bool
	GetOwner() sdk.AccAddress
	GetMetaData() Metadata
}

// NewFanToken constructs a new FanToken instance
func NewFanToken(name, symbol, uri string, maxSupply sdk.Int, owner sdk.AccAddress) FanToken {
	return FanToken{
		Denom:     GetFantokenDenom(owner, symbol, name),
		MaxSupply: maxSupply,
		Mintable:  true,
		Owner:     owner.String(),
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

// GetOwner implements exported.FanTokenI
func (ft FanToken) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(ft.Owner)
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
