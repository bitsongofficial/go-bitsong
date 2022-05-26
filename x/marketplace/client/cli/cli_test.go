package cli_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"

	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/nft/client/testutil"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	_, err = testutil.CreateNFT(clientCtx, val.Address.String(), s.cfg.BondDenom)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// TODO: add test for GetCmdQueryAuctions
// TODO: add test for GetCmdQueryAuction
// TODO: add test for GetCmdQueryBidsByAuction
// TODO: add test for GetCmdQueryBidsByBidder
// TODO: add test for GetCmdQueryBidderMetadata

// TODO: add test for GetCmdCreateAuction
// TODO: add test for GetCmdSetAuctionAuthority
// TODO: add test for GetCmdStartAuction
// TODO: add test for GetCmdEndAuction
// TODO: add test for GetCmdPlaceBid
// TODO: add test for GetCmdCancelBid
// TODO: add test for GetCmdClaimBid
