package termination

import (
	"github.com/Okabe-Junya/gago/pkg/ga"
)

// DiversityThresholdTermination returns a termination condition that stops evolution
// when the population diversity falls below a specified threshold.
//
// This is useful for problems where maintaining population diversity is important,
// and termination should occur when the population becomes too homogeneous.
//
// Parameters:
//   - threshold: The diversity threshold below which evolution should stop (0.0 to 1.0).
//
// Returns:
//   - A TerminationCondition that evaluates to true when diversity falls below the threshold.
func DiversityThresholdTermination(threshold float64) ga.TerminationCondition {
	if threshold < 0.0 {
		threshold = 0.0
	} else if threshold > 1.0 {
		threshold = 1.0
	}

	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		if ga == nil || ga.Population == nil || ga.Population.Statistics == nil {
			return false
		}
		return ga.Population.Statistics.Diversity < threshold
	})
}

// DiversityStagnationTermination returns a termination condition that stops evolution
// when the population diversity has not improved for a specified number of generations.
//
// This is useful for detecting when diversity evolution has plateaued, indicating
// that the algorithm may be stuck in a local optimum.
//
// Parameters:
//   - generations: The number of consecutive generations without diversity improvement
//     after which evolution should stop.
//
// Returns:
//   - A TerminationCondition that evaluates to true when diversity has stagnated
//     for the specified number of generations.
func DiversityStagnationTermination(generations int) ga.TerminationCondition {
	if generations < 1 {
		generations = 1
	}

	bestDiversity := 0.0
	stagnationCount := 0

	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		if ga == nil || ga.Population == nil || ga.Population.Statistics == nil {
			return false
		}

		currentDiversity := ga.Population.Statistics.Diversity

		if currentDiversity > bestDiversity {
			bestDiversity = currentDiversity
			stagnationCount = 0
		} else {
			stagnationCount++
		}

		return stagnationCount >= generations
	})
}

// DiversityImprovementTermination returns a termination condition that stops evolution
// when the rate of diversity improvement falls below a specified threshold.
//
// This is useful for detecting when diversity improvements are becoming marginal,
// indicating diminishing returns from further evolution.
//
// Parameters:
//   - threshold: The minimum acceptable improvement rate (can be negative for diversity decrease).
//
// Returns:
//   - A TerminationCondition that evaluates to true when the diversity improvement
//     rate falls below the threshold.
func DiversityImprovementTermination(threshold float64) ga.TerminationCondition {
	prevDiversity := -1.0

	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		if ga == nil || ga.Population == nil || ga.Population.Statistics == nil {
			return false
		}

		currentDiversity := ga.Population.Statistics.Diversity

		// Initialize on first call
		if prevDiversity < 0 {
			prevDiversity = currentDiversity
			return false
		}

		// Avoid division by zero
		if prevDiversity == 0 {
			// If previous diversity was 0, any positive value is an improvement
			improvement := currentDiversity > 0
			prevDiversity = currentDiversity
			return !improvement
		}

		improvement := (currentDiversity - prevDiversity) / prevDiversity
		prevDiversity = currentDiversity

		return improvement < threshold
	})
}
