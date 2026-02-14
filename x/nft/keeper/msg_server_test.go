package keeper_test

import (
	"github.com/bitsongofficial/go-bitsong/x/nft/types"
)

func (suite *KeeperTestSuite) TestMsgCreateCollection() {
	testCases := []struct {
		name             string
		msg              *types.MsgCreateCollection
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "valid message",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  minter1.String(),
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol + "0",
				Uri:     testCollection1.Uri,
			},
		},
		{
			name: "invalid creator address",
			msg: &types.MsgCreateCollection{
				Creator: "invalid_address",
				Minter:  minter1.String(),
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol,
				Uri:     testCollection1.Uri,
			},
			expectError:      true,
			expectedErrorMsg: "invalid creator address: decoding bech32 failed",
		},
		{
			name: "should fail on empty creator address",
			msg: &types.MsgCreateCollection{
				Creator: "",
				Minter:  minter1.String(),
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol,
				Uri:     testCollection1.Uri,
			},
			expectError:      true,
			expectedErrorMsg: "invalid creator address: empty address string",
		},
		{
			name: "invalid minter address",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  "invalid_address",
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol,
				Uri:     testCollection1.Uri,
			},
			expectError:      true,
			expectedErrorMsg: "invalid minter address: decoding bech32 failed",
		},
		{
			name: "valid empty minter address",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  "",
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol,
				Uri:     testCollection1.Uri,
			},
		},
		{
			name: "should fail on empty symbol",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  minter1.String(),
				Name:    testCollection1.Name,
				Symbol:  "",
				Uri:     testCollection1.Uri,
			},
			expectError:      true,
			expectedErrorMsg: "symbol cannot be empty",
		},
		{
			name: "valid on empty name",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  minter1.String(),
				Name:    "",
				Symbol:  testCollection1.Symbol + "1",
				Uri:     testCollection1.Uri,
			},
		},
		{
			name: "valid on empty uri",
			msg: &types.MsgCreateCollection{
				Creator: creator1.String(),
				Minter:  minter1.String(),
				Name:    testCollection1.Name,
				Symbol:  testCollection1.Symbol + "2",
				Uri:     "",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.CreateCollection(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				if tc.expectedErrorMsg != "" {
					suite.Require().Contains(err.Error(), tc.expectedErrorMsg)
				}
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgMintNFT() {
	// first create a collection to mint into
	createCollectionMsg := &types.MsgCreateCollection{
		Creator: creator1.String(),
		Minter:  minter1.String(),
		Name:    testCollection1.Name,
		Symbol:  testCollection1.Symbol,
		Uri:     testCollection1.Uri,
	}
	res, err := suite.msgServer.CreateCollection(suite.ctx, createCollectionMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.Denom)

	testCases := []struct {
		name             string
		msg              *types.MsgMintNFT
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "valid message",
			msg: &types.MsgMintNFT{
				Minter:     minter1.String(),
				Recipient:  owner1.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
				Name:       testNft1.Name,
				Uri:        testNft1.Uri,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.MintNFT(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				if tc.expectedErrorMsg != "" {
					suite.Require().Contains(err.Error(), tc.expectedErrorMsg)
				}
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgMintNFT_MaxLimit() {
	createCollectionMsg := &types.MsgCreateCollection{
		Creator: creator1.String(),
		Minter:  minter1.String(),
		Name:    testCollection1.Name,
		Symbol:  testCollection1.Symbol,
		Uri:     testCollection1.Uri,
	}
	res, err := suite.msgServer.CreateCollection(suite.ctx, createCollectionMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.Denom)

	for i := 0; i < types.MaxNftsInCollection; i++ {
		mintMsg := &types.MsgMintNFT{
			Minter:     minter1.String(),
			Recipient:  owner1.String(),
			Collection: res.Denom,
			TokenId:    "token" + string(rune(i)),
			Name:       "Token " + string(rune(i)),
			Uri:        "http://example.com/token" + string(rune(i)),
		}
		_, err := suite.msgServer.MintNFT(suite.ctx, mintMsg)
		suite.Require().NoError(err)
	}

	// attempt to mint one more NFT beyond the limit
	mintMsg := &types.MsgMintNFT{
		Minter:     minter1.String(),
		Recipient:  owner1.String(),
		Collection: res.Denom,
		TokenId:    "token_overflow",
		Name:       "Token Overflow",
		Uri:        "http://example.com/token_overflow",
	}
	_, err = suite.msgServer.MintNFT(suite.ctx, mintMsg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "max supply reached")
}

func (suite *KeeperTestSuite) TestMsgSendNFT() {
	createCollectionMsg := &types.MsgCreateCollection{
		Creator: creator1.String(),
		Minter:  minter1.String(),
		Name:    testCollection1.Name,
		Symbol:  testCollection1.Symbol,
		Uri:     testCollection1.Uri,
	}
	res, err := suite.msgServer.CreateCollection(suite.ctx, createCollectionMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.Denom)

	mintMsg := &types.MsgMintNFT{
		Minter:     minter1.String(),
		Recipient:  owner1.String(),
		Collection: res.Denom,
		TokenId:    testNft1.TokenId,
		Name:       testNft1.Name,
		Uri:        testNft1.Uri,
	}
	_, err = suite.msgServer.MintNFT(suite.ctx, mintMsg)
	suite.Require().NoError(err)

	testCases := []struct {
		name             string
		msg              *types.MsgSendNFT
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "valid message",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
		},
		{
			name: "should fail on same sender and recipient",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  owner1.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "cannot transfer NFT to the same owner",
		},
		{
			name: "invalid sender address",
			msg: &types.MsgSendNFT{
				Sender:     "invalid_address",
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid sender address: decoding bech32 failed",
		},
		{
			name: "should fail on empty sender address",
			msg: &types.MsgSendNFT{
				Sender:     "",
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid sender address: empty address string",
		},
		{
			name: "invalid recipient address",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  "invalid_address",
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid recipient address: decoding bech32 failed",
		},
		{
			name: "should fail on empty recipient address",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  "",
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid recipient address: empty address string",
		},
		{
			name: "should fail on non existing collection",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  owner2.String(),
				Collection: "non_existing_collection",
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "collection or token_id does not exist",
		},
		{
			name: "should fail on non existing token",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    "non_existing_token",
			},
			expectError:      true,
			expectedErrorMsg: "collection or token_id does not exist",
		},
		{
			name: "should fail when sender is not the owner",
			msg: &types.MsgSendNFT{
				Sender:     owner1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "only the owner can transfer the NFT",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.SendNFT(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				if tc.expectedErrorMsg != "" {
					suite.Require().Contains(err.Error(), tc.expectedErrorMsg)
				}
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgPrintEdition() {
	createCollectionMsg := &types.MsgCreateCollection{
		Creator: creator1.String(),
		Minter:  minter1.String(),
		Name:    testCollection1.Name,
		Symbol:  testCollection1.Symbol + "3",
		Uri:     testCollection1.Uri,
	}
	res, err := suite.msgServer.CreateCollection(suite.ctx, createCollectionMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.Denom)

	mintMsg := &types.MsgMintNFT{
		Minter:     minter1.String(),
		Recipient:  owner1.String(),
		Collection: res.Denom,
		TokenId:    testNft1.TokenId,
		Name:       testNft1.Name,
		Uri:        testNft1.Uri,
	}
	_, err = suite.msgServer.MintNFT(suite.ctx, mintMsg)
	suite.Require().NoError(err)

	testCases := []struct {
		name             string
		msg              *types.MsgPrintEdition
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "valid message",
			msg: &types.MsgPrintEdition{
				Minter:     minter1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
		},
		{
			name: "invalid minter address",
			msg: &types.MsgPrintEdition{
				Minter:     "invalid_address",
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid minter address: decoding bech32 failed",
		},
		{
			name: "should fail on empty minter address",
			msg: &types.MsgPrintEdition{
				Minter:     "",
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid minter address: empty address string",
		},
		{
			name: "invalid recipient address",
			msg: &types.MsgPrintEdition{
				Minter:     minter1.String(),
				Recipient:  "invalid_address",
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid recipient address: decoding bech32 failed",
		},
		{
			name: "should fail on empty recipient address",
			msg: &types.MsgPrintEdition{
				Minter:     minter1.String(),
				Recipient:  "",
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "invalid recipient address: empty address string",
		},
		{
			name: "should fail on non existing collection",
			msg: &types.MsgPrintEdition{
				Minter:     minter1.String(),
				Recipient:  owner2.String(),
				Collection: "non_existing_collection",
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "NFT with token ID",
		},
		{
			name: "should fail on non existing token",
			msg: &types.MsgPrintEdition{
				Minter:     minter1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    "non_existing_token",
			},
			expectError:      true,
			expectedErrorMsg: "NFT with token ID",
		},
		{
			name: "should fail when minter is not the collection minter",
			msg: &types.MsgPrintEdition{
				Minter:     creator1.String(),
				Recipient:  owner2.String(),
				Collection: res.Denom,
				TokenId:    testNft1.TokenId,
			},
			expectError:      true,
			expectedErrorMsg: "only the collection minter can print editions",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.PrintEdition(suite.ctx, tc.msg)
			if tc.expectError {
				suite.Require().Error(err)
				if tc.expectedErrorMsg != "" {
					suite.Require().Contains(err.Error(), tc.expectedErrorMsg)
				}
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgPrintNFT_MaxLimit() {
	createCollectionMsg := &types.MsgCreateCollection{
		Creator: creator1.String(),
		Minter:  minter1.String(),
		Name:    testCollection1.Name,
		Symbol:  testCollection1.Symbol,
		Uri:     testCollection1.Uri,
	}
	res, err := suite.msgServer.CreateCollection(suite.ctx, createCollectionMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.Denom)

	mintMsg := &types.MsgMintNFT{
		Minter:     minter1.String(),
		Recipient:  owner1.String(),
		Collection: res.Denom,
		TokenId:    testNft1.TokenId,
		Name:       testNft1.Name,
		Uri:        testNft1.Uri,
	}
	_, err = suite.msgServer.MintNFT(suite.ctx, mintMsg)
	suite.Require().NoError(err)

	printEditionMsg := &types.MsgPrintEdition{
		Minter:     minter1.String(),
		Recipient:  owner2.String(),
		Collection: res.Denom,
		TokenId:    testNft1.TokenId,
	}

	for i := 0; i < types.MaxEditions; i++ {
		_, err := suite.msgServer.PrintEdition(suite.ctx, printEditionMsg)
		suite.Require().NoError(err)
	}

	// attempt to print one more edition beyond the limit
	_, err = suite.msgServer.PrintEdition(suite.ctx, printEditionMsg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "max editions reached")
}
