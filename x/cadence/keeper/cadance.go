package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/bitsongofficial/go-bitsong/x/cadence/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// Store Keys for cadence contract s (both jailed and unjailed)
var (
	StoreKeyContracts = []byte("contracts")
)

// Get the store for the cadence contract s.
func (k Keeper) getStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), StoreKeyContracts)
}

// Set a cadence contract  address in the KV store.
func (k Keeper) SetCadenceContract(ctx sdk.Context, contract types.CadenceContract) error {
	// Get store, marshal content
	store := k.getStore(ctx)
	bz, err := k.cdc.Marshal(&contract)
	if err != nil {
		return err
	}

	// Set the contract
	store.Set([]byte(contract.ContractAddress), bz)
	return nil
}

// Check if a cadence contract  address is in the KV store.
func (k Keeper) IsCadenceContract(ctx sdk.Context, contractAddress string) bool {
	store := k.getStore(ctx)
	return store.Has([]byte(contractAddress))
}

// Get a cadence contract  address from the KV store.
func (k Keeper) GetCadenceContract(ctx sdk.Context, contractAddress string) (*types.CadenceContract, error) {
	// Check if the contract is registered
	if !k.IsCadenceContract(ctx, contractAddress) {
		return nil, types.ErrContractNotRegistered
	}

	// Get the KV store
	store := k.getStore(ctx)
	bz := store.Get([]byte(contractAddress))

	// Unmarshal the contract
	var contract types.CadenceContract
	err := k.cdc.Unmarshal(bz, &contract)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract %s: %w", contractAddress, err)
	}

	// Return the contract
	return &contract, nil
}

// Get all cadence contract  addresses from the KV store.
func (k Keeper) GetAllContracts(ctx sdk.Context) ([]types.CadenceContract, error) {
	// Get the KV store
	store := k.getStore(ctx)

	// Create iterator for contracts
	iterator := storetypes.KVStorePrefixIterator(store, []byte(nil))
	defer iterator.Close()

	// Iterate over all contracts
	contracts := []types.CadenceContract{}
	for ; iterator.Valid(); iterator.Next() {

		// Unmarshal iterator
		var contract types.CadenceContract
		err := k.cdc.Unmarshal(iterator.Value(), &contract)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	// Return array of contracts
	return contracts, nil
}

// Get all registered fee pay contracts
func (k Keeper) GetPaginatedContracts(ctx sdk.Context, pag *query.PageRequest) (*types.QueryCadenceContractsResponse, error) {
	store := k.getStore(ctx)

	// Filter and paginate all contracts
	results, pageRes, err := query.GenericFilteredPaginate(
		k.cdc,
		store,
		pag,
		func(_ []byte, value *types.CadenceContract) (*types.CadenceContract, error) {
			return value, nil
		},
		func() *types.CadenceContract {
			return &types.CadenceContract{}
		},
	)
	if err != nil {
		return nil, err
	}

	// Dereference pointer array of contracts
	var contracts []types.CadenceContract
	for _, contract := range results {
		contracts = append(contracts, *contract)
	}

	// Return paginated contracts
	return &types.QueryCadenceContractsResponse{
		CadenceContracts: contracts,
		Pagination:       pageRes,
	}, nil
}

// Remove a cadence contract  address from the KV store.
func (k Keeper) RemoveContract(ctx sdk.Context, contractAddress string) {
	store := k.getStore(ctx)
	key := []byte(contractAddress)

	if store.Has(key) {
		store.Delete(key)
	}
}

// Register a cadence contract  address in the KV store.
func (k Keeper) RegisterContract(ctx sdk.Context, senderAddress string, contractAddress string) error {
	// Check if the contract is already registered
	if k.IsCadenceContract(ctx, contractAddress) {
		return types.ErrContractAlreadyRegistered
	}

	// Ensure the sender is the contract admin or creator
	if ok, err := k.IsContractManager(ctx, senderAddress, contractAddress); !ok {
		return err
	}

	// Register contract
	return k.SetCadenceContract(ctx, types.CadenceContract{
		ContractAddress: contractAddress,
		IsJailed:        false,
	})
}

// Unregister a cadence contract  from either the jailed or unjailed KV store.
func (k Keeper) UnregisterContract(ctx sdk.Context, senderAddress string, contractAddress string) error {
	// Check if the contract is registered in either store
	if !k.IsCadenceContract(ctx, contractAddress) {
		return types.ErrContractNotRegistered
	}

	// Ensure the sender is the contract admin or creator
	if ok, err := k.IsContractManager(ctx, senderAddress, contractAddress); !ok {
		return err
	}

	// Remove contract from both stores
	k.RemoveContract(ctx, contractAddress)
	return nil
}

// Set the jail status of a cadence contract  in the KV store.
func (k Keeper) SetJailStatus(ctx sdk.Context, contractAddress string, isJailed bool) error {
	// Get the contract
	contract, err := k.GetCadenceContract(ctx, contractAddress)
	if err != nil {
		return err
	}

	// Check if the contract is already jailed or unjailed
	if contract.IsJailed == isJailed {
		if isJailed {
			return types.ErrContractAlreadyJailed
		}

		return types.ErrContractNotJailed
	}

	// Set the jail status
	contract.IsJailed = isJailed

	// Set the contract
	return k.SetCadenceContract(ctx, *contract)
}

// Set the jail status of a cadence contract  by the sender address.
func (k Keeper) SetJailStatusBySender(ctx sdk.Context, senderAddress string, contractAddress string, jailStatus bool) error {
	// Ensure the sender is the contract admin or creator
	if ok, err := k.IsContractManager(ctx, senderAddress, contractAddress); !ok {
		return err
	}

	return k.SetJailStatus(ctx, contractAddress, jailStatus)
}

// Check if the sender is the designated contract manager for the FeePay contract. If
// an admin is present, they are considered the manager. If there is no admin, the
// contract creator is considered the manager.
func (k Keeper) IsContractManager(ctx sdk.Context, senderAddress string, contractAddress string) (bool, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddress)

	// Ensure the contract is a cosm wasm contract
	if ok := k.wasmKeeper.HasContractInfo(ctx, contractAddr); !ok {
		return false, types.ErrInvalidCWContract
	}

	// Get the contract info
	contractInfo := k.wasmKeeper.GetContractInfo(ctx, contractAddr)

	// Flags for admin existence & sender being admin/creator
	adminExists := len(contractInfo.Admin) > 0
	isSenderAdmin := contractInfo.Admin == senderAddress
	isSenderCreator := contractInfo.Creator == senderAddress

	// Check if the sender is the admin or creator
	if adminExists && !isSenderAdmin {
		return false, types.ErrContractNotAdmin
	} else if !adminExists && !isSenderCreator {
		return false, types.ErrContractNotCreator
	}

	return true, nil
}
