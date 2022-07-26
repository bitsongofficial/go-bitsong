package cli_test

import (
	"encoding/json"
	"fmt"
	"github.com/bitsongofficial/go-bitsong/app"
	"github.com/bitsongofficial/go-bitsong/app/params"
	"github.com/bitsongofficial/go-bitsong/x/merkledrop/client/cli"
	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.cfg = app.DefaultConfig()

	s.network = network.New(s.T(), s.cfg)

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

func accountJson(t *testing.T, addr sdk.AccAddress) (string, map[string]cli.ClaimInfo, error) {
	list := fmt.Sprintf(`
		{
		  "%s": "100000",
		  "bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw": "200000",
		  "bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw": "300000"
		}
		`, addr.String(),
	)

	jsonFile := testutil.WriteToNewTempFile(t, list)

	var stringList map[string]string
	if err := json.Unmarshal([]byte(list), &stringList); err != nil {
		return "", map[string]cli.ClaimInfo{}, fmt.Errorf("Could not unmarshal json: %v", err)
	}

	accMap, err := cli.AccountsFromMap(stringList)
	if err != nil {
		return "", map[string]cli.ClaimInfo{}, fmt.Errorf("Could not get accounts from map")
	}

	_, claimInfo, _, err := cli.CreateDistributionList(accMap)
	if err != nil {
		return "", map[string]cli.ClaimInfo{}, fmt.Errorf("Could not create distribution list: %v", err)
	}

	return jsonFile.Name(), claimInfo, nil
}

func addCommonFlags(args []string, from sdk.AccAddress) []string {
	args = append(args, []string{
		fmt.Sprintf("--%s=%s", sdkflags.FlagFrom, from.String()),
		fmt.Sprintf("--%s=true", sdkflags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", sdkflags.FlagBroadcastMode, sdkflags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", sdkflags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
	}...)

	return args
}

func (s *IntegrationTestSuite) TestGetCmdQueryMerkledrop() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	//merkleRoot := "3452cae72dab475d017c1c46d289f9dc458a9fccf79add3e49347f2fc984e463"
	startHeight := 1
	endHeight := 40

	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("1000%s", s.cfg.BondDenom))
	s.Require().NoError(err)

	//------test GetCmdCreate()-------------
	cmd := cli.GetCmdCreate()
	fileJson, claimInfo, err := accountJson(s.T(), val.Address)
	outJson := "out.json"
	s.Require().NoError(err)

	args := addCommonFlags(
		[]string{
			fileJson,
			outJson,
			fmt.Sprintf("--%s=%s", cli.FlagDenom, coin.Denom),
			fmt.Sprintf("--%s=%d", cli.FlagStartHeight, startHeight),
			fmt.Sprintf("--%s=%d", cli.FlagEndHeight, endHeight),
		}, val.Address,
	)

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	os.Remove(outJson)

	//------test GetCmdQueryMerkledrop()-------------
	cmd = cli.GetCmdQueryMerkledrop()
	args = []string{
		"1",
	}

	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	s.Require().NoError(err, out.String())

	//------test GetCmdClaim()-------------
	s.network.WaitForHeight(2)

	cmd = cli.GetCmdClaim()
	args = addCommonFlags(
		[]string{
			"1",
			fmt.Sprintf("--%s=%s", cli.FlagProofs, strings.Join(claimInfo[val.Address.String()].Proof[:], ",")),
			fmt.Sprintf("--%s=%s", cli.FlagAmount, claimInfo[val.Address.String()].Amount),
			fmt.Sprintf("--%s=%d", cli.FlagIndex, claimInfo[val.Address.String()].Index),
		}, val.Address,
	)

	respType = proto.Message(&sdk.TxResponse{})
	expectedCode = uint32(0)
	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
}

func TestSimpleProof(t *testing.T) {
	leafs := [][]byte{
		[]byte("0bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj21000000"),
		[]byte("1bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw2000000"),
		[]byte("2bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw3000000"),
	}

	tree := cli.NewTree(leafs...)
	merkleRootStr := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "3452cae72dab475d017c1c46d289f9dc458a9fccf79add3e49347f2fc984e463", merkleRootStr)
}

func TestCreateProof(t *testing.T) {
	params.SetAddressPrefixes()

	accounts := map[string]string{
		"bitsong1vgpsha4f8grmsqr6krfdxwpcf3x20h0q3ztaj2": "1000000",
		"bitsong1zm6wlhr622yr9d7hh4t70acdfg6c32kcv34duw": "2000000",
		"bitsong1nzxmsks45e55d5edj4mcd08u8dycaxq5eplakw": "3000000",
	}

	accMap, err := cli.AccountsFromMap(accounts)
	assert.NoError(t, err)

	tree, _, _, err := cli.CreateDistributionList(accMap)
	assert.NoError(t, err)

	merkleRoot := fmt.Sprintf("%x", tree.Root())
	assert.Equal(t, "3452cae72dab475d017c1c46d289f9dc458a9fccf79add3e49347f2fc984e463", merkleRoot)
}
