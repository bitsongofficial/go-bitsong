package keeper_test

import (
	_ "embed"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

// Test register cadence contract .
func (s *IntegrationTestSuite) TestRegisterCadenceContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	// Store code
	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		isJailed bool
		success  bool
	}{
		{
			desc:     "Success - Register Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Register Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Register Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Contract Already Jailed",
			sender:   addr.String(),
			contract: contractAddress,
			isJailed: true,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Set params
			params := types.DefaultParams()
			err := s.App.CadenceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Jail contract if needed
			if tc.isJailed {
				s.RegisterCadenceContract(tc.sender, tc.contract)
				err := s.App.CadenceKeeper.SetJailStatus(s.Ctx, tc.contract, true)
				s.Require().NoError(err)
			}

			// Try to register contract
			res, err := s.cadenceMsgServer.RegisterCadenceContract(s.Ctx, &types.MsgRegisterCadenceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgRegisterCadenceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}

// Test standard unregistration of cadence contract s.
func (s *IntegrationTestSuite) TestUnregisterCadenceContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		success  bool
	}{
		{
			desc:     "Success - Unregister Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Unregister Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Unregister Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.RegisterCadenceContract(addr.String(), contractAddress)
			s.RegisterCadenceContract(addr2.String(), contractAddressWithAdmin)

			// Set params
			params := types.DefaultParams()
			err := s.App.CadenceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.cadenceMsgServer.UnregisterCadenceContract(s.Ctx, &types.MsgUnregisterCadenceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnregisterCadenceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}

// Test duplicate register/unregister cadence contract s.
func (s *IntegrationTestSuite) TestDuplicateRegistrationChecks() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")

	// Test double register, first succeed, second fail
	_, err := s.cadenceMsgServer.RegisterCadenceContract(s.Ctx, &types.MsgRegisterCadenceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.cadenceMsgServer.RegisterCadenceContract(s.Ctx, &types.MsgRegisterCadenceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)

	// Test double unregister, first succeed, second fail
	_, err = s.cadenceMsgServer.UnregisterCadenceContract(s.Ctx, &types.MsgUnregisterCadenceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.cadenceMsgServer.UnregisterCadenceContract(s.Ctx, &types.MsgUnregisterCadenceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)
}

// Test unjailing cadence contract s.
func (s *IntegrationTestSuite) TestUnjailCadenceContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		unjail   bool
		success  bool
	}{
		{
			desc:     "Success - Unjail Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Unjail Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Unjail Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Contract Not Jailed",
			sender:   addr.String(),
			contract: contractAddress,
			unjail:   true,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.RegisterCadenceContract(addr.String(), contractAddress)
			s.JailCadenceContract(contractAddress)
			s.RegisterCadenceContract(addr2.String(), contractAddressWithAdmin)
			s.JailCadenceContract(contractAddressWithAdmin)

			// Unjail contract if needed
			if tc.unjail {
				s.UnjailCadenceContract(addr.String(), contractAddress)
				s.UnjailCadenceContract(addr2.String(), contractAddressWithAdmin)
			}

			// Set params
			params := types.DefaultParams()
			err := s.App.CadenceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.cadenceMsgServer.UnjailCadenceContract(s.Ctx, &types.MsgUnjailCadenceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnjailCadenceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadenceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}
