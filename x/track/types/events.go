package types

// Track module event types
const (
	EventTypeCreateTrack    = "create_track"
	EventTypePlayTrack      = "play_track"
	EventTypeSetTrackStatus = "set_track_status"
	EventTypeDepositTrack   = "deposit_track"

	// Track attributes
	AttributeValueCategory  = ModuleName
	AttributeKeyTrackID     = "id"
	AttributeKeyTrackTitle  = "title"
	AttributeKeyTrackOwner  = "owner"
	AttributeKeyTrackStatus = "status"

	AttributeKeyPlayAccAddr = "play_acc_addr"
)
