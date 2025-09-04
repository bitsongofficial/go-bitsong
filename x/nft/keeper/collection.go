package keeper

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	tmcrypto "github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

const MaxDenomLength = 43

func (k Keeper) CreateCollection(ctx sdk.Context, creator sdk.AccAddress, coll types.Collection) (denom string, err error) {
	denom, err = k.validateCollectionDenom(ctx, creator, coll.Symbol)
	if err != nil {
		return "", err
	}

	// TODO: charge fee

	if err := k.setCollection(ctx, denom, coll); err != nil {
		return "", err
	}

	return denom, nil
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) math.Int {
	supply, err := k.Supply.Get(ctx, denom)
	if err != nil {
		return math.ZeroInt()
	}

	return supply
}

func (k Keeper) HasSupply(ctx sdk.Context, denom string) bool {
	has, err := k.Supply.Has(ctx, denom)
	return has && err == nil
}

func (k Keeper) setSupply(ctx sdk.Context, denom string, supply math.Int) error {
	return k.Supply.Set(ctx, denom, supply)
}

func (k Keeper) incrementSupply(ctx sdk.Context, denom string) error {
	supply := k.GetSupply(ctx, denom)
	supply = supply.Add(math.NewInt(1))

	return k.setSupply(ctx, denom, supply)
}

func (k Keeper) createCollectionDenom(creator sdk.AccAddress, symbol string) string {
	// TODO: if necessary add a salt field

	bz := []byte(fmt.Sprintf("%s/%s", creator.String(), symbol))
	return fmt.Sprintf("nft%x", tmcrypto.AddressHash(bz))
}

func (k Keeper) validateCollectionDenom(ctx sdk.Context, creator sdk.AccAddress, symbol string) (string, error) {
	denom := k.createCollectionDenom(creator, symbol)

	if err := sdk.ValidateDenom(denom); err != nil {
		return "", err
	}

	if k.HasCollection(ctx, denom) {
		return "", types.ErrCollectionAlreadyExists
	}

	return denom, nil
}

func (k Keeper) setCollection(ctx sdk.Context, denom string, coll types.Collection) error {
	return k.Collections.Set(ctx, denom, coll)
}

func (k Keeper) getCollection(ctx sdk.Context, denom string) (types.Collection, error) {
	coll, err := k.Collections.Get(ctx, denom)
	if err != nil {
		return types.Collection{}, types.ErrCollectionNotFound
	}

	return coll, nil
}

func (k Keeper) HasCollection(ctx sdk.Context, denom string) bool {
	has, err := k.Collections.Has(ctx, denom)
	return has && err == nil
}

func LengthDenomPrefix(bz []byte) ([]byte, error) {
	bzLen := len(bz)
	if bzLen == 0 {
		return bz, nil
	}

	if bzLen > MaxDenomLength {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidType, "denom length should be max %d bytes, got %d", MaxDenomLength, bzLen)
	}

	return append([]byte{byte(bzLen)}, bz...), nil
}

func MustLengthDenomPrefix(bz []byte) []byte {
	res, err := LengthDenomPrefix(bz)
	if err != nil {
		panic(err)
	}

	return res
}

func SplitNftLengthPrefixedKey(key []byte) (denom, tokenId []byte, err error) {
	parts := bytes.SplitN(key, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid composite key format: expected 2 parts, got %d", len(parts))
	}

	denomLen := len(parts[0])

	if denomLen > MaxDenomLength {
		return nil, nil, errors.Wrapf(sdkerrors.ErrInvalidType, "decoded denom key length %d exceeds max allowed length %d", denomLen, MaxDenomLength)
	}

	if len(key)-1 < denomLen {
		return nil, nil, fmt.Errorf("key is malformed: length prefix %d is greater than tokenId bytes %d", denomLen, len(key)-1)
	}

	denom = parts[0]
	tokenId = parts[1]

	return denom, tokenId, nil
}

func MustSplitNftLengthPrefixedKey(key []byte) (denom, tokenId []byte) {
	denom, tokenId, err := SplitNftLengthPrefixedKey(key)
	if err != nil {
		panic(err)
	}

	return denom, tokenId
}
