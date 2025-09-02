package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_createCollectionDenom(t *testing.T) {
	creator := sdk.AccAddress(tmhash.SumTruncated([]byte("creator")))
	symbol := "MYNFT"

	expectedDenom := "nftF1D9FE89CCE1FAD3F83FFCBA6F496EFD30855C42"
	k := Keeper{}

	denom := k.createCollectionDenom(creator, symbol)
	if denom != expectedDenom {
		t.Errorf("expected %s, got %s", expectedDenom, denom)
	}
}
