package track

import (
	"fmt"
	"testing"
)

func TestEndBlocker(t *testing.T) {
	input := SetupTestInput(t)
	ctx := input.ctx
	trackKeeper := input.trackKeeper

	playPool := trackKeeper.GetFeePlayPool(ctx)
	fmt.Printf("%s", playPool)
}
