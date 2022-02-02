package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

func GetFantokenDenom(creator sdk.AccAddress, symbol, name string) string {
	return "ft" + tmcrypto.AddressHash([]byte(creator.String()+symbol+name)).String()
}
