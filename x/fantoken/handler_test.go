package fantoken_test

import (
	"fmt"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/x/fantoken"
	fantokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *simapp.BitsongApp
	ctx        sdk.Context
	govHandler govtypes.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.govHandler = params.NewParamChangeProposalHandler(suite.app.ParamsKeeper)
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

	err := suite.govHandler(suite.ctx, tp)
	suite.Require().NoError(err)
}

func TestProposalHandlerPassed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	params := app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, params, fantokentypes.DefaultParams())

	newIssueFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2))
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
	)

	h := fantoken.NewProposalHandler(app.FanTokenKeeper)
	require.NoError(t, h(ctx, proposal))

	params = app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, newIssueFee, params.IssueFee)
	require.Equal(t, newMintFee, params.MintFee)
	require.Equal(t, newBurnFee, params.BurnFee)
}

func TestProposalHandlerFailed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	params := app.FanTokenKeeper.GetParamSet(ctx)
	require.Equal(t, params, fantokentypes.DefaultParams())

	newIssueFee := sdk.Coin{
		Denom:  sdk.DefaultBondDenom,
		Amount: sdk.NewInt(-1),
	}
	newMintFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
	newBurnFee := sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())

	proposal := fantokentypes.NewUpdateFeesProposal(
		"Test",
		"description",
		newIssueFee,
		newMintFee,
		newBurnFee,
	)

	h := fantoken.NewProposalHandler(app.FanTokenKeeper)
	require.Error(t, h(ctx, proposal))
}
