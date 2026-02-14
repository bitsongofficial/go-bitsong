package rules

import (
	"fmt"
	"time"
)

type StartTimeRule struct {
	StartTime time.Time
}

func (r *StartTimeRule) Name() string {
	return "start_time"
}

func (r *StartTimeRule) Validate(req *RuleEngineRequest) error {
	if req.currentTime.Before(r.StartTime) {
		return fmt.Errorf("current time %s is before start time %s", req.currentTime, r.StartTime)
	}

	return nil
}
