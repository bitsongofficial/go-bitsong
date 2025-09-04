package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/store/prefix"
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Collection(ctx context.Context, req *types.QueryCollectionRequest) (*types.QueryCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	coll, err := k.Collections.Get(sdkCtx, req.Collection)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "collection %s not found", req.Collection)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCollectionResponse{
		Collection: &coll,
	}, nil
}

func (k Keeper) OwnerOf(ctx context.Context, req *types.QueryOwnerOfRequest) (*types.QueryOwnerOfResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection cannot be empty")
	}
	if req.TokenId == "" {
		return nil, status.Error(codes.InvalidArgument, "token_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	nft, err := k.NFTs.Get(sdkCtx, collections.Join(req.Collection, req.TokenId))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "nft with collection %s and token_id %s not found", req.Collection, req.TokenId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOwnerOfResponse{
		Owner: nft.Owner,
	}, nil
}

func (k Keeper) NumTokens(ctx context.Context, req *types.QueryNumTokensRequest) (*types.QueryNumTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	supply := k.GetSupply(sdkCtx, req.Collection)

	return &types.QueryNumTokensResponse{
		Count: supply.Uint64(),
	}, nil
}

func (k Keeper) NftInfo(ctx context.Context, req *types.QueryNftInfoRequest) (*types.QueryNftInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection cannot be empty")
	}
	if req.TokenId == "" {
		return nil, status.Error(codes.InvalidArgument, "token_id cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	nft, err := k.NFTs.Get(sdkCtx, collections.Join(req.Collection, req.TokenId))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "nft with collection %s and token_id %s not found", req.Collection, req.TokenId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryNftInfoResponse{
		Nft: &nft,
	}, nil
}

func (k Keeper) Nfts(ctx context.Context, req *types.QueryNftsRequest) (*types.QueryNftsResponse, error) {
	if req == nil || req.Collection == "" {
		return nil, status.Error(codes.InvalidArgument, "collection cannot be empty")
	}

	nfts, pageRes, err := query.CollectionPaginate(
		ctx,
		k.NFTs,
		req.Pagination,
		func(key collections.Pair[string, string], value types.Nft) (types.Nft, error) {
			return value, nil
		},
		query.WithCollectionPaginationPairPrefix[string, string](req.Collection),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryNftsResponse{
		Nfts:       nfts,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) AllNftsByOwner(ctx context.Context, req *types.QueryAllNftsByOwnerRequest) (*types.QueryAllNftsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner address: %s", err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := prefix.NewStore(
		sdkCtx.KVStore(k.storeKey),
		append(types.NFTsByOwnerPrefix, address.MustLengthPrefix(owner)...),
	)

	var nfts []types.Nft

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		denom, tokenId := types.MustSplitNftLengthPrefixedKey(key)

		nft, err := k.NFTs.Get(ctx, collections.Join(string(denom), string(tokenId)))
		if err != nil {
			return err
		}

		nfts = append(nfts, nft)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	/*iter, err := k.NFTs.Indexes.Owner.MatchExact(ctx, owner)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iter.Close()

	nfts, err := indexes.CollectValues(ctx, k.NFTs, iter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}*/
	return &types.QueryAllNftsByOwnerResponse{
		Nfts:       nfts,
		Pagination: pageRes,
	}, nil
}
