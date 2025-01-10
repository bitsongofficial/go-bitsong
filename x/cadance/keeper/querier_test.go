package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
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
			err := s.app.AppKeepers.CadanceKeeper.SetParams(s.ctx, tc.params)
			s.Require().NoError(err)

			// Query params
			goCtx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.queryClient.Params(goCtx, &types.QueryParamsRequest{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().Equal(tc.params, *resp.Params)
		})
	}
}

// Query Clock Contracts
func (s *IntegrationTestSuite) TestQueryCadanceContracts() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

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
				s.RegisterCadanceContract(addr.String(), contract)
			}

			// Contracts check
			goCtx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.queryClient.CadanceContracts(goCtx, &types.QueryCadanceContracts{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			for _, contract := range resp.CadanceContracts {
				s.Require().Contains(tc.contracts, contract.ContractAddress)
				s.Require().False(contract.IsJailed)
			}

			// Remove all contracts
			for _, contract := range tc.contracts {
				s.app.AppKeepers.CadanceKeeper.RemoveContract(s.ctx, contract)
			}
		})
	}
}

// Query Jailed Clock Contracts
func (s *IntegrationTestSuite) TestQueryJailedCadanceContracts() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

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
				s.RegisterCadanceContract(addr.String(), contract)
				s.JailCadanceContract(contract)
			}

			// Contracts check
			goCtx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.queryClient.CadanceContracts(goCtx, &types.QueryCadanceContracts{})

			// Response check
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			for _, contract := range resp.CadanceContracts {
				s.Require().Contains(tc.contracts, contract.ContractAddress)
				s.Require().True(contract.IsJailed)
			}

			// Remove all contracts
			for _, contract := range tc.contracts {
				s.app.AppKeepers.CadanceKeeper.RemoveContract(s.ctx, contract)
			}
		})
	}
}

// Query Clock Contract
func (s *IntegrationTestSuite) TestQueryCadanceContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))
	_, _, invalidAddr := testdata.KeyTestPubAddr()

	s.StoreCode()

	unjailedContract := s.InstantiateContract(addr.String(), "")
	_ = s.app.AppKeepers.CadanceKeeper.SetCadanceContract(s.ctx, types.CadanceContract{
		ContractAddress: unjailedContract,
		IsJailed:        false,
	})

	jailedContract := s.InstantiateContract(addr.String(), "")
	_ = s.app.AppKeepers.CadanceKeeper.SetCadanceContract(s.ctx, types.CadanceContract{
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
			resp, err := s.queryClient.CadanceContract(s.ctx, &types.QueryCadanceContract{
				ContractAddress: tc.contract,
			})

			// Validate responses
			if tc.success {
				s.Require().NoError(err)
				s.Require().Equal(resp.CadanceContract.ContractAddress, tc.contract)
				s.Require().Equal(resp.CadanceContract.IsJailed, tc.isJailed)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
