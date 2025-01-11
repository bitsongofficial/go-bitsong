package types

import (
	fmt "fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	db "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
)

const (
	// ModuleName defines the module name
	ModuleName = "smartaccount"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	KeySeparator = "|"

	AttributeValueCategory        = ModuleName
	AttributeKeyAuthenticatorType = "authenticator_type"
	AttributeKeyAuthenticatorId   = "authenticator_id"

	AtrributeKeyIsSmartAccountActive = "is_smart_account_active"

	AttributeKeyAccountSequenceAuthenticator = "authenticator_acc_seq"
	AttributeKeySignatureAuthenticator       = "authenticator_signature"
)

var (
	// Store prefix keys
	KeyNextAccountAuthenticatorIdPrefix = []byte{0x01}
	KeyAccountAuthenticatorsPrefix      = []byte{0x02}

	// Parameter keys
	KeyMaximumUnauthenticatedGas = []byte("MaximumUnauthenticatedGas")
	KeyIsSmartAccountActive      = []byte("IsSmartAccountActive")
	KeyCircuitBreakerControllers = []byte("CircuitBreakerControllers")
)

func KeyAccount(account sdk.AccAddress) []byte {
	return BuildKey(KeyAccountAuthenticatorsPrefix, account.String())
}

func KeyAccountId(account sdk.AccAddress, id uint64) []byte {
	return BuildKey(KeyAccountAuthenticatorsPrefix, account.String(), id)
}

func KeyNextAccountAuthenticatorId() []byte {
	return BuildKey(KeyNextAccountAuthenticatorIdPrefix)
}

func KeyAccountAuthenticatorsPrefixId() []byte {
	return BuildKey(KeyAccountAuthenticatorsPrefix)
}

// BuildKey creates a key by concatenating the provided elements with the key separator.
func BuildKey(elements ...interface{}) []byte {
	strElements := make([]string, len(elements))
	for i, element := range elements {
		strElements[i] = fmt.Sprint(element)
	}
	return []byte(strings.Join(strElements, KeySeparator) + KeySeparator)
}

func noStopFn([]byte) bool {
	return false
}

// MustSet runs store.Set(key, proto.Marshal(value))
// but panics on any error.
func MustSet(storeObj storetypes.KVStore, key []byte, value proto.Message) {
	bz, err := proto.Marshal(value)
	if err != nil {
		panic(err)
	}

	storeObj.Set(key, bz)
}

// MustGet
// GatherValuesFromStorePrefix is a decorator around GatherValuesFromStorePrefixWithKeyParser. It overwrites the parse function to
// disable parsing keys, only keeping values
func GatherValuesFromStorePrefix[T any](storeObj storetypes.KVStore, prefix []byte, parseValue func([]byte) (T, error)) ([]T, error) {
	// Replace a callback with the one that takes both key and value
	// but ignores the key.
	parseOnlyValue := func(_ []byte, value []byte) (T, error) {
		return parseValue(value)
	}
	return GatherValuesFromStorePrefixWithKeyParser(storeObj, prefix, parseOnlyValue)
}

// GatherValuesFromStorePrefixWithKeyParser is a helper function that gathers values from a given store prefix. While iterating through
// the entries, it parses both key and the value using the provided parse function to return the desired type.
// Returns error if:
// - the parse function returns an error.
// - internal database error
func GatherValuesFromStorePrefixWithKeyParser[T any](storeObj storetypes.KVStore, prefix []byte, parse func(key []byte, value []byte) (T, error)) ([]T, error) {
	iterator := storetypes.KVStorePrefixIterator(storeObj, prefix)
	defer iterator.Close()
	return gatherValuesFromIteratorWithKeyParser(iterator, parse, noStopFn)
}

func gatherValuesFromIteratorWithKeyParser[T any](iterator db.Iterator, parse func(key []byte, value []byte) (T, error), stopFn func([]byte) bool) ([]T, error) {
	values := []T{}
	for ; iterator.Valid(); iterator.Next() {
		if stopFn(iterator.Key()) {
			break
		}
		val, err := parse(iterator.Key(), iterator.Value())
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

// Get returns a value at key by mutating the result parameter. Returns true if the value was found and the
// result mutated correctly. If the value is not in the store, returns false.
// Returns error only when database or serialization errors occur. (And when an error occurs, returns false)
func Get(store storetypes.KVStore, key []byte, result proto.Message) (found bool, err error) {
	b := store.Get(key)
	if b == nil {
		return false, nil
	}
	if err := proto.Unmarshal(b, result); err != nil {
		return true, err
	}
	return true, nil
}
