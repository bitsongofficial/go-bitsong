package app

import (
	"github.com/bitsongofficial/go-bitsong/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	defaultState := stakingGenesisState()
	require.Equal(t, types.BondDenom, defaultState.Params.BondDenom)
}
