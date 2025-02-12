package fantoken_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	apptesting.KeeperTestHelper
	govHandler govv1beta1.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.Setup()
	// suite.ctx = suite.app.BaseApp.NewContext(false)
	suite.govHandler = params.NewParamChangeProposalHandler(suite.App.AppKeepers.ParamsKeeper)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func testProposal(changes ...proposal.ParamChange) *proposal.ParameterChangeProposal {
	return proposal.NewParameterChangeProposal("title", "description", changes)
}

func (suite *HandlerTestSuite) TestParamChangeProposal() {
	tp := testProposal(
		proposal.NewParamChange(
			fantokentypes.ModuleName,
			string(fantokentypes.KeyMintFee),
			"{\"denom\":\"utsg\",\"amount\":\"0\"}",
		),
	)

	fmt.Println(tp.String())

	err := suite.govHandler(suite.Ctx, tp)
	suite.Require().NoError(err)
}

func (suite *HandlerTestSuite) TestProposalHandlerPassed() {

	params := suite.App.AppKeepers.FanTokenKeeper.GetParamSet(suite.Ctx)
	require.Equal(suite.T(), params, fantokentypes.DefaultParams())

	newIssueFee := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1))
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(2))
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(3))

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
	)

	h := fantoken.NewProposalHandler(suite.App.AppKeepers.FanTokenKeeper)
	require.NoError(suite.T(), h(suite.Ctx, proposal))

	params = suite.App.AppKeepers.FanTokenKeeper.GetParamSet(suite.Ctx)
	require.Equal(suite.T(), newIssueFee, params.IssueFee)
	require.Equal(suite.T(), newMintFee, params.MintFee)
	require.Equal(suite.T(), newBurnFee, params.BurnFee)
}

func (suite *HandlerTestSuite) TestProposalHandlerFailed() {

	params := suite.App.AppKeepers.FanTokenKeeper.GetParamSet(suite.Ctx)
	require.Equal(suite.T(), params, fantokentypes.DefaultParams())

	newIssueFee := sdk.Coin{
		Denom:  sdk.DefaultBondDenom,
		Amount: math.NewInt(-1),
	}
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt())
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, math.ZeroInt())

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
	)

	h := fantoken.NewProposalHandler(suite.App.AppKeepers.FanTokenKeeper)
	require.Error(suite.T(), h(suite.Ctx, proposal))
}
