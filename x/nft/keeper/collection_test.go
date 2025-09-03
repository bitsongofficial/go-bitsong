package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_createCollectionDenom(t *testing.T) {
	creator := sdk.AccAddress(tmhash.SumTruncated([]byte("creator")))
	symbol := "MYNFT"

	expectedDenom := "nft653AF6715F0C4EE2E24A54B191EBD0AD5DB33723"
	k := Keeper{}

	denom := k.createCollectionDenom(creator, symbol)
	if denom != expectedDenom {
		t.Errorf("expected %s, got %s", expectedDenom, denom)
	}
}
