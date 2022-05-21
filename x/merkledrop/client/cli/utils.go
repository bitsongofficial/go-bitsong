package cli

import (
	"errors"
	"strconv"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	var startTime time.Time
	if timeStr == "" { // empty start time
		startTime = time.Unix(0, 0)
	} else if timeUnix, err := strconv.ParseInt(timeStr, 10, 64); err == nil { // unix time
		startTime = time.Unix(timeUnix, 0)
	} else if timeRFC, err := time.Parse(time.RFC3339, timeStr); err == nil { // RFC time
		startTime = timeRFC
	} else { // invalid input
		return startTime, errors.New("invalid start time format")
	}

	return startTime, nil
}
