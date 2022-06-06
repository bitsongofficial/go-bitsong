package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

func GetFantokenDenom(height int64, creator sdk.AccAddress, symbol, name string) string {
	bz := []byte(fmt.Sprintf("%d%s%s%s", height, creator.String(), symbol, name))
	return "ft" + tmcrypto.AddressHash(bz).String()
}
