package rules

import (
	"fmt"
	"time"

	"github.com/bitsongofficial/go-bitsong/x/drop/types"
)

type RuleEngineRequest struct {
	currentTime time.Time
}

type RuleI interface {
	Name() string
	Validate(req *RuleEngineRequest) error
}

type RuleEngine struct {
	rules map[string]RuleI
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: make(map[string]RuleI),
	}
}

func (e *RuleEngine) Register(rule RuleI) error {
	name := rule.Name()
	if _, exists := e.rules[name]; exists {
		return nil
	}

	e.rules[name] = rule
	return nil
}

func (e *RuleEngine) Validate(req *RuleEngineRequest) error {
	for name, rule := range e.rules {
		if err := rule.Validate(req); err != nil {
			return fmt.Errorf("rule %s validation failed: %w", name, err)
		}
	}

	return nil
}

func NewRuleEngineFromRule(rule types.Rule) (*RuleEngine, error) {
	engine := NewRuleEngine()

	if !rule.StartTime.IsZero() {
		err := engine.Register(&StartTimeRule{
			StartTime: rule.StartTime,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register start time rule: %w", err)
		}
	}

	if !rule.EndTime.IsZero() {
		err := engine.Register(&EndTimeRule{
			EndTime: rule.EndTime,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register end time rule: %w", err)
		}
	} else {
		err := engine.Register(&EndTimeRule{
			EndTime: rule.StartTime.Add(7 * 24 * time.Hour),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to register default end time rule: %w", err)
		}
	}

	return engine, nil
}
