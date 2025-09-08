package rules

import (
	"fmt"
	"time"
)

type EndTimeRule struct {
	EndTime time.Time
}

func (r *EndTimeRule) Name() string {
	return "end_time"
}

func (r *EndTimeRule) Validate(req *RuleEngineRequest) error {
	if req.currentTime.After(r.EndTime) {
		return fmt.Errorf("current time %s is after end time %s", req.currentTime, r.EndTime)
	}

	return nil
}
