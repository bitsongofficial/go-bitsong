package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// TODO: test CreateCandyMachine
// TODO: test UpdateCandyMachine
// TODO: test CloseCandyMachine
// TODO: test MintNFT

func (suite *KeeperTestSuite) TestCandyMachineGetSetDelete() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available candymachine
	_, err := suite.app.CandyMachineKeeper.GetCandyMachineByCollId(suite.ctx, 0)
	suite.Require().Error(err)

	// get all candymachines when not available
	allMachines := suite.app.CandyMachineKeeper.GetAllCandyMachines(suite.ctx)
	suite.Require().Len(allMachines, 0)

	// get ending queue candy machines when not available
	endingMachines := suite.app.CandyMachineKeeper.GetCandyMachinesToEndByTime(suite.ctx)
	suite.Require().Len(endingMachines, 0)

	// set candy machines
	machines := []types.CandyMachine{
		{
			CollId:     1,
			Price:      100,
			Treasury:   addr.String(),
			Denom:      "ubtsg",
			GoLiveDate: uint64(now.Unix()),
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Time,
				Value:   uint64(now.Unix()),
			},
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:     2,
			Price:      100,
			Treasury:   addr.String(),
			Denom:      "ubtsg",
			GoLiveDate: uint64(now.Unix()),
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Time,
				Value:   uint64(now.Unix()) + 100,
			},
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:     3,
			Price:      100,
			Treasury:   addr.String(),
			Denom:      "ubtsg",
			GoLiveDate: uint64(now.Unix()),
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Time,
				Value:   uint64(now.Unix()) - 100,
			},
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
		{
			CollId:     4,
			Price:      100,
			Treasury:   addr.String(),
			Denom:      "ubtsg",
			GoLiveDate: uint64(now.Unix()),
			EndSettings: types.EndSettings{
				EndType: types.EndSettingType_Mint,
				Value:   10,
			},
			Minted:               0,
			Authority:            addr.String(),
			MetadataBaseUrl:      "https://punk.com/metadata",
			Mutable:              false,
			SellerFeeBasisPoints: 100,
		},
	}

	for _, machine := range machines {
		suite.app.CandyMachineKeeper.SetCandyMachine(suite.ctx, machine)
	}

	for _, machine := range machines {
		m, err := suite.app.CandyMachineKeeper.GetCandyMachineByCollId(suite.ctx, machine.CollId)
		suite.Require().NoError(err)
		suite.Require().Equal(machine, m)
	}

	allMachines = suite.app.CandyMachineKeeper.GetAllCandyMachines(suite.ctx)
	suite.Require().Len(allMachines, 4)
	suite.Require().Equal(machines, allMachines)

	endingMachines = suite.app.CandyMachineKeeper.GetCandyMachinesToEndByTime(suite.ctx)
	suite.Require().Len(endingMachines, 1)

	for _, machine := range machines {
		suite.app.CandyMachineKeeper.DeleteCandyMachine(suite.ctx, machine)
	}

	// get all candymachines after all deletion
	allMachines = suite.app.CandyMachineKeeper.GetAllCandyMachines(suite.ctx)
	suite.Require().Len(allMachines, 0)

	// get ending queue candy machines after all deletion
	endingMachines = suite.app.CandyMachineKeeper.GetCandyMachinesToEndByTime(suite.ctx)
	suite.Require().Len(endingMachines, 0)
}
