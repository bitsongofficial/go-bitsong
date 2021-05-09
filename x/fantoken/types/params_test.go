package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	defaultToken := GetNativeToken()
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{
			"Minimum value",
			Params{
				IssuePrice: sdk.NewCoin(defaultToken.Denom, sdk.ZeroInt()),
			},
			true,
		}, {
			"Maximum value",
			Params{
				IssuePrice: sdk.NewCoin(defaultToken.Denom, sdk.NewInt(math.MaxInt64)),
			},
			true,
		}, {
			"IssuePrice is negative",
			Params{
				IssuePrice: sdk.Coin{Denom: defaultToken.Denom, Amount: sdk.NewInt(-1)},
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
