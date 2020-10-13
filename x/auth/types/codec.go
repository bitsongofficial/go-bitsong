package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	authtypes.RegisterCodec(cdc)

	cdc.RegisterConcrete(&BitSongAccount{}, "bitsong/Account", nil)
	cdc.RegisterConcrete(MsgRegisterHandle{}, "bitsong/MsgRegisterHandle", nil)
}

func RegisterAccountTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
