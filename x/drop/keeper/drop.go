package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/bitsongofficial/go-bitsong/x/drop/types"

	nfttypes "github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (k Keeper) CreateDrop(
	ctx context.Context,
	collectionDenom string,
	maxAvailable uint64,
	rules []types.Rule,
) error {
	// 1. check if collections is not in our drop store
	hasDrop, err := k.HasDrop(ctx, collectionDenom)
	if err != nil {
		return fmt.Errorf("failed to check drop existence: %w", err)
	}

	if hasDrop {
		return fmt.Errorf("drop already exists for collection %s", collectionDenom)
	}

	// 2. check if collections has nfts, if yes return error, only empty collections can be dropped
	collSupply := k.nftKeeper.GetSupply(ctx, collectionDenom)

	if collSupply.GT(math.ZeroInt()) {
		return fmt.Errorf("collection %s is not empty, cannot create drop", collectionDenom)
	}

	// 3. check drop max available (max is MaxNftsInCollection)
	if math.NewUint(maxAvailable).GT(math.NewUint(nfttypes.MaxNftsInCollection)) {
		return fmt.Errorf("max available %d exceeds max allowed %d", maxAvailable, nfttypes.MaxNftsInCollection)
	}

	// 4. check max rules (5)
	if len(rules) > types.MaxRulesPerDrop {
		return fmt.Errorf("number of rules %d exceeds max allowed %d", len(rules), types.MaxRulesPerDrop)
	}

	if err := k.validateRules(rules); err != nil {
		return err
	}

	// store drop
	// store rules

	return nil
}

func (k Keeper) GetDrop(ctx context.Context, collectionDenom string) (types.Drop, error) {
	return k.Drops.Get(ctx, collectionDenom)
}

func (k Keeper) HasDrop(ctx context.Context, collectionDenom string) (bool, error) {
	return k.Drops.Has(ctx, collectionDenom)
}

func (k Keeper) validateRules(rules []types.Rule) error {
	return nil
}
