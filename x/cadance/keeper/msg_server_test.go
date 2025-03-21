package keeper_test

import (
	_ "embed"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadance/types"
)

// Test register cadance contract .
func (s *IntegrationTestSuite) TestRegisterCadanceContract() {
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
			err := s.App.CadanceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Jail contract if needed
			if tc.isJailed {
				s.RegisterCadanceContract(tc.sender, tc.contract)
				err := s.App.CadanceKeeper.SetJailStatus(s.Ctx, tc.contract, true)
				s.Require().NoError(err)
			}

			// Try to register contract
			res, err := s.cadanceMsgServer.RegisterCadanceContract(s.Ctx, &types.MsgRegisterCadanceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgRegisterCadanceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}

// Test standard unregistration of cadance contract s.
func (s *IntegrationTestSuite) TestUnregisterCadanceContract() {
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
			s.RegisterCadanceContract(addr.String(), contractAddress)
			s.RegisterCadanceContract(addr2.String(), contractAddressWithAdmin)

			// Set params
			params := types.DefaultParams()
			err := s.App.CadanceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.cadanceMsgServer.UnregisterCadanceContract(s.Ctx, &types.MsgUnregisterCadanceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnregisterCadanceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}

// Test duplicate register/unregister cadance contract s.
func (s *IntegrationTestSuite) TestDuplicateRegistrationChecks() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.Ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")

	// Test double register, first succeed, second fail
	_, err := s.cadanceMsgServer.RegisterCadanceContract(s.Ctx, &types.MsgRegisterCadanceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.cadanceMsgServer.RegisterCadanceContract(s.Ctx, &types.MsgRegisterCadanceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)

	// Test double unregister, first succeed, second fail
	_, err = s.cadanceMsgServer.UnregisterCadanceContract(s.Ctx, &types.MsgUnregisterCadanceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.cadanceMsgServer.UnregisterCadanceContract(s.Ctx, &types.MsgUnregisterCadanceContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)
}

// Test unjailing cadance contract s.
func (s *IntegrationTestSuite) TestUnjailCadanceContract() {
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
			s.RegisterCadanceContract(addr.String(), contractAddress)
			s.JailCadanceContract(contractAddress)
			s.RegisterCadanceContract(addr2.String(), contractAddressWithAdmin)
			s.JailCadanceContract(contractAddressWithAdmin)

			// Unjail contract if needed
			if tc.unjail {
				s.UnjailCadanceContract(addr.String(), contractAddress)
				s.UnjailCadanceContract(addr2.String(), contractAddressWithAdmin)
			}

			// Set params
			params := types.DefaultParams()
			err := s.App.CadanceKeeper.SetParams(s.Ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.cadanceMsgServer.UnjailCadanceContract(s.Ctx, &types.MsgUnjailCadanceContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnjailCadanceContractResponse{})
			}

			// Ensure contract is unregistered
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddress)
			s.App.CadanceKeeper.RemoveContract(s.Ctx, contractAddressWithAdmin)
		})
	}
}
