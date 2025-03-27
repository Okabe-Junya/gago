package ga

import (
	"testing"
	"time"
)

func TestGenerationCountTermination(t *testing.T) {
	maxGenerations := 5
	termCondition := GenerationCountTermination(maxGenerations)

	// Setup mock GA instance
	ga := &GA{
		History: make([]*Statistics, 0),
	}

	// Add history entries (simulating generations)
	for i := 0; i < maxGenerations-1; i++ {
		ga.History = append(ga.History, &Statistics{})
		if termCondition.Evaluate(ga) {
			t.Errorf("Termination condition returned true after %d generations, expected false until %d generations", i+1, maxGenerations)
		}
	}

	// Add one more to reach maxGenerations
	ga.History = append(ga.History, &Statistics{})
	if !termCondition.Evaluate(ga) {
		t.Errorf("Termination condition returned false after %d generations, expected true", maxGenerations)
	}
}

func TestConvergenceTermination(t *testing.T) {
	noImprovementGens := 3
	improvementThreshold := 0.01
	termCondition := ConvergenceTermination(noImprovementGens, improvementThreshold)

	// Setup mock GA instance
	ga := &GA{
		History: make([]*Statistics, 0),
	}

	// Add history entries with improving fitness
	fitnessValues := []float64{1.0, 1.05, 1.1, 1.15, 1.2, 1.25, 1.3, 1.305}
	for _, fitness := range fitnessValues {
		ga.History = append(ga.History, &Statistics{BestFitness: fitness})

		// Should only return true when we've had noImprovementGens with less than improvementThreshold improvement
		shouldTerminate := len(ga.History) > noImprovementGens &&
			(ga.History[len(ga.History)-1].BestFitness-ga.History[len(ga.History)-1-noImprovementGens].BestFitness) < improvementThreshold

		if termCondition.Evaluate(ga) != shouldTerminate {
			t.Errorf(
				"Termination condition incorrect at fitness %f (history length: %d); got %v, expected %v",
				fitness, len(ga.History), termCondition.Evaluate(ga), shouldTerminate,
			)
		}
	}
}

func TestFitnessThresholdTermination(t *testing.T) {
	threshold := 10.0
	termCondition := FitnessThresholdTermination(threshold)

	// Setup mock GA instance with population
	ga := &GA{
		Population: &Population{
			Statistics: &Statistics{
				BestFitness: 5.0, // Initial fitness below threshold
			},
		},
	}

	// Test below threshold
	if termCondition.Evaluate(ga) {
		t.Errorf("Termination condition returned true for fitness %.2f, expected false for threshold %.2f",
			ga.Population.Statistics.BestFitness, threshold)
	}

	// Test at threshold
	ga.Population.Statistics.BestFitness = threshold
	if !termCondition.Evaluate(ga) {
		t.Errorf("Termination condition returned false for fitness %.2f, expected true for threshold %.2f",
			ga.Population.Statistics.BestFitness, threshold)
	}

	// Test above threshold
	ga.Population.Statistics.BestFitness = 15.0
	if !termCondition.Evaluate(ga) {
		t.Errorf("Termination condition returned false for fitness %.2f, expected true for threshold %.2f",
			ga.Population.Statistics.BestFitness, threshold)
	}
}

func TestTimeBasedTermination(t *testing.T) {
	// Using a test interface that doesn't depend on actual time
	// We use mockTimeNow to control time within the test
	originalTimeFunc := timeNow
	defer func() { timeNow = originalTimeFunc }() // Restore original function after test

	mockTime := time.Now()
	timeNow = func() time.Time { return mockTime }

	// Set a short duration for testing
	duration := 100 * time.Millisecond
	termCondition := TimeBasedTermination(duration)

	// Setup mock GA instance with start time
	ga := &GA{
		StartTime: mockTime,
	}

	// Test before duration has elapsed
	if termCondition.Evaluate(ga) {
		t.Error("Termination condition returned true immediately, expected false")
	}

	// Advance time (without sleep)
	mockTime = mockTime.Add(duration + 10*time.Millisecond)

	// Test after duration has elapsed
	if !termCondition.Evaluate(ga) {
		t.Error("Termination condition returned false after time advance, expected true")
	}
}
