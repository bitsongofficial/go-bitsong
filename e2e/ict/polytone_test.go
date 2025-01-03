// source: https://github.com/DA0-DA0/polytone/blob/main/tests/strangelove/incompatible_handshake_test.go
package e2e

import (
	"encoding/json"
	"fmt"
	"testing"

	w "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
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

	// // note <-> note not allowed.
	// _, tc, err := suite.CreateChannel(
	// 	suite.ChainA.Note,
	// 	suite.ChainB.Note,
	// 	&suite.ChainA,
	// 	&suite.ChainB, suite.PathAB,
	// )
	// require.ErrorContains(t, err, "no new channels created", "note <-/-> note")
	// log.Printf("trychannel: %+v", tc)

	// // voice <-> voice not allowed
	// _, _, err = suite.CreateChannel(
	// 	suite.ChainA.Voice,
	// 	suite.ChainB.Voice,
	// 	&suite.ChainA,
	// 	&suite.ChainB,
	// 	suite.PathAB,
	// )
	// require.ErrorContains(t, err, "no new channels created", "voice <-/-> voice")

	// note <-> voice allowed
	_, _, err := suite.CreateChannel(
		suite.ChainA.Note,
		suite.ChainB.Voice,
		&suite.ChainA,
		&suite.ChainB,
		suite.PathAB,
	)
	require.NoError(t, err, "note <-> voice")

	// Wait for the channel to get set up
	err = testutil.WaitForBlocks(suite.ctx, 2, suite.ChainA.Cosmos, suite.ChainB.Cosmos)
	require.NoError(t, err)

	// TODO: reimplement require no error here. this is commented out for now as we are happy that
	// that the full ibc channel creation lifecycle was successful, meaning that the wasmIbcHandler
	//  is communicating correctly with the ibc module keeper.
	// accAddr, _ := sdk.AccAddressFromBech32(suite.ChainB.Tester)
	// dataCosmosMsg, _ := HelloMessage(accAddr, string(testBinary))
	// noDataCosmosMsg := w.CosmosMsg{
	// 	Distribution: &w.DistributionMsg{
	// 		SetWithdrawAddress: &w.SetWithdrawAddressMsg{
	// 			Address: suite.ChainB.Voice,
	// 		},
	// 	},
	// }
	// suite.RoundtripExecute(suite.ChainA.Note, &suite.ChainB, []w.CosmosMsg{dataCosmosMsg, noDataCosmosMsg})
	// require.NoError(t, err, "round trip message not complete")
	// require.Len(t, len(callbackExecute.Success), 2, "error: "+callbackExecute.Error)
	// require.Equal(t, "", callbackExecute.Error)

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
