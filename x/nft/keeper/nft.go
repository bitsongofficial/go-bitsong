package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintNFT(
	ctx context.Context,
	minter sdk.AccAddress,
	owner sdk.AccAddress,
	collectionDenom,
	tokenId,
	name,
	description,
	uri string,
) error {
	if err := k.validateNftMetadata(tokenId, name, description, uri); err != nil {
		return err
	}

	nftKey := collections.Join(collectionDenom, tokenId)
	has, err := k.NFTs.Has(ctx, nftKey)
	if err != nil {
		return fmt.Errorf("failed to check NFT: %w", err)
	}
	if has {
		return fmt.Errorf("NFT with token ID %s already exists in collection %s", tokenId, collectionDenom)
	}

	coll, err := k.Collections.Get(ctx, collectionDenom)
	if err != nil {
		return types.ErrCollectionNotFound
	}

	if coll.Minter == "" {
		return fmt.Errorf("minting disabled for this collection")
	}

	collectionMinter, err := sdk.AccAddressFromBech32(coll.Minter)
	if err != nil {
		return fmt.Errorf("invalid minter address: %w", err)
	}

	if !minter.Equals(collectionMinter) {
		return fmt.Errorf("only the collection minter can mint NFTs")
	}

	// TODO: Charge fee if necessary

	nft := types.Nft{
		Collection:  collectionDenom,
		TokenId:     tokenId,
		Name:        name,
		Description: description,
		Uri:         uri,
		Owner:       owner.String(),
		Editions:    0,
	}

	if err := k.setNft(ctx, nft); err != nil {
		return fmt.Errorf("failed to set NFT: %w", err)
	}

	// TODO: add events

	return k.incrementSupply(ctx, collectionDenom)
}

func (k Keeper) SendNft(ctx context.Context, fromAddr, toAddr sdk.AccAddress, collectionDenom, tokenId string) error {
	err := k.changeNftOwner(ctx, fromAddr, toAddr, collectionDenom, tokenId)
	if err != nil {
		return err
	}

	// Same as https://github.com/cosmos/cosmos-sdk/blob/v0.53.4/x/bank/keeper/send.go
	// Create account if recipient does not exist.
	//
	// NOTE: This should ultimately be removed in favor a more flexible approach
	// such as delegated fee messages.
	accExists := k.ak.HasAccount(ctx, toAddr)
	if !accExists {
		defer telemetry.IncrCounter(1, "new", "account")
		k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, toAddr))
	}

	// Same as https://github.com/cosmos/cosmos-sdk/blob/v0.53.4/x/bank/keeper/send.go
	// bech32 encoding is expensive! Only do it once for fromAddr
	fromAddrString := fromAddr.String()
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransferNft,
			sdk.NewAttribute(types.AttributeKeyReceiver, toAddr.String()),
			sdk.NewAttribute(types.AttributeKeySender, fromAddrString),
			sdk.NewAttribute(types.AttributeKeyCollection, collectionDenom),
			sdk.NewAttribute(types.AttributeKeyTokenId, tokenId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddrString),
		),
	})

	return nil
}

func (k Keeper) GetNft(ctx context.Context, collectionDenom, tokenId string) (*types.Nft, error) {
	nftKey := collections.Join(collectionDenom, tokenId)
	has, err := k.NFTs.Has(ctx, nftKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check NFT: %w", err)
	}
	if !has {
		return nil, fmt.Errorf("NFT with token ID %s does not exist in collection %s", tokenId, collectionDenom)
	}

	nft, err := k.NFTs.Get(ctx, nftKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT: %w", err)
	}

	return &nft, nil
}

func (k Keeper) validateNftMetadata(tokenId, name, description, uri string) error {
	if strings.TrimSpace(tokenId) == "" {
		return fmt.Errorf("token ID cannot be empty")
	}

	if len(tokenId) > types.MaxTokenIdLength {
		return fmt.Errorf("token ID length exceeds maximum of %d", types.MaxTokenIdLength)
	}

	if len(name) > types.MaxNameLength {
		return fmt.Errorf("name length exceeds maximum of %d", types.MaxNameLength)

	}
	if len(description) > types.MaxDescriptionLength {
		return fmt.Errorf("description length exceeds maximum of %d", types.MaxDescriptionLength)
	}

	if len(uri) > types.MaxURILength {
		return fmt.Errorf("URI length exceeds maximum of %d", types.MaxURILength)
	}

	return nil
}

func (k Keeper) createNftDenom(ctx context.Context, collectionDenom string) string {
	supply := k.GetSupply(ctx, collectionDenom)
	return fmt.Sprintf("%s-%d", collectionDenom, supply.Uint64()+1)
}

func (k Keeper) setNft(ctx context.Context, nft types.Nft) error {
	pk := collections.Join(nft.Collection, nft.TokenId)
	return k.NFTs.Set(ctx, pk, nft)
}

func (k Keeper) changeNftOwner(ctx context.Context, oldOwner, newOwner sdk.AccAddress, collectionDenom string, tokenId string) error {
	nft, err := k.NFTs.Get(ctx, collections.Join(collectionDenom, tokenId))
	if err != nil {
		return fmt.Errorf("failed to get NFT: %w", err)
	}

	if nft.Owner != oldOwner.String() {
		return fmt.Errorf("only the owner can transfer the NFT")
	}

	nft.Owner = newOwner.String()
	err = k.setNft(ctx, nft)
	if err != nil {
		return fmt.Errorf("failed to set NFT Owner: %w", err)
	}

	// emit nft received event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		types.NewNftReceivedEvent(newOwner, collectionDenom, tokenId),
	)

	return nil
}

func (k Keeper) incrementEdition(ctx context.Context, collectionDenom, tokenId string) error {
	nft, err := k.NFTs.Get(ctx, collections.Join(collectionDenom, tokenId))
	if err != nil {
		return fmt.Errorf("failed to get NFT: %w", err)
	}

	nft.Editions += 1
	return k.setNft(ctx, nft)
}
