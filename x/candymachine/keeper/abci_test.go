package keeper_test

import (
	"time"

	"github.com/bitsongofficial/go-bitsong/x/candymachine/types"
	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

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
		moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)
		suite.app.NFTKeeper.SetCollection(suite.ctx, nfttypes.Collection{
			Id:              machine.CollId,
			Symbol:          "PUNK",
			UpdateAuthority: moduleAddr.String(),
		})
		suite.app.CandyMachineKeeper.SetCandyMachine(suite.ctx, machine)
	}

	endingMachines := suite.app.CandyMachineKeeper.GetCandyMachinesToEndByTime(suite.ctx)
	suite.Require().Len(endingMachines, 1)

	allMachines := suite.app.CandyMachineKeeper.GetAllCandyMachines(suite.ctx)
	suite.Require().Len(allMachines, 4)
	suite.Require().Equal(machines, allMachines)

	suite.app.CandyMachineKeeper.EndBlocker(suite.ctx)

	// check one machine automatically closed
	allMachines = suite.app.CandyMachineKeeper.GetAllCandyMachines(suite.ctx)
	suite.Require().Len(allMachines, 3)
}
