package client

//NormalizeArtistStatus - normalize user specified artist status
func NormalizeArtistStatus(status string) string {
	switch status {
	case "Verified", "verified":
		return "Verified"
	case "Rejected", "rejected":
		return "Rejected"
	case "Failed", "failed":
		return "Failed"
	}
	return ""
}
