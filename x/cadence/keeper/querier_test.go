package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

// Query Clock Params
func (s *IntegrationTestSuite) TestQueryClockParams() {
	for _, tc := range []struct {
		desc   string
		params types.Params
	}{
		{
			desc:   "On default",
			params: types.DefaultParams(),
		},
		{
			desc: "On 500_000",
			params: types.Params{
				ContractGasLimit: 500_000,
			},
		},
		{
			desc: "On 1_000_000",
			params: types.Params{
				ContractGasLimit: 1_000_000,
			},
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Set params
			err := s.App.CadenceKeeper.SetParams(s.Ctx, tc.params)
			s.Require().NoError(err)

			// Query params
			resp, err := s.queryClient.Params(s.Ctx, &types.QueryParamsRequest{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tc.params, *resp.Params)
		})
	}
}

// Query Clock Contracts
func (s *IntegrationTestSuite) TestQueryCadenceContracts() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()

	for _, tc := range []struct {
		desc      string
		contracts []string
	}{
		{
			desc:      "On empty",
			contracts: []string(nil),
		},
		{
			desc: "On Single",
			contracts: []string{
				s.InstantiateContract(addr.String(), ""),
			},
		},
		{
			desc: "On Multiple",
			contracts: []string{
				s.InstantiateContract(addr.String(), ""),
				s.InstantiateContract(addr.String(), ""),
				s.InstantiateContract(addr.String(), ""),
			},
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Loop through contracts & register
			for _, contract := range tc.contracts {
				s.RegisterCadenceContract(addr.String(), contract)
			}

			// Contracts check
			resp, err := s.queryClient.CadenceContracts(s.Ctx, &types.QueryCadenceContracts{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			for _, contract := range resp.CadenceContracts {
				s.Require().Contains(tc.contracts, contract.ContractAddress)
				s.Require().False(contract.IsJailed)
			}

			// Remove all contracts
			for _, contract := range tc.contracts {
				s.App.CadenceKeeper.RemoveContract(s.Ctx, contract)
			}
		})
	}
}

// Query Jailed Clock Contracts
func (s *IntegrationTestSuite) TestQueryJailedCadenceContracts() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()

	for _, tc := range []struct {
		desc      string
		contracts []string
	}{
		{
			desc:      "On empty",
			contracts: []string(nil),
		},
		{
			desc: "On Single",
			contracts: []string{
				s.InstantiateContract(addr.String(), ""),
			},
		},
		{
			desc: "On Multiple",
			contracts: []string{
				s.InstantiateContract(addr.String(), ""),
				s.InstantiateContract(addr.String(), ""),
				s.InstantiateContract(addr.String(), ""),
			},
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Loop through contracts & register
			for _, contract := range tc.contracts {
				s.RegisterCadenceContract(addr.String(), contract)
				s.JailCadenceContract(contract)
			}

			// Contracts check
			resp, err := s.queryClient.CadenceContracts(s.Ctx, &types.QueryCadenceContracts{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			for _, contract := range resp.CadenceContracts {
				s.Require().Contains(tc.contracts, contract.ContractAddress)
				s.Require().True(contract.IsJailed)
			}

			// Remove all contracts
			for _, contract := range tc.contracts {
				s.App.CadenceKeeper.RemoveContract(s.Ctx, contract)
			}
		})
	}
}

// Query Clock Contract
func (s *IntegrationTestSuite) TestQueryCadenceContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))
	_, _, invalidAddr := testdata.KeyTestPubAddr()

	s.StoreCode()

	unjailedContract := s.InstantiateContract(addr.String(), "")
	_ = s.App.CadenceKeeper.SetCadenceContract(s.Ctx, types.CadenceContract{
		ContractAddress: unjailedContract,
		IsJailed:        false,
	})

	jailedContract := s.InstantiateContract(addr.String(), "")
	_ = s.App.CadenceKeeper.SetCadenceContract(s.Ctx, types.CadenceContract{
		ContractAddress: jailedContract,
		IsJailed:        true,
	})

	for _, tc := range []struct {
		desc     string
		contract string
		isJailed bool
		success  bool
	}{
		{
			desc:     "On Unjailed",
			contract: unjailedContract,
			isJailed: false,
			success:  true,
		},
		{
			desc:     "On Jailed",
			contract: jailedContract,
			isJailed: true,
			success:  true,
		},
		{
			desc:     "Invalid Contract - Unjailed",
			contract: invalidAddr.String(),
			isJailed: false,
			success:  false,
		},
		{
			desc:     "Invalid Contract - Jailed",
			contract: invalidAddr.String(),
			isJailed: true,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Query contract
			resp, err := s.queryClient.CadenceContract(s.Ctx, &types.QueryCadenceContract{
				ContractAddress: tc.contract,
			})

			// Validate responses
			if tc.success {
				s.Require().NoError(err)
				s.Require().Equal(resp.CadenceContract.ContractAddress, tc.contract)
				s.Require().Equal(resp.CadenceContract.IsJailed, tc.isJailed)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
