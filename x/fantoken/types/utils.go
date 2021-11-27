package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func GetFantokenDenom(creator sdk.AccAddress, symbol, name string) string {
	return "ft" + crypto.AddressHash([]byte(creator.String()+symbol+name)).String()
}
