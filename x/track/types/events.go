package types

// Track module event types
const (
	EventTypeCreateTrack    = "create_track"
	EventTypeSetTrackStatus = "set_track_status"

	// Track attributes
	AttributeValueCategory  = ModuleName
	AttributeKeyTrackID     = ModuleName + "_id"
	AttributeKeyTrackTitle  = ModuleName + "_title"
	AttributeKeyTrackOwner  = ModuleName + "_owner"
	AttributeKeyTrackStatus = ModuleName + "_status"
)
