package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PrintEdition(
	ctx context.Context,
	minter,
	owner sdk.AccAddress,
	collectionDenom,
	tokenId string,
) (uint64, error) {
	// TODO: this is temporary, must be improved!

	nft, err := k.GetNft(ctx, collectionDenom, tokenId)
	if err != nil {
		return 0, err
	}
	if nft == nil {
		return 0, fmt.Errorf("NFT with token ID %s does not exist in collection %s", tokenId, collectionDenom)
	}

	collectionMinter, err := k.GetMinter(ctx, collectionDenom)
	if err != nil {
		return 0, err
	}

	if !minter.Equals(collectionMinter) {
		return 0, fmt.Errorf("only the collection minter can print editions")
	}

	// TODO: Charge fee if necessary

	edition := types.Edition{
		Collection: collectionDenom,
		TokenId:    tokenId,
		Seq:        nft.Editions + 1,
		Owner:      owner.String(),
	}

	if err := k.setEdition(ctx, edition); err != nil {
		return 0, fmt.Errorf("failed to set edition: %w", err)
	}

	if err := k.incrementEdition(ctx, collectionDenom, tokenId); err != nil {
		return 0, fmt.Errorf("failed to increment edition: %w", err)
	}

	return edition.Seq, nil
}

func (k Keeper) setEdition(ctx context.Context, edition types.Edition) error {
	editionKey := collections.Join3(edition.Collection, edition.TokenId, edition.Seq)
	return k.Editions.Set(ctx, editionKey, edition)
}
