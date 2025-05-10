package cadence_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitsongofficial/go-bitsong/x/cadence"
	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(CadenceModuleSuite))
}

func (suite *CadenceModuleSuite) SetupTest() {
	suite.Setup()
}

func (suite *CadenceModuleSuite) TestClockInitGenesis() {
	testCases := []struct {
		name    string
		genesis types.GenesisState
		success bool
	}{
		{
			"Success - Default Genesis",
			*cadence.DefaultGenesisState(),
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
					cadence.InitGenesis(suite.Ctx, suite.App.CadenceKeeper, tc.genesis)
				})

				params := suite.App.CadenceKeeper.GetParams(suite.Ctx)
				suite.Require().Equal(tc.genesis.Params, params)
			} else {
				suite.Require().Panics(func() {
					cadence.InitGenesis(suite.Ctx, suite.App.CadenceKeeper, tc.genesis)
				})
			}
		})
	}
}
