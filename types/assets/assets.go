package assets

//nolint
const (
	MicroBitSongDenom = "ubtsg"

	MicroUnit = int64(1e8)
)

// IsValidDenom returns the given denom is valid or not
func IsValidDenom(denom string) bool {
	return denom == MicroBitSongDenom
}