package keeper

import (
	"context"

	"github.com/bitsongofficial/go-bitsong/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateAuction(goCtx context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
	// ctx.EventManager().EmitTypedEvent(&types.EventCreateAuction{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// create auction
	// pub fn create_auction(
	//     program_id: &Pubkey,
	//     accounts: &[AccountInfo],
	//     args: CreateAuctionArgs,
	//     instant_sale_price: Option<u64>,
	//     name: Option<AuctionName>,
	// ) -> ProgramResult {
	//     msg!("+ Processing CreateAuction");
	//     let accounts = parse_accounts(program_id, accounts)?;

	//     let auction_path = [
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         &args.resource.to_bytes(),
	//     ];

	//     // Derive the address we'll store the auction in, and confirm it matches what we expected the
	//     // user to provide.
	//     let (auction_key, bump) = Pubkey::find_program_address(&auction_path, program_id);
	//     if auction_key != *accounts.auction.key {
	//         return Err(AuctionError::InvalidAuctionAccount.into());
	//     }
	//     // The data must be large enough to hold at least the number of winners.
	//     let auction_size = match args.winners {
	//         WinnerLimit::Capped(n) => {
	//             mem::size_of::<Bid>() * BidState::max_array_size_for(n) + BASE_AUCTION_DATA_SIZE
	//         }
	//         WinnerLimit::Unlimited(_) => BASE_AUCTION_DATA_SIZE,
	//     };

	//     let bid_state = match args.winners {
	//         WinnerLimit::Capped(n) => BidState::new_english(n),
	//         WinnerLimit::Unlimited(_) => BidState::new_open_edition(),
	//     };

	//     if let Some(gap_tick) = args.gap_tick_size_percentage {
	//         if gap_tick > 100 {
	//             return Err(AuctionError::InvalidGapTickSizePercentage.into());
	//         }
	//     }

	//     // Create auction account with enough space for a winner tracking.
	//     create_or_allocate_account_raw(
	//         *program_id,
	//         accounts.auction,
	//         accounts.rent,
	//         accounts.system,
	//         accounts.payer,
	//         auction_size,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             &args.resource.to_bytes(),
	//             &[bump],
	//         ],
	//     )?;

	//     let auction_ext_bump = assert_derivation(
	//         program_id,
	//         accounts.auction_extended,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             &args.resource.to_bytes(),
	//             EXTENDED.as_bytes(),
	//         ],
	//     )?;

	//     create_or_allocate_account_raw(
	//         *program_id,
	//         accounts.auction_extended,
	//         accounts.rent,
	//         accounts.system,
	//         accounts.payer,
	//         MAX_AUCTION_DATA_EXTENDED_SIZE,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             &args.resource.to_bytes(),
	//             EXTENDED.as_bytes(),
	//             &[auction_ext_bump],
	//         ],
	//     )?;

	//     // Configure extended
	//     AuctionDataExtended {
	//         total_uncancelled_bids: 0,
	//         tick_size: args.tick_size,
	//         gap_tick_size_percentage: args.gap_tick_size_percentage,
	//         instant_sale_price,
	//         name,
	//     }
	//     .serialize(&mut *accounts.auction_extended.data.borrow_mut())?;

	//     // Configure Auction.
	//     AuctionData {
	//         authority: args.authority,
	//         bid_state: bid_state,
	//         end_auction_at: args.end_auction_at,
	//         end_auction_gap: args.end_auction_gap,
	//         ended_at: None,
	//         last_bid: None,
	//         price_floor: args.price_floor,
	//         state: AuctionState::create(),
	//         token_mint: args.token_mint,
	//     }
	//     .serialize(&mut *accounts.auction.data.borrow_mut())?;

	//     Ok(())
	// }

	return &types.MsgCreateAuctionResponse{}, nil
}

func (m msgServer) StartAuction(goCtx context.Context, msg *types.MsgStartAuction) (*types.MsgStartAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// start auction codebase
	// pub fn start_auction<'a, 'b: 'a>(
	//     program_id: &Pubkey,
	//     accounts: &'a [AccountInfo<'b>],
	//     args: StartAuctionArgs,
	// ) -> ProgramResult {
	//     msg!("+ Processing StartAuction");
	//     let accounts = parse_accounts(program_id, accounts)?;
	//     let clock = Clock::from_account_info(accounts.clock_sysvar)?;

	//     // Derive auction address so we can make the modifications necessary to start it.
	//     assert_derivation(
	//         program_id,
	//         accounts.auction,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             &args.resource.as_ref(),
	//         ],
	//     )?;

	//     // Initialise a new auction. The end time is calculated relative to now.
	//     let mut auction = AuctionData::from_account_info(accounts.auction)?;

	//     // Check authority is correct.
	//     if auction.authority != *accounts.authority.key {
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     // Calculate the relative end time.
	//     let ended_at = if let Some(end_auction_at) = auction.end_auction_at {
	//         match clock.unix_timestamp.checked_add(end_auction_at) {
	//             Some(val) => Some(val),
	//             None => return Err(AuctionError::NumericalOverflowError.into()),
	//         }
	//     } else {
	//         None
	//     };

	//     AuctionData {
	//         ended_at,
	//         state: auction.state.start()?,
	//         ..auction
	//     }
	//     .serialize(&mut *accounts.auction.data.borrow_mut())?;

	//     Ok(())
	// }

	// ctx.EventManager().EmitTypedEvent(&types.EventStartAuction{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	return &types.MsgStartAuctionResponse{}, nil
}

func (m msgServer) SetAuctionAuthority(goCtx context.Context, msg *types.MsgSetAuctionAuthority) (*types.MsgSetAuctionAuthorityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// ctx.EventManager().EmitTypedEvent(&types.EventSetAuctionAuthority{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// 	pub fn set_authority(program_id: &Pubkey, accounts: &[AccountInfo]) -> ProgramResult {
	//     msg!("+ Processing SetAuthority");
	//     let account_iter = &mut accounts.iter();
	//     let auction_act = next_account_info(account_iter)?;
	//     let current_authority = next_account_info(account_iter)?;
	//     let new_authority = next_account_info(account_iter)?;

	//     let mut auction = AuctionData::from_account_info(auction_act)?;
	//     assert_owned_by(auction_act, program_id)?;

	//     if auction.authority != *current_authority.key {
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     if !current_authority.is_signer {
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     // Make sure new authority actually exists in some form.
	//     if new_authority.data_is_empty() || new_authority.lamports() == 0 {
	//         msg!("Disallowing new authority because it does not exist.");
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     auction.authority = *new_authority.key;
	//     auction.serialize(&mut *auction_act.data.borrow_mut())?;
	//     Ok(())
	// }

	return &types.MsgSetAuctionAuthorityResponse{}, nil
}

func (m msgServer) EndAuction(goCtx context.Context, msg *types.MsgEndAuction) (*types.MsgEndAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// ctx.EventManager().EmitTypedEvent(&types.EventEndAuction{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// pub fn end_auction<'a, 'b: 'a>(
	//     program_id: &Pubkey,
	//     accounts: &'a [AccountInfo<'b>],
	//     args: EndAuctionArgs,
	// ) -> ProgramResult {
	//     msg!("+ Processing EndAuction");
	//     let accounts = parse_accounts(program_id, accounts)?;
	//     let clock = Clock::from_account_info(accounts.clock_sysvar)?;

	//     assert_derivation(
	//         program_id,
	//         accounts.auction,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             &args.resource.as_ref(),
	//         ],
	//     )?;

	//     // End auction.
	//     let mut auction = AuctionData::from_account_info(accounts.auction)?;

	//     // Check authority is correct.
	//     if auction.authority != *accounts.authority.key {
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     // As long as it hasn't already ended.
	//     if auction.ended_at.is_some() {
	//         return Err(AuctionError::AuctionTransitionInvalid.into());
	//     }

	//     AuctionData {
	//         ended_at: Some(clock.unix_timestamp),
	//         state: auction.state.end()?,
	//         price_floor: reveal(auction.price_floor, args.reveal)?,
	//         ..auction
	//     }
	//     .serialize(&mut *accounts.auction.data.borrow_mut())?;

	//     Ok(())
	// }

	return &types.MsgEndAuctionResponse{}, nil
}

func (m msgServer) PlaceBid(goCtx context.Context, msg *types.MsgPlaceBid) (*types.MsgPlaceBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// ctx.EventManager().EmitTypedEvent(&types.EventPlaceBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// #[allow(clippy::absurd_extreme_comparisons)]
	// pub fn place_bid<'r, 'b: 'r>(
	//     program_id: &Pubkey,
	//     accounts: &'r [AccountInfo<'b>],
	//     args: PlaceBidArgs,
	// ) -> ProgramResult {
	//     msg!("+ Processing PlaceBid");
	//     let accounts = parse_accounts(program_id, accounts)?;

	//     let auction_path = [
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         &args.resource.to_bytes(),
	//     ];
	//     assert_derivation(program_id, accounts.auction, &auction_path)?;

	//     // Load the auction and verify this bid is valid.
	//     let mut auction = AuctionData::from_account_info(accounts.auction)?;

	//     // Load the clock, used for various auction timing.
	//     let clock = Clock::from_account_info(accounts.clock_sysvar)?;

	//     // Verify auction has not ended.
	//     if auction.ended(clock.unix_timestamp)? {
	//         auction.state = auction.state.end()?;
	//         auction.serialize(&mut *accounts.auction.data.borrow_mut())?;
	//         msg!("Auction ended!");
	//         return Ok(());
	//     }
	//     // Derive Metadata key and load it.
	//     let metadata_bump = assert_derivation(
	//         program_id,
	//         accounts.bidder_meta,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             accounts.auction.key.as_ref(),
	//             accounts.bidder.key.as_ref(),
	//             "metadata".as_bytes(),
	//         ],
	//     )?;

	//     // If metadata doesn't exist, create it.
	//     if accounts.bidder_meta.owner != program_id {
	//         create_or_allocate_account_raw(
	//             *program_id,
	//             accounts.bidder_meta,
	//             accounts.rent,
	//             accounts.system,
	//             accounts.payer,
	//             // For whatever reason, using Mem function here returns 7, which is wholly wrong for this struct
	//             // seems to be issues with UnixTimestamp
	//             BIDDER_METADATA_LEN,
	//             &[
	//                 PREFIX.as_bytes(),
	//                 program_id.as_ref(),
	//                 accounts.auction.key.as_ref(),
	//                 accounts.bidder.key.as_ref(),
	//                 "metadata".as_bytes(),
	//                 &[metadata_bump],
	//             ],
	//         )?;
	//     } else {
	//         // Verify the last bid was cancelled before continuing.
	//         let bidder_metadata: BidderMetadata =
	//             BidderMetadata::from_account_info(accounts.bidder_meta)?;
	//         if bidder_metadata.cancelled == false {
	//             return Err(AuctionError::BidAlreadyActive.into());
	//         }
	//     };

	//     // Derive Pot address, this account wraps/holds an SPL account to transfer tokens into and is
	//     // also used as the authoriser of the SPL pot.
	//     let pot_bump = assert_derivation(
	//         program_id,
	//         accounts.bidder_pot,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             accounts.auction.key.as_ref(),
	//             accounts.bidder.key.as_ref(),
	//         ],
	//     )?;
	//     // The account within the pot must be new

	//     // Can't bid on an auction that isn't running.
	//     if auction.state != AuctionState::Started {
	//         return Err(AuctionError::InvalidState.into());
	//     }

	//     let bump_authority_seeds = &[
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         accounts.auction.key.as_ref(),
	//         accounts.bidder.key.as_ref(),
	//         &[pot_bump],
	//     ];

	//     // If the bidder pot account is empty, we need to generate one.
	//     if accounts.bidder_pot.data_is_empty() {
	//         create_or_allocate_account_raw(
	//             *program_id,
	//             accounts.bidder_pot,
	//             accounts.rent,
	//             accounts.system,
	//             accounts.payer,
	//             mem::size_of::<BidderPot>(),
	//             bump_authority_seeds,
	//         )?;

	//         // Attach SPL token address to pot account.
	//         let mut pot = BidderPot::from_account_info(accounts.bidder_pot)?;
	//         pot.bidder_pot = *accounts.bidder_pot_token.key;
	//         pot.bidder_act = *accounts.bidder.key;
	//         pot.auction_act = *accounts.auction.key;
	//         pot.serialize(&mut *accounts.bidder_pot.data.borrow_mut())?;

	//         assert_uninitialized::<Account>(accounts.bidder_pot_token)?;
	//         let bidder_token_account_bump = assert_derivation(
	//             program_id,
	//             accounts.bidder_pot_token,
	//             &[
	//                 PREFIX.as_bytes(),
	//                 &accounts.bidder_pot.key.as_ref(),
	//                 BIDDER_POT_TOKEN.as_bytes(),
	//             ],
	//         )?;
	//         let bidder_token_account_seeds = &[
	//             PREFIX.as_bytes(),
	//             &accounts.bidder_pot.key.as_ref(),
	//             BIDDER_POT_TOKEN.as_bytes(),
	//             &[bidder_token_account_bump],
	//         ];

	//         spl_token_create_account(TokenCreateAccount {
	//             payer: accounts.payer.clone(),
	//             authority: accounts.auction.clone(),
	//             authority_seeds: bidder_token_account_seeds,
	//             token_program: accounts.token_program.clone(),
	//             mint: accounts.mint.clone(),
	//             account: accounts.bidder_pot_token.clone(),
	//             system_program: accounts.system.clone(),
	//             rent: accounts.rent.clone(),
	//         })?;
	//     } else {
	//         // Already exists, verify that the pot contains the specified SPL address.
	//         let bidder_pot = BidderPot::from_account_info(accounts.bidder_pot)?;
	//         if bidder_pot.bidder_pot != *accounts.bidder_pot_token.key {
	//             return Err(AuctionError::BidderPotTokenAccountOwnerMismatch.into());
	//         }
	//         assert_initialized::<Account>(accounts.bidder_pot_token)?;
	//     }

	//     // Update now we have new bid.
	//     assert_derivation(
	//         program_id,
	//         accounts.auction_extended,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             args.resource.as_ref(),
	//             EXTENDED.as_bytes(),
	//         ],
	//     )?;
	//     let mut auction_extended: AuctionDataExtended =
	//         AuctionDataExtended::from_account_info(accounts.auction_extended)?;
	//     auction_extended.total_uncancelled_bids = auction_extended
	//         .total_uncancelled_bids
	//         .checked_add(1)
	//         .ok_or(AuctionError::NumericalOverflowError)?;
	//     auction_extended.serialize(&mut *accounts.auction_extended.data.borrow_mut())?;

	//     let mut bid_price = args.amount;

	//     if let Some(instant_sale_price) = auction_extended.instant_sale_price {
	//         if args.amount > instant_sale_price {
	//             msg!("Received amount is more than instant_sale_price so it was reduced to instant_sale_price - {:?}", instant_sale_price);
	//             bid_price = instant_sale_price;
	//         }
	//     }

	//     // Confirm payers SPL token balance is enough to pay the bid.
	//     let account: Account = Account::unpack_from_slice(&accounts.bidder_token.data.borrow())?;
	//     if account.amount.saturating_sub(bid_price) < 0 {
	//         msg!(
	//             "Amount is too small: {:?}, compared to account amount of {:?}",
	//             bid_price,
	//             account.amount
	//         );
	//         return Err(AuctionError::BalanceTooLow.into());
	//     }

	//     // Transfer amount of SPL token to bid account.
	//     spl_token_transfer(TokenTransferParams {
	//         source: accounts.bidder_token.clone(),
	//         destination: accounts.bidder_pot_token.clone(),
	//         authority: accounts.transfer_authority.clone(),
	//         authority_signer_seeds: bump_authority_seeds,
	//         token_program: accounts.token_program.clone(),
	//         amount: bid_price,
	//     })?;

	//     // Serialize new Auction State
	//     auction.last_bid = Some(clock.unix_timestamp);
	//     auction.place_bid(
	//         Bid(*accounts.bidder.key, bid_price),
	//         auction_extended.tick_size,
	//         auction_extended.gap_tick_size_percentage,
	//         clock.unix_timestamp,
	//         auction_extended.instant_sale_price,
	//     )?;
	//     auction.serialize(&mut *accounts.auction.data.borrow_mut())?;

	//     // Update latest metadata with results from the bid.
	//     BidderMetadata {
	//         bidder_pubkey: *accounts.bidder.key,
	//         auction_pubkey: *accounts.auction.key,
	//         last_bid: bid_price,
	//         last_bid_timestamp: clock.unix_timestamp,
	//         cancelled: false,
	//     }
	//     .serialize(&mut *accounts.bidder_meta.data.borrow_mut())?;

	//     Ok(())
	// }

	return &types.MsgPlaceBidResponse{}, nil
}

func (m msgServer) CancelBid(goCtx context.Context, msg *types.MsgCancelBid) (*types.MsgCancelBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// ctx.EventManager().EmitTypedEvent(&types.EventCancelBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// pub fn cancel_bid(
	//     program_id: &Pubkey,
	//     accounts: &[AccountInfo],
	//     args: CancelBidArgs,
	// ) -> ProgramResult {
	//     msg!("+ Processing Cancelbid");
	//     let accounts = parse_accounts(program_id, accounts)?;

	//     // The account within the pot must be owned by us.
	//     let actual_account: Account = assert_initialized(accounts.bidder_pot_token)?;
	//     if actual_account.owner != *accounts.auction.key {
	//         return Err(AuctionError::BidderPotTokenAccountOwnerMismatch.into());
	//     }

	//     // Derive and load Auction.
	//     let auction_bump = assert_derivation(
	//         program_id,
	//         accounts.auction,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             args.resource.as_ref(),
	//         ],
	//     )?;

	//     let auction_seeds = &[
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         args.resource.as_ref(),
	//         &[auction_bump],
	//     ];

	//     // Load the auction and verify this bid is valid.
	//     let mut auction = AuctionData::from_account_info(accounts.auction)?;
	//     // The mint provided in this bid must match the one the auction was initialized with.
	//     if auction.token_mint != *accounts.mint.key {
	//         return Err(AuctionError::IncorrectMint.into());
	//     }

	//     // Load auction extended account to check instant_sale_price
	//     // and update cancelled bids if auction still active
	//     let mut auction_extended = AuctionDataExtended::from_account_info(accounts.auction_extended)?;

	//     // Load the clock, used for various auction timing.
	//     let clock = Clock::from_account_info(accounts.clock_sysvar)?;

	//     // Derive Metadata key and load it.
	//     let metadata_bump = assert_derivation(
	//         program_id,
	//         accounts.bidder_meta,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             accounts.auction.key.as_ref(),
	//             accounts.bidder.key.as_ref(),
	//             "metadata".as_bytes(),
	//         ],
	//     )?;

	//     // If metadata doesn't exist, error, can't cancel a bid that doesn't exist and metadata must
	//     // exist if a bid was placed.
	//     if accounts.bidder_meta.owner != program_id {
	//         return Err(AuctionError::MetadataInvalid.into());
	//     }

	//     // Derive Pot address, this account wraps/holds an SPL account to transfer tokens out of.
	//     let pot_seeds = [
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         accounts.auction.key.as_ref(),
	//         accounts.bidder.key.as_ref(),
	//     ];

	//     let pot_bump = assert_derivation(program_id, accounts.bidder_pot, &pot_seeds)?;

	//     let bump_authority_seeds = &[
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         accounts.auction.key.as_ref(),
	//         accounts.bidder.key.as_ref(),
	//         &[pot_bump],
	//     ];

	//     // If the bidder pot account is empty, this bid is invalid.
	//     if accounts.bidder_pot.data_is_empty() {
	//         return Err(AuctionError::BidderPotDoesNotExist.into());
	//     }

	//     // Refuse to cancel if the auction ended and this person is a winning account.
	//     let winner_bid_index = auction.is_winner(accounts.bidder.key);
	//     if auction.ended(clock.unix_timestamp)? && winner_bid_index.is_some() {
	//         return Err(AuctionError::InvalidState.into());
	//     }

	//     // Refuse to cancel if bidder set price above or equal instant_sale_price
	//     if let Some(bid_index) = winner_bid_index {
	//         if let Some(instant_sale_price) = auction_extended.instant_sale_price {
	//             if auction.bid_state.amount(bid_index) >= instant_sale_price {
	//                 return Err(AuctionError::InvalidState.into());
	//             }
	//         }
	//     }

	//     // Confirm we're looking at the real SPL account for this bidder.
	//     let bidder_pot = BidderPot::from_account_info(accounts.bidder_pot)?;
	//     if bidder_pot.bidder_pot != *accounts.bidder_pot_token.key {
	//         return Err(AuctionError::BidderPotTokenAccountOwnerMismatch.into());
	//     }

	//     // Transfer SPL bid balance back to the user.
	//     let account: Account = Account::unpack_from_slice(&accounts.bidder_pot_token.data.borrow())?;
	//     spl_token_transfer(TokenTransferParams {
	//         source: accounts.bidder_pot_token.clone(),
	//         destination: accounts.bidder_token.clone(),
	//         authority: accounts.auction.clone(),
	//         authority_signer_seeds: auction_seeds,
	//         token_program: accounts.token_program.clone(),
	//         amount: account.amount,
	//     })?;

	//     // Update Metadata
	//     let metadata = BidderMetadata::from_account_info(accounts.bidder_meta)?;
	//     let already_cancelled = metadata.cancelled;
	//     BidderMetadata {
	//         cancelled: true,
	//         ..metadata
	//     }
	//     .serialize(&mut *accounts.bidder_meta.data.borrow_mut())?;

	//     // Update Auction

	//     if auction.state != AuctionState::Ended {
	//         // Once ended we want uncancelled bids to retain it's pre-ending count
	//         assert_derivation(
	//             program_id,
	//             accounts.auction_extended,
	//             &[
	//                 PREFIX.as_bytes(),
	//                 program_id.as_ref(),
	//                 args.resource.as_ref(),
	//                 EXTENDED.as_bytes(),
	//             ],
	//         )?;

	//         msg!("Already cancelled is {:?}", already_cancelled);

	//         if !already_cancelled && auction_extended.total_uncancelled_bids > 0 {
	//             auction_extended.total_uncancelled_bids = auction_extended
	//                 .total_uncancelled_bids
	//                 .checked_sub(1)
	//                 .ok_or(AuctionError::NumericalOverflowError)?;
	//         }
	//         auction_extended.serialize(&mut *accounts.auction_extended.data.borrow_mut())?;

	//         // Only cancel the bid if the auction has not ended yet
	//         auction.bid_state.cancel_bid(*accounts.bidder.key);
	//         auction.serialize(&mut *accounts.auction.data.borrow_mut())?;
	//     }

	//     Ok(())
	// }

	return &types.MsgCancelBidResponse{}, nil
}

func (m msgServer) ClaimBid(goCtx context.Context, msg *types.MsgClaimBid) (*types.MsgClaimBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	// ctx.EventManager().EmitTypedEvent(&types.EventClaimBid{
	// 	Creator:   msg.Sender,
	// 	AuctionId: metadata.Id,
	// })

	// pub fn claim_bid(
	//     program_id: &Pubkey,
	//     accounts: &[AccountInfo],
	//     args: ClaimBidArgs,
	// ) -> ProgramResult {
	//     msg!("+ Processing ClaimBid");
	//     let accounts = parse_accounts(program_id, accounts)?;
	//     let clock = Clock::from_account_info(accounts.clock_sysvar)?;

	//     // The account within the pot must be owned by us.
	//     let actual_account: Account = assert_initialized(accounts.bidder_pot_token)?;
	//     if actual_account.owner != *accounts.auction.key {
	//         return Err(AuctionError::BidderPotTokenAccountOwnerMismatch.into());
	//     }

	//     // Derive and load Auction.
	//     let auction_bump = assert_derivation(
	//         program_id,
	//         accounts.auction,
	//         &[
	//             PREFIX.as_bytes(),
	//             program_id.as_ref(),
	//             args.resource.as_ref(),
	//         ],
	//     )?;

	//     let auction_seeds = &[
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         args.resource.as_ref(),
	//         &[auction_bump],
	//     ];

	//     // Load the auction and verify this bid is valid.
	//     let auction = AuctionData::from_account_info(accounts.auction)?;

	//     if auction.authority != *accounts.authority.key {
	//         return Err(AuctionError::InvalidAuthority.into());
	//     }

	//     // User must have won the auction in order to claim their funds. Check early as the rest of the
	//     // checks will be for nothing otherwise.
	//     let bid_index = auction.is_winner(accounts.bidder.key);
	//     if bid_index.is_none() {
	//         msg!("User {:?} is not winner", accounts.bidder.key);
	//         return Err(AuctionError::InvalidState.into());
	//     }

	//     let instant_sale_price = accounts.auction_extended.and_then(|info| {
	//         assert_derivation(
	//             program_id,
	//             info,
	//             &[
	//                 PREFIX.as_bytes(),
	//                 program_id.as_ref(),
	//                 args.resource.as_ref(),
	//                 EXTENDED.as_bytes(),
	//             ],
	//         )
	//         .ok()?;

	//         AuctionDataExtended::from_account_info(info)
	//             .ok()?
	//             .instant_sale_price
	//     });

	//     // Auction either must have ended or bidder pay instant_sale_price
	//     if !auction.ended(clock.unix_timestamp)? {
	//         match instant_sale_price {
	//             Some(instant_sale_price)
	//                 if auction.bid_state.amount(bid_index.unwrap()) < instant_sale_price =>
	//             {
	//                 return Err(AuctionError::InvalidState.into())
	//             }
	//             None => return Err(AuctionError::InvalidState.into()),
	//             _ => (),
	//         }
	//     }

	//     // The mint provided in this claim must match the one the auction was initialized with.
	//     if auction.token_mint != *accounts.mint.key {
	//         return Err(AuctionError::IncorrectMint.into());
	//     }

	//     // Derive Pot address, this account wraps/holds an SPL account to transfer tokens into.
	//     let pot_seeds = [
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         accounts.auction.key.as_ref(),
	//         accounts.bidder.key.as_ref(),
	//     ];

	//     let pot_bump = assert_derivation(program_id, accounts.bidder_pot, &pot_seeds)?;

	//     let bump_authority_seeds = &[
	//         PREFIX.as_bytes(),
	//         program_id.as_ref(),
	//         accounts.auction.key.as_ref(),
	//         accounts.bidder.key.as_ref(),
	//         &[pot_bump],
	//     ];

	//     // If the bidder pot account is empty, this bid is invalid.
	//     if accounts.bidder_pot.data_is_empty() {
	//         return Err(AuctionError::BidderPotDoesNotExist.into());
	//     }

	//     // Confirm we're looking at the real SPL account for this bidder.
	//     let mut bidder_pot = BidderPot::from_account_info(accounts.bidder_pot)?;
	//     if bidder_pot.bidder_pot != *accounts.bidder_pot_token.key {
	//         return Err(AuctionError::BidderPotTokenAccountOwnerMismatch.into());
	//     }

	//     // Transfer SPL bid balance back to the user.
	//     spl_token_transfer(TokenTransferParams {
	//         source: accounts.bidder_pot_token.clone(),
	//         destination: accounts.destination.clone(),
	//         authority: accounts.auction.clone(),
	//         authority_signer_seeds: auction_seeds,
	//         token_program: accounts.token_program.clone(),
	//         amount: actual_account.amount,
	//     })?;

	//     bidder_pot.emptied = true;
	//     bidder_pot.serialize(&mut *accounts.bidder_pot.data.borrow_mut())?;

	//     Ok(())
	// }

	return &types.MsgClaimBidResponse{}, nil
}
