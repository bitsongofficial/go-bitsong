package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
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
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
				BurnedCoins: nil,
			},
			valid: true,
		},
		{
			desc: "empty authority",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
				BurnedCoins: nil,
			},
			valid: true,
		},
		{
			desc: "issue fee 0",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "test fantoken",
							Symbol: "symbol",
							URI:    "ipfs://...",
						},
					},
				},
				BurnedCoins: nil,
			},
			valid: true,
		},
		{
			desc: "no fantokens",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens:   nil,
				BurnedCoins: nil,
			},
			valid: true,
		},
		{
			desc: "no metadata",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData:  Metadata{},
					},
				},
				BurnedCoins: nil,
			},
			valid: false,
		},
		{
			desc: "invalid symbol",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "test token",
							Symbol: "",
						},
					},
				},
				BurnedCoins: nil,
			},
			valid: false,
		},
		{
			desc: "empty name",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "",
							Symbol: "fttest",
						},
					},
				},
				BurnedCoins: nil,
			},
			valid: true,
		},
		{
			desc: "empty uri",
			genState: &GenesisState{
				Params: Params{
					IssueFee: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
					MintFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
					BurnFee:  sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()),
				},
				FanTokens: []FanToken{
					{
						Denom:     "fttest",
						MaxSupply: sdk.NewInt(1),
						MetaData: Metadata{
							Name:   "test token",
							Symbol: "fttest",
							URI:    "",
						},
					},
				},
				BurnedCoins: nil,
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
