package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	tokencli "github.com/bitsongofficial/bitsong/x/fantoken/client/cli"
)

func IssueFanTokenExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdIssueFanToken(), args)
}

func UpdateFanTokenMintableExec(clientCtx client.Context, from string, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdUpdateFanTokenMintable(), args)
}

func MintFanTokenExec(clientCtx client.Context, from string, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdMintFanToken(), args)
}

func BurnFanTokenExec(clientCtx client.Context, from string, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdBurnFanToken(), args)
}

func TransferFanTokenOwnerExec(clientCtx client.Context, from string, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdTransferFanTokenOwner(), args)
}

func QueryFanTokenExec(clientCtx client.Context, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanToken(), args)
}

func QueryFanTokensExec(clientCtx client.Context, owner string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		owner,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFanTokens(), args)
}

func QueryParamsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryParams(), args)
}
