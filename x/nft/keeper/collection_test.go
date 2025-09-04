package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_createCollectionDenom(t *testing.T) {
	creator := sdk.AccAddress(tmhash.SumTruncated([]byte("creator1")))
	symbol := "MYNFT"

	expectedDenom := "nft9436DDD23FB751AEA7BC6C767F20F943DD735E06"
	k := Keeper{}

	denom, err := k.createCollectionDenom(creator, symbol)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if denom != expectedDenom {
		t.Errorf("expected %s, got %s", expectedDenom, denom)
	}
}
