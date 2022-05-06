package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{
			"Minimum value",
			Params{
				IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
			},
			true,
		}, {
			"Maximum value",
			Params{
				IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(math.MaxInt64)),
			},
			true,
		}, {
			"IssueFee is negative",
			Params{
				IssueFee: sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(-1)},
			},
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, ValidateParams(tc.Params), "test: %v", tc.testCase)
		}
	}
}
