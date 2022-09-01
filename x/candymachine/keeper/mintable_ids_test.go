package keeper_test

import (
	"time"
)

func (suite *KeeperTestSuite) TestMintableMetadataIdsGetSetDelete() {
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available metadataIds
	metadataIds := suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)

	// set mintable metadata ids and check
	suite.app.CandyMachineKeeper.SetMintableMetadataIds(suite.ctx, 1, []uint64{1, 2, 3, 4})
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 4)

	// delete mintable metadata ids and check
	suite.app.CandyMachineKeeper.DeleteMintableMetadataIds(suite.ctx, 1)
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)
}

func (suite *KeeperTestSuite) TestTakeOutFirstMintableMetadataId() {
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available metadataIds
	metadataIds := suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)

	// set mintable metadata ids and check
	suite.app.CandyMachineKeeper.SetMintableMetadataIds(suite.ctx, 1, []uint64{1, 2, 3, 4})
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 4)

	// take out first mintable metadata id and check
	metadataId := suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().Equal(metadataId, uint64(1))
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 3)

	// take out 3 more
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().Equal(metadataId, uint64(2))
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().Equal(metadataId, uint64(3))
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().Equal(metadataId, uint64(4))
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)
}

func (suite *KeeperTestSuite) TestTakeOutRandomMintableMetadataId() {
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available metadataIds
	metadataIds := suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)

	// set mintable metadata ids and check
	suite.app.CandyMachineKeeper.SetMintableMetadataIds(suite.ctx, 1, []uint64{1, 2, 3, 4})
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 4)

	// take out random mintable metadata id and check
	metadataId := suite.app.CandyMachineKeeper.TakeOutRandomMintableMetadataId(suite.ctx, 1, 4)
	suite.Require().GreaterOrEqual(metadataId, uint64(1))
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 3)

	// take out 3 more
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().GreaterOrEqual(metadataId, uint64(1))
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().GreaterOrEqual(metadataId, uint64(1))
	metadataId = suite.app.CandyMachineKeeper.TakeOutFirstMintableMetadataId(suite.ctx, 1)
	suite.Require().GreaterOrEqual(metadataId, uint64(1))
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)
}

func (suite *KeeperTestSuite) TestShuffleMintableMetadataIds() {
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	// get not available metadataIds
	metadataIds := suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 0)

	// set mintable metadata ids and check
	suite.app.CandyMachineKeeper.SetMintableMetadataIds(suite.ctx, 1, []uint64{1, 2, 3, 4})
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 4)

	// shuffle mintable metadata ids and check
	suite.app.CandyMachineKeeper.ShuffleMintableMetadataIds(suite.ctx, 1)
	metadataIds = suite.app.CandyMachineKeeper.GetMintableMetadataIds(suite.ctx, 1)
	suite.Require().Len(metadataIds, 4)
}
