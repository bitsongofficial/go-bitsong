package app

import (
	"github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/cosmos/cosmos-sdk/std"
)

var encodingConfig params.EncodingConfig = MakeEncodingConfig()

func GetEncodingConfig() params.EncodingConfig {
	return encodingConfig
}

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	AppModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	AppModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
