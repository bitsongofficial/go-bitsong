// source: https://github.com/DA0-DA0/polytone/blob/main/tests/strangelove/incompatible_handshake_test.go
package e2e

import (
	"encoding/json"
	"fmt"
	"testing"

	w "github.com/CosmWasm/wasmvm/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	testBinary string = "aGVsbG8=" // "hello" in base64
	testText   string = "hello"
)

// Tests that a note may only ever connect to a voice, and a voice
// only to a note.
func TestPolytoneOnBitsong(t *testing.T) {
	suite := NewPolytoneSuite(t)

	// note <-> note not allowed.
	_, tc, err := suite.CreateChannel(
		suite.ChainA.Note,
		suite.ChainB.Note,
		&suite.ChainA,
		&suite.ChainB, suite.PathAB,
	)
	require.ErrorContains(t, err, "no new channels created", "note <-/-> note")
	fmt.Println("trychannel: ", tc)
	channels := suite.QueryChannelsInState(&suite.ChainB, CHANNEL_STATE_TRY)
	require.Len(t, channels, 1, "try note stops in first step")
	channels = suite.QueryChannelsInState(&suite.ChainB, CHANNEL_STATE_INIT)
	require.Len(t, channels, 1, "init note doesn't advance")

	// voice <-> voice not allowed
	_, _, err = suite.CreateChannel(
		suite.ChainA.Voice,
		suite.ChainB.Voice,
		&suite.ChainA,
		&suite.ChainB,
		suite.PathAB,
	)
	require.ErrorContains(t, err, "no new channels created", "voice <-/-> voice")
	accAddr, _ := sdk.AccAddressFromBech32(suite.ChainB.Tester)
	dataCosmosMsg, _ := HelloMessage(accAddr, string(testBinary))

	noDataCosmosMsg := w.CosmosMsg{
		Distribution: &w.DistributionMsg{
			SetWithdrawAddress: &w.SetWithdrawAddressMsg{
				Address: suite.ChainB.Voice,
			},
		},
	}

	callbackExecute, err := suite.RoundtripExecute(suite.ChainA.Note, &suite.ChainB, []w.CosmosMsg{dataCosmosMsg, noDataCosmosMsg})
	if err != nil {
		t.Fatal(err)
	}

	require.Len(t, len(callbackExecute.Success), 2, "error: "+callbackExecute.Error)
	require.Equal(t, "", callbackExecute.Error)
	// note <-> voice allowed
	//
	// FIXME: below errors with:
	//
	// `exit code 1:  Error: channel {channel-1} with port {wasm.juno14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9skjuwg8} already exists on chain {juno1-1}`
	//
	// See `TestHandshakeBetweenSameModule` where this channel
	// creation also fails in ibctesting.

	// _, _, err = suite.CreateChannel(
	// 	suite.ChainA.Note,
	// 	suite.ChainB.Voice,
	// 	&suite.ChainA,
	// 	&suite.ChainB,
	// )
	// require.NoError(t, err, "note <-> voice")

}
func HelloMessage(to sdk.AccAddress, data string) (w.CosmosMsg, error) {
	msgContent := map[string]interface{}{"hello": map[string]string{"data": data}}
	msgBytes, err := json.Marshal(msgContent)
	if err != nil {
		return w.CosmosMsg{}, fmt.Errorf("failed to marshal message: %w", err)
	}
	return w.CosmosMsg{
		Wasm: &w.WasmMsg{
			Execute: &w.ExecuteMsg{
				ContractAddr: to.String(),
				Msg:          msgBytes,
				Funds:        []w.Coin{},
			},
		},
	}, nil
}
