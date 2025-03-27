package termination

import (
	"testing"
	"time"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

func TestDiversityThresholdTermination(t *testing.T) {
	// Setup a mock GA instance
	mockGA := &ga.GA{
		Population: &ga.Population{
			Statistics: &ga.Statistics{
				Diversity: 0.5, // Initial diversity at 50%
			},
		},
	}

	t.Run("terminates below threshold", func(t *testing.T) {
		// Diversity is 0.5, threshold is 0.6
		termCondition := DiversityThresholdTermination(0.6)
		if !termCondition.Evaluate(mockGA) {
			t.Error("Should terminate when diversity (0.5) is below threshold (0.6)")
		}
	})

	t.Run("continues above threshold", func(t *testing.T) {
		// Diversity is 0.5, threshold is 0.4
		termCondition := DiversityThresholdTermination(0.4)
		if termCondition.Evaluate(mockGA) {
			t.Error("Should continue when diversity (0.5) is above threshold (0.4)")
		}
	})

	t.Run("handles nil GA", func(t *testing.T) {
		termCondition := DiversityThresholdTermination(0.5)
		if termCondition.Evaluate(nil) {
			t.Error("Should return false when GA is nil")
		}
	})

	t.Run("normalizes threshold values", func(t *testing.T) {
		// Test with out of bounds thresholds
		// For negative threshold, we expect diversity (0.5) to be greater than 0.0
		negativeThreshold := DiversityThresholdTermination(-0.1)
		if negativeThreshold.Evaluate(mockGA) {
			t.Error("Normalized negative threshold (0.0) should not terminate when diversity is 0.5")
		}

		// For high threshold, we expect diversity (0.5) to be less than 1.0
		highThreshold := DiversityThresholdTermination(1.5)
		if !highThreshold.Evaluate(mockGA) {
			t.Error("Normalized high threshold (1.0) should terminate when diversity is 0.5")
		}
	})
}

func TestDiversityStagnationTermination(t *testing.T) {
	// Setup a mock GA instance
	mockGA := &ga.GA{
		Population: &ga.Population{
			Statistics: &ga.Statistics{
				Diversity: 0.5, // Initial diversity
			},
		},
	}

	t.Run("terminates after stagnation", func(t *testing.T) {
		termCondition := DiversityStagnationTermination(3)

		// First call, diversity doesn't change
		if termCondition.Evaluate(mockGA) {
			t.Error("Should not terminate on first generation")
		}

		// Change diversity to improve, then stagnate
		mockGA.Population.Statistics.Diversity = 0.6
		if termCondition.Evaluate(mockGA) {
			t.Error("Should not terminate when diversity improves")
		}

		// Stagnant for 1 generation
		if termCondition.Evaluate(mockGA) {
			t.Error("Should not terminate after 1 generation of stagnation")
		}

		// Stagnant for 2 generations
		if termCondition.Evaluate(mockGA) {
			t.Error("Should not terminate after 2 generations of stagnation")
		}

		// Stagnant for 3 generations (should terminate)
		if !termCondition.Evaluate(mockGA) {
			t.Error("Should terminate after 3 generations of stagnation")
		}
	})

	t.Run("resets counter on improvement", func(t *testing.T) {
		termCondition := DiversityStagnationTermination(2)

		// First call
		termCondition.Evaluate(mockGA)

		// Stagnant for 1 generation
		termCondition.Evaluate(mockGA)

		// Diversity improves
		mockGA.Population.Statistics.Diversity = 0.7
		if termCondition.Evaluate(mockGA) {
			t.Error("Should reset counter when diversity improves")
		}

		// Stagnant for 1 generation again
		if termCondition.Evaluate(mockGA) {
			t.Error("Should not terminate after 1 generation of stagnation")
		}

		// Stagnant for 2 generations (should terminate)
		if !termCondition.Evaluate(mockGA) {
			t.Error("Should terminate after 2 generations of stagnation")
		}
	})

	t.Run("handles invalid generations parameter", func(t *testing.T) {
		termCondition := DiversityStagnationTermination(-1)

		// Should normalize to 1
		termCondition.Evaluate(mockGA)

		// Should terminate after 1 generation of stagnation
		if !termCondition.Evaluate(mockGA) {
			t.Error("Should normalize negative generations to 1")
		}
	})
}

func TestCompositeTermination(t *testing.T) {
	// Setup a mock GA instance
	mockGA := &ga.GA{
		StartTime: time.Now().Add(-1 * time.Minute),
		Population: &ga.Population{
			Statistics: &ga.Statistics{
				Diversity:   0.5,
				BestFitness: 0.7,
			},
		},
	}

	t.Run("any operator (OR logic)", func(t *testing.T) {
		// Create conditions: one that passes, one that fails
		condition1 := DiversityThresholdTermination(0.6) // Passes (diversity 0.5 < threshold 0.6)
		condition2 := FitnessThresholdTermination(0.8)   // Fails (fitness 0.7 < threshold 0.8)

		composite := NewCompositeTermination(AnyTermination, condition1, condition2)

		if !composite.Evaluate(mockGA) {
			t.Error("OR composite should terminate when at least one condition is met")
		}
	})

	t.Run("all operator (AND logic)", func(t *testing.T) {
		// Create conditions: one that passes, one that fails
		condition1 := DiversityThresholdTermination(0.6) // Passes (diversity 0.5 < threshold 0.6)
		condition2 := FitnessThresholdTermination(0.8)   // Fails (fitness 0.7 < threshold 0.8)

		composite := NewCompositeTermination(AllTermination, condition1, condition2)

		if composite.Evaluate(mockGA) {
			t.Error("AND composite should not terminate unless all conditions are met")
		}

		// Now make both conditions pass
		condition3 := FitnessThresholdTermination(0.6) // Passes (fitness 0.7 > threshold 0.6)

		composite = NewCompositeTermination(AllTermination, condition1, condition3)

		if !composite.Evaluate(mockGA) {
			t.Error("AND composite should terminate when all conditions are met")
		}
	})

	t.Run("handles nil conditions", func(t *testing.T) {
		// Create a composite with nil conditions
		composite := NewCompositeTermination(AnyTermination, nil)

		if composite.Evaluate(mockGA) {
			t.Error("Composite with nil conditions should not terminate")
		}
	})

	t.Run("add and remove conditions", func(t *testing.T) {
		// Create an empty composite
		composite := &CompositeTermination{
			operator:   AnyTermination,
			conditions: []ga.TerminationCondition{},
		}

		// Add a condition that would pass
		condition := DiversityThresholdTermination(0.6)
		composite.AddCondition(condition)

		if !composite.Evaluate(mockGA) {
			t.Error("Should terminate after adding a passing condition")
		}

		// Remove the condition
		composite.RemoveCondition(0)

		if composite.Evaluate(mockGA) {
			t.Error("Should not terminate after removing the condition")
		}
	})
}
