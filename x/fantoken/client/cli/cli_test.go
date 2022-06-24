package cli_test

import (
	"context"
	"fmt"
	simapp "github.com/bitsongofficial/go-bitsong/app"
	tokentypes "github.com/bitsongofficial/go-bitsong/x/fantoken/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/cli"
	"testing"

	tokencli "github.com/bitsongofficial/go-bitsong/x/fantoken/client/cli"
)

var (
	name      = "Bitcoin"
	symbol    = "btc"
	uri       = "ipfs://"
	maxSupply = sdk.NewInt(200000000)
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.DefaultConfig()
	cfg.NumValidators = 2

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func issueCmd(s *IntegrationTestSuite, ctx client.Context, name, symbol string, maxSupply sdk.Int, uri string, from sdk.AccAddress) string {
	//------test GetCmdIssue()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%s", tokencli.FlagURI, uri),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdIssue(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	denom := string(txResp.Events[12].Attributes[0].Value)
	return denom[1 : len(denom)-1]
}

func mintCmd(s *IntegrationTestSuite, ctx client.Context, coin string, rcpt, from sdk.AccAddress) {
	args := []string{
		coin,
		fmt.Sprintf("--%s=%s", tokencli.FlagRecipient, rcpt),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdMint(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func burnCmd(s *IntegrationTestSuite, ctx client.Context, coin string, from sdk.AccAddress) {
	args := []string{
		coin,

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdBurn(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func setAuthorityCmd(s *IntegrationTestSuite, ctx client.Context, denom, new, old string) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", tokencli.FlagNewAuthority, new),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, old),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdSetAuthority(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func setMinterCmd(s *IntegrationTestSuite, ctx client.Context, denom, new, old string) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", tokencli.FlagNewMinter, new),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, old),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdSetMinter(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func disableMintCmd(s *IntegrationTestSuite, ctx client.Context, denom, from string) {
	args := []string{
		denom,

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdDisableMint(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func setUriCmd(s *IntegrationTestSuite, ctx client.Context, denom, uri, from string) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", tokencli.FlagURI, uri),

		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := clitestutil.ExecTestCLICmd(ctx, tokencli.GetCmdSetUri(), args)
	s.Require().NoError(err)
	s.Require().NoError(ctx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())

	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func queryBalance(s *IntegrationTestSuite, ctx client.Context, denom string, addr sdk.AccAddress) *sdk.Coin {
	bankClient := banktypes.NewQueryClient(ctx)

	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{
			Address: addr.String(),
			Denom:   denom,
		},
	)
	s.Require().NoError(err)

	return bankRes.Balance
}

func (s *IntegrationTestSuite) TestCmdIssue() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)
}

func (s *IntegrationTestSuite) TestCmdMint() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	// issue a new fantoken
	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	// mint 10 tokens
	coin := fmt.Sprintf("%d%s", 10, denom)
	mintCmd(s, clientCtx, coin, from, from)

	// query balance
	balance := queryBalance(s, clientCtx, denom, from)
	s.Require().Equal(coin, balance.String())
}

func (s *IntegrationTestSuite) TestCmdBurn() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	// issue a new fantoken
	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	// mint 10 tokens
	coin := fmt.Sprintf("%d%s", 10, denom)
	mintCmd(s, clientCtx, coin, from, from)

	// burn 6 tokens
	coin = fmt.Sprintf("%d%s", 6, denom)
	burnCmd(s, clientCtx, coin, from)

	// query balance
	expBalance := fmt.Sprintf("%d%s", 4, denom)
	balance := queryBalance(s, clientCtx, denom, from)
	s.Require().Equal(expBalance, balance.String())
}

func (s *IntegrationTestSuite) TestCmdSetAuthority() {
	val := s.network.Validators[0]
	val2 := s.network.Validators[1]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	setAuthorityCmd(s, clientCtx, denom, val2.Address.String(), from.String())

	var response tokentypes.QueryFanTokenResponse

	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.MetaData.Authority, val2.Address.String())
}

func (s *IntegrationTestSuite) TestCmdSetMinter() {
	val := s.network.Validators[0]
	val2 := s.network.Validators[1]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	setMinterCmd(s, clientCtx, denom, val2.Address.String(), from.String())

	var response tokentypes.QueryFanTokenResponse

	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.Minter, val2.Address.String())
}

func (s *IntegrationTestSuite) TestCmdSetUri() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	newUri := "ipfs://newuri"
	setUriCmd(s, clientCtx, denom, newUri, from.String())

	var response tokentypes.QueryFanTokenResponse

	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.MetaData.URI, newUri)
}

func (s *IntegrationTestSuite) TestCmdDisableMint() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	var response tokentypes.QueryFanTokenResponse

	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.Minter, from.String())

	disableMintCmd(s, clientCtx, denom, from.String())

	resp, err = clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.Minter, "")
}

func (s *IntegrationTestSuite) TestCmdQueryFanToken() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	var response tokentypes.QueryFanTokenResponse

	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().Equal(response.Fantoken.Denom, denom)
}

func (s *IntegrationTestSuite) TestCmdQueryFanTokens() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	var response tokentypes.QueryFanTokensResponse

	args := []string{
		from.String(),
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanTokens(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &response))
	s.Require().True(len(response.Fantokens) >= 1)
}

func (s *IntegrationTestSuite) TestCmdQueryParams() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	var params tokentypes.Params

	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryParams(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &params))
	s.Require().Equal(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)), params.IssueFee)
	s.Require().Equal(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()), params.MintFee)
	s.Require().Equal(sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt()), params.BurnFee)
}

func (s *IntegrationTestSuite) TestCmdQueryTotalBurn() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	from := val.Address

	denom := issueCmd(s, clientCtx, name, symbol, maxSupply, uri, from)

	var totalBurn tokentypes.QueryTotalBurnResponse

	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}

	resp, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryTotalBurn(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &totalBurn))
	s.Require().GreaterOrEqual(len(totalBurn.BurnedCoins), 0)

	// mint
	coin := fmt.Sprintf("%d%s", 10, denom)
	mintCmd(s, clientCtx, coin, from, from)

	// burn
	coin = fmt.Sprintf("%d%s", 10, denom)
	burnCmd(s, clientCtx, coin, from)

	// query again
	resp, err = clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryTotalBurn(), args)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp.Bytes(), &totalBurn))
	s.Require().GreaterOrEqual(len(totalBurn.BurnedCoins), 1)
	s.Require().Equal(sdk.NewInt(10), totalBurn.BurnedCoins[0].Amount)
	s.Require().Equal(denom, totalBurn.BurnedCoins[0].Denom)
}
