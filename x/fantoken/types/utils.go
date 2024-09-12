package types

import (
	"fmt"

	tmcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetFantokenDenom(height int64, minter sdk.AccAddress, symbol, name string) string {
	bz := []byte(fmt.Sprintf("%d%s%s%s", height, minter.String(), symbol, name))
	return "ft" + tmcrypto.AddressHash(bz).String()
}
