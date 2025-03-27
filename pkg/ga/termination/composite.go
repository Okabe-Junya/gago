package termination

import (
	"github.com/Okabe-Junya/gago/pkg/ga"
)

// CompositeOperator defines how multiple termination conditions should be combined.
type CompositeOperator int

const (
	// AnyTermination stops evolution if any condition is met (logical OR).
	AnyTermination CompositeOperator = iota

	// AllTermination stops evolution only if all conditions are met (logical AND).
	AllTermination
)

// CompositeTermination represents a combination of multiple termination conditions.
// It allows combining simple conditions into more complex termination logic.
type CompositeTermination struct {
	conditions []ga.TerminationCondition
	operator   CompositeOperator
}

// NewCompositeTermination creates a new composite termination condition.
//
// Parameters:
//   - operator: Determines how the conditions are combined (AnyTermination or AllTermination).
//   - conditions: A variable number of termination conditions to combine.
//
// Returns:
//   - A TerminationCondition that combines the provided conditions according to the operator.
func NewCompositeTermination(operator CompositeOperator, conditions ...ga.TerminationCondition) ga.TerminationCondition {
	// Filter out nil conditions
	validConditions := make([]ga.TerminationCondition, 0, len(conditions))
	for _, cond := range conditions {
		if cond != nil {
			validConditions = append(validConditions, cond)
		}
	}

	return &CompositeTermination{
		conditions: validConditions,
		operator:   operator,
	}
}

// Evaluate implements the TerminationCondition interface.
// It evaluates all the contained conditions according to the composite operator.
func (ct *CompositeTermination) Evaluate(ga *ga.GA) bool {
	if ct == nil || len(ct.conditions) == 0 || ga == nil {
		return false
	}

	switch ct.operator {
	case AnyTermination:
		// Return true if any condition is met (logical OR)
		for _, condition := range ct.conditions {
			if condition.Evaluate(ga) {
				return true
			}
		}
		return false

	case AllTermination:
		// Return true only if all conditions are met (logical AND)
		for _, condition := range ct.conditions {
			if !condition.Evaluate(ga) {
				return false
			}
		}
		return true

	default:
		return false
	}
}

// AddCondition adds a new termination condition to the composite.
//
// Parameters:
//   - condition: The termination condition to add.
func (ct *CompositeTermination) AddCondition(condition ga.TerminationCondition) {
	if ct == nil || condition == nil {
		return
	}
	ct.conditions = append(ct.conditions, condition)
}

// RemoveCondition removes a termination condition from the composite at the specified index.
//
// Parameters:
//   - index: The index of the condition to remove.
func (ct *CompositeTermination) RemoveCondition(index int) {
	if ct == nil || index < 0 || index >= len(ct.conditions) {
		return
	}
	ct.conditions = append(ct.conditions[:index], ct.conditions[index+1:]...)
}

// SetOperator changes the composite operator.
//
// Parameters:
//   - operator: The new operator to use (AnyTermination or AllTermination).
func (ct *CompositeTermination) SetOperator(operator CompositeOperator) {
	if ct == nil {
		return
	}
	ct.operator = operator
}

// GetConditions returns all termination conditions.
//
// Returns:
//   - A slice containing all the termination conditions.
func (ct *CompositeTermination) GetConditions() []ga.TerminationCondition {
	if ct == nil {
		return nil
	}
	return ct.conditions
}

// GetOperator returns the current composite operator.
//
// Returns:
//   - The current composite operator.
func (ct *CompositeTermination) GetOperator() CompositeOperator {
	if ct == nil {
		return AnyTermination // Default
	}
	return ct.operator
}
