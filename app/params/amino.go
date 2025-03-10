//go:build test_amino
// +build test_amino

package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	cdc := codec.New()
	interfaceRegistry := testutil.CodecOptions{AccAddressPrefix: "bitsong", ValAddressPrefix: "bitsongvaloper"}.NewInterfaceRegistry()
	marshaler := codec.NewAminoCodec(cdc)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          authtypes.StdTxConfig{Cdc: cdc},
		Amino:             cdc,
	}
}
