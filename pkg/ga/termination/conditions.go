package termination

import (
	"math"
	"time"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

// GenerationCountTermination returns a termination condition that terminates after a specified number of generations.
//
// Parameters:
//   - maxGenerations: the maximum number of generations after which evolution should terminate.
//
// Returns:
//   - A TerminationCondition that evaluates to true when the specified number of generations is reached.
func GenerationCountTermination(maxGenerations int) ga.TerminationCondition {
	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		return len(ga.History) >= maxGenerations
	})
}

// ConvergenceTermination returns a termination condition that terminates when
// the best fitness hasn't improved by the specified threshold over the specified number of generations.
//
// Parameters:
//   - noImprovementGens: the number of consecutive generations without significant improvement.
//   - improvementThreshold: the minimum improvement considered significant.
//
// Returns:
//   - A TerminationCondition that evaluates to true when fitness improvement stagnates.
func ConvergenceTermination(noImprovementGens int, improvementThreshold float64) ga.TerminationCondition {
	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		if len(ga.History) <= noImprovementGens {
			return false
		}

		currentBest := ga.History[len(ga.History)-1].BestFitness
		pastBest := ga.History[len(ga.History)-1-noImprovementGens].BestFitness
		improvement := math.Abs(currentBest - pastBest)

		return improvement < improvementThreshold
	})
}

// TimeBasedTermination returns a termination condition that terminates after a specified duration.
//
// Parameters:
//   - duration: the maximum runtime after which evolution should terminate.
//
// Returns:
//   - A TerminationCondition that evaluates to true when the specified duration is exceeded.
func TimeBasedTermination(duration time.Duration) ga.TerminationCondition {
	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		return time.Since(ga.StartTime) >= duration
	})
}

// FitnessThresholdTermination returns a termination condition that terminates when
// the best fitness reaches or exceeds the specified threshold.
//
// Parameters:
//   - threshold: the fitness threshold above which evolution should terminate.
//
// Returns:
//   - A TerminationCondition that evaluates to true when the specified fitness threshold is reached.
func FitnessThresholdTermination(threshold float64) ga.TerminationCondition {
	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		return ga.Population.Statistics.BestFitness >= threshold
	})
}

// FitnessStagnationTermination returns a termination condition that stops evolution
// when the best fitness has not improved for a specified number of generations.
func FitnessStagnationTermination(generations int) ga.TerminationCondition {
	if generations < 1 {
		generations = 1
	}

	bestFitness := -1.0
	stagnationCount := 0

	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		currentFitness := ga.Population.Statistics.BestFitness

		if currentFitness > bestFitness {
			bestFitness = currentFitness
			stagnationCount = 0
		} else {
			stagnationCount++
		}

		return stagnationCount >= generations
	})
}

// FitnessImprovementTermination returns a termination condition that stops evolution
// when the rate of fitness improvement falls below a threshold.
func FitnessImprovementTermination(threshold float64) ga.TerminationCondition {
	prevFitness := -1.0

	return ga.TerminationConditionFunc(func(ga *ga.GA) bool {
		currentFitness := ga.Population.Statistics.BestFitness
		if prevFitness < 0 {
			prevFitness = currentFitness
			return false
		}

		improvement := (currentFitness - prevFitness) / prevFitness
		prevFitness = currentFitness

		return improvement < threshold
	})
}
