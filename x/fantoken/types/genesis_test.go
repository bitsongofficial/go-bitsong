package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidateGenesis(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: DefaultGenesisState(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "empty authority",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "issue fee 0",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "no fantokens",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: nil,
			},
			valid: true,
		},
		{
			desc: "no metadata",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData:  Metadata{},
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid symbol",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "test token",
							Symbol: "",
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "empty name",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "",
							Symbol: "fttest",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "empty uri",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: math.NewInt(1),
						MetaData: Metadata{
							Name:   "test token",
							Symbol: "fttest",
							URI:    "",
						},
					},
				},
			},
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
