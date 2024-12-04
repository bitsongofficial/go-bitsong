// source : https://github.com/DA0-DA0/polytone/blob/main/tests/strangelove/incompatible_handshake_test.go
package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests that a note may only ever connect to a voice, and a voice
// only to a note.
func TestPolytoneOnBitsong(t *testing.T) {
	suite := NewPolytoneSuite(t)

	// note <-> note not allowed.
	_, _, err := suite.CreateChannel(
		suite.ChainA.Note,
		suite.ChainB.Note,
		&suite.ChainA,
		&suite.ChainB, suite.PathAB,
	)
	require.ErrorContains(t, err, "no new channels created", "note <-/-> note")

	// channels := suite.QueryChannelsInState(&suite.ChainB, CHANNEL_STATE_TRY)
	// require.Len(t, channels, 1, "try note stops in first step")
	// channels = suite.QueryChannelsInState(&suite.ChainB, CHANNEL_STATE_INIT)
	// require.Len(t, channels, 1, "init note doesn't advance")

	// voice <-> voice not allowed
	_, _, err = suite.CreateChannel(
		suite.ChainA.Voice,
		suite.ChainB.Voice,
		&suite.ChainA,
		&suite.ChainB,
		suite.PathAB,
	)
	require.ErrorContains(t, err, "no new channels created", "voice <-/-> voice")

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
