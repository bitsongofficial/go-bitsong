package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/candymachine/keeper"
	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// TODO: test UpdateCandyMachine
// TODO: test CloseCandyMachine
// TODO: test MintNFT

func (suite *KeeperTestSuite) TestMsgServerCreateCandyMachine() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase        string
		sender          sdk.AccAddress
		collectionOwner sdk.AccAddress
		collectionFee   sdk.Coin
		collectionId    uint64
		expectPass      bool
	}{
		{
			"when collection is not available",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 0),
			0,
			false,
		},
		{
			"when collection is not owned by the sender",
			addr1,
			addr2,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			false,
		},
		{
			"successful candymachine creation when fee is positive",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 1000),
			1,
			true,
		},
		{
			"successful candymachine creation when fee is zero",
			addr1,
			addr1,
			sdk.NewInt64Coin("ubtsg", 0),
			1,
			true,
		},
	}

	for _, tc := range tests {
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.NewCoins(tc.collectionFee))
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, tc.sender, sdk.NewCoins(tc.collectionFee))
		suite.Require().NoError(err)

		if tc.collectionId > 0 {
			suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
				Id:              tc.collectionId,
				Symbol:          "PUNK",
				UpdateAuthority: tc.collectionOwner.String(),
			})
		}

		params := suite.app.CandyMachineKeeper.GetParamSet(suite.ctx)
		params.CandymachineCreationPrice = tc.collectionFee
		suite.app.CandyMachineKeeper.SetParamSet(suite.ctx, params)

		msgServer := keeper.NewMsgServerImpl(suite.app.CandyMachineKeeper)
		machine := types.CandyMachine{
			CollId:     tc.collectionId,
			Price:      0,
			Treasury:   addr1.String(),
			Denom:      "ubtsg",
			GoLiveDate: 1659870342,
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Mint,
				Value:   1000,
			},
			Minted:               0,
			Authority:            tc.sender.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		oldSenderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "ubtsg")
		_, err = msgServer.CreateCandyMachine(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateCandyMachine(
			tc.sender, machine,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check collection authority is upgraded
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(collection.UpdateAuthority, moduleAddr.String())

			// check candymachine object is created from params
			savedMachine, err := suite.app.CandyMachineKeeper.GetCandyMachineByCollId(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(machine, savedMachine)

			// check fee spent when it is positive
			newSenderBalance := suite.app.BankKeeper.GetBalance(suite.ctx, tc.sender, "ubtsg")
			suite.Require().Equal(newSenderBalance.Add(tc.collectionFee), oldSenderBalance)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerUpdateCandyMachine() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase         string
		sender           sdk.AccAddress
		machineAuthority sdk.AccAddress
		collectionId     uint64
		expectPass       bool
	}{
		{
			"when candy machine is not available case",
			addr1,
			addr1,
			0,
			false,
		},
		{
			"when sender is not the authority",
			addr1,
			addr2,
			1,
			false,
		},
		{
			"successful candymachine upgrade",
			addr1,
			addr1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.CandyMachineKeeper)
		machine := types.CandyMachine{
			CollId:     1,
			Price:      0,
			Treasury:   tc.machineAuthority.String(),
			Denom:      "ubtsg",
			GoLiveDate: 1659870342,
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Mint,
				Value:   1000,
			},
			Minted:               0,
			Authority:            tc.machineAuthority.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		if tc.collectionId > 0 {
			suite.app.CandyMachineKeeper.SetCandyMachine(suite.ctx, machine)
		}

		machine.MetadataBaseUrl = "https://punk.com/newmeatadata"
		_, err := msgServer.UpdateCandyMachine(sdk.WrapSDKContext(suite.ctx), types.NewMsgUpdateCandyMachine(
			tc.sender, machine,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check candy machine upgraded
			savedMachine, err := suite.app.CandyMachineKeeper.GetCandyMachineByCollId(suite.ctx, machine.CollId)
			suite.Require().NoError(err)
			suite.Require().Equal(machine, savedMachine)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestMsgServerCloseCandyMachine() {
	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tests := []struct {
		testCase         string
		sender           sdk.AccAddress
		machineAuthority sdk.AccAddress
		collectionId     uint64
		expectPass       bool
	}{
		{
			"when candy machine is not available case",
			addr1,
			addr1,
			0,
			false,
		},
		{
			"when sender is not the authority",
			addr1,
			addr2,
			1,
			false,
		},
		{
			"successful candymachine close",
			addr1,
			addr1,
			1,
			true,
		},
	}

	for _, tc := range tests {
		msgServer := keeper.NewMsgServerImpl(suite.app.CandyMachineKeeper)
		machine := types.CandyMachine{
			CollId:     1,
			Price:      0,
			Treasury:   tc.machineAuthority.String(),
			Denom:      "ubtsg",
			GoLiveDate: 1659870342,
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Mint,
				Value:   1000,
			},
			Minted:               0,
			Authority:            tc.machineAuthority.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              true,
			SellerFeeBasisPoints: 100,
			Creators:             []nfttypes.Creator(nil),
		}

		if tc.collectionId > 0 {
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
			suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
				Id:              tc.collectionId,
				Symbol:          "PUNK",
				UpdateAuthority: moduleAddr.String(),
			})
			suite.app.CandyMachineKeeper.SetCandyMachine(suite.ctx, machine)
		}

		machine.MetadataBaseUrl = "https://punk.com/newmeatadata"
		_, err := msgServer.CloseCandyMachine(sdk.WrapSDKContext(suite.ctx), types.NewMsgCloseCandyMachine(
			tc.sender, tc.collectionId,
		))
		if tc.expectPass {
			suite.Require().NoError(err)

			// check candy machine deleted
			_, err := suite.app.CandyMachineKeeper.GetCandyMachineByCollId(suite.ctx, machine.CollId)
			suite.Require().Error(err)

			// check collection ownership transferred
			collection, err := suite.app.NFTKeeper.GetCollectionById(suite.ctx, tc.collectionId)
			suite.Require().NoError(err)
			suite.Require().Equal(collection.UpdateAuthority, tc.sender.String())
		} else {
			suite.Require().Error(err)
		}
	}
}
