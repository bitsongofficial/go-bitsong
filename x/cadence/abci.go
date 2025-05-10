package cadence

import (
	"time"

	"cosmossdk.io/log"

	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitsongofficial/go-bitsong/x/cadence/keeper"
	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
)

var endBlockSudoMessage = []byte(types.EndBlockSudoMessage)

// EndBlocker executes on contracts at the end of the block.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)
	p := k.GetParams(ctx)

	// Get all contracts
	contracts, err := k.GetAllContracts(ctx)
	if err != nil {
		logger.Error("Failed to get contracts", "error", err)
		return
	}

	// Track errors
	errorExecs := make([]string, len(contracts))
	errorExists := false

	// Execute all contracts that are not jailed
	for idx, contract := range contracts {

		// Skip jailed contracts
		if contract.IsJailed {
			continue
		}

		// Get sdk.AccAddress from contract address
		contractAddr := sdk.MustAccAddressFromBech32(contract.ContractAddress)
		if handleError(ctx, k, logger, errorExecs, &errorExists, err, idx, contract.ContractAddress) {
			continue
		}

		// Create context with gas limit
		childCtx := ctx.WithGasMeter(storetypes.NewGasMeter(p.ContractGasLimit))

		// Execute contract
		ExecuteContract(k.GetContractKeeper(), childCtx, contractAddr, endBlockSudoMessage, &err)
		if handleError(ctx, k, logger, errorExecs, &errorExists, err, idx, contract.ContractAddress) {
			continue
		}
	}

	// Log errors if present
	if errorExists {
		logger.Error("Failed to execute contracts", "contracts", errorExecs)
	}
}

// Function to handle contract execution errors. Returns true if error is present, false otherwise.
func handleError(
	ctx sdk.Context,
	k keeper.Keeper,
	logger log.Logger,
	errorExecs []string,
	errorExists *bool,
	err error,
	idx int,
	contractAddress string,
) bool {
	// Check if error is present
	if err != nil {

		// Flag error
		*errorExists = true
		errorExecs[idx] = contractAddress

		// Attempt to jail contract, log error if present
		err := k.SetJailStatus(ctx, contractAddress, true)
		if err != nil {
			logger.Error("Failed to jail contract", "contract", contractAddress, "error", err)
		}
	}

	return err != nil
}

// Execute contract, recover from panic
func ExecuteContract(k wasmtypes.ContractOpsKeeper, childCtx sdk.Context, contractAddr sdk.AccAddress, msgBz []byte, err *error) {
	// Recover from panic, return error
	defer func() {
		if recoveryError := recover(); recoveryError != nil {
			// Determine error associated with panic
			if isOutofGas, msg := IsOutOfGasError(recoveryError); isOutofGas {
				*err = types.ErrOutOfGas.Wrapf("%s", msg)
			} else {
				*err = types.ErrContractExecutionPanic.Wrapf("%s", recoveryError)
			}
		}
	}()

	// Execute contract with sudo
	_, *err = k.Sudo(childCtx, contractAddr, msgBz)
}

// Check if error is out of gas error
func IsOutOfGasError(err any) (bool, string) {
	switch e := err.(type) {
	case storetypes.ErrorOutOfGas:
		return true, e.Descriptor
	case storetypes.ErrorGasOverflow:
		return true, e.Descriptor
	default:
		return false, ""
	}
}
