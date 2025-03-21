package cadance_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitsongofficial/go-bitsong/x/cadance"
	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(CadanceModuleSuite))
}

func (suite *CadanceModuleSuite) SetupTest() {
	suite.Setup()
}

func (suite *CadanceModuleSuite) TestClockInitGenesis() {
	testCases := []struct {
		name    string
		genesis types.GenesisState
		success bool
	}{
		{
			"Success - Default Genesis",
			*cadance.DefaultGenesisState(),
			true,
		},
		{
			"Success - Custom Genesis",
			types.GenesisState{
				Params: types.Params{
					ContractGasLimit: 500_000,
				},
			},
			true,
		},
		{
			"Fail - Invalid Gas Amount",
			types.GenesisState{
				Params: types.Params{
					ContractGasLimit: 1,
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			if tc.success {
				suite.Require().NotPanics(func() {
					cadance.InitGenesis(suite.Ctx, suite.App.CadanceKeeper, tc.genesis)
				})

				params := suite.App.CadanceKeeper.GetParams(suite.Ctx)
				suite.Require().Equal(tc.genesis.Params, params)
			} else {
				suite.Require().Panics(func() {
					cadance.InitGenesis(suite.Ctx, suite.App.CadanceKeeper, tc.genesis)
				})
			}
		})
	}
}
