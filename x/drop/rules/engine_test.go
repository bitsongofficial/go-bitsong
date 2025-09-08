package rules

import (
	"testing"
	"time"

	"github.com/bitsongofficial/go-bitsong/x/drop/types"
)

func TestRulesEngine(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)

	rules := types.Rule{
		StartTime: startTime,
		EndTime:   endTime,
	}

	engine, err := NewRuleEngineFromRule(rules)
	if err != nil {
		t.Fatalf("failed to create rules engine: %v", err)
	}

	req := &RuleEngineRequest{
		currentTime: time.Now(),
	}

	err = engine.Validate(req)
	if err != nil {
		t.Fatalf("rules validation failed: %v", err)
	}
}
