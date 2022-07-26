package types

func (a Auction) IsActive() bool {
	return a.State == AuctionState_Started
}
