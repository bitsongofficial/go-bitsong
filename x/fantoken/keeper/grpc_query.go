package keeper

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/bitsongofficial/go-bitsong/x/fantoken/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) FanToken(c context.Context, req *types.QueryFanTokenRequest) (*types.QueryFanTokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	token, err := k.GetFanToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "token %s not found", req.Denom)
	}
	msg, ok := token.(proto.Message)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't protomarshal %T", token)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFanTokenResponse{FanToken: any}, nil
}

func (k Keeper) FanTokens(c context.Context, req *types.QueryFanTokensRequest) (*types.QueryFanTokensResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var owner sdk.AccAddress
	var err error
	if len(req.Owner) > 0 {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid owner address (%s)", err))
		}
	}

	var tokens []types.FanTokenI
	var pageRes *query.PageResponse
	store := ctx.KVStore(k.storeKey)
	if owner == nil {
		tokenStore := prefix.NewStore(store, types.PrefixFanTokenForSymbol)
		pageRes, err = query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
			var token types.FanToken
			k.cdc.MustUnmarshalBinaryBare(value, &token)
			tokens = append(tokens, &token)
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	} else {
		tokenStore := prefix.NewStore(store, types.KeyFanTokens(owner, ""))
		pageRes, err = query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
			var symbol gogotypes.StringValue
			k.cdc.MustUnmarshalBinaryBare(value, &symbol)
			token, err := k.GetFanToken(ctx, symbol.Value)
			if err == nil {
				tokens = append(tokens, token)
			}
			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}
	result := make([]*codectypes.Any, len(tokens))
	for i, token := range tokens {
		msg, ok := token.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.Internal, "%T does not implement proto.Message", token)
		}

		var err error
		if result[i], err = codectypes.NewAnyWithValue(msg); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryFanTokensResponse{FanTokens: result, Pagination: pageRes}, nil
}

// Params return the all the parameter in tonken module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// TotalBurn return the all burn coin
func (k Keeper) TotalBurn(c context.Context, req *types.QueryTotalBurnRequest) (*types.QueryTotalBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryTotalBurnResponse{
		BurnedCoins: k.GetAllBurnCoin(ctx),
	}, nil
}
