// Package ga: this file defines the Result returned by Evolve, the StopReason
// enum, and the EarlyStopping configuration.
package ga

// StopReason records why an evolution run exited.
type StopReason string

const (
	// StopMaxGenerations means the loop ran for ga.Generations and exited normally.
	StopMaxGenerations StopReason = "max_generations"
	// StopTermCondition means GA.TermCondition.Evaluate returned true.
	StopTermCondition StopReason = "term_condition"
	// StopPatience means EarlyStopping.Patience generations passed without improvement > Tol.
	StopPatience StopReason = "patience"
	// StopTargetFitness means the best fitness reached EarlyStopping.TargetFitness.
	StopTargetFitness StopReason = "target_fitness"
	// StopTimeLimit means wall-clock time exceeded EarlyStopping.TimeLimit.
	StopTimeLimit StopReason = "time_limit"
	// StopError means evolution terminated due to an unrecoverable error;
	// in that case Evolve returns a non-nil error and Result.Best may be nil.
	StopError StopReason = "error"
)

// Result is the outcome of a call to GA.Evolve.
//
// Best is the highest-fitness individual ever observed across all generations
// (not necessarily the best in the final generation).
// StoppedAtGeneration is the index of the last generation that ran; 0 means
// only the initial population was evaluated.
type Result struct {
	Best                *Individual
	Population          *Population
	History             []*Statistics
	StopReason          StopReason
	StoppedAtGeneration int
}

// EarlyStopping configures additional stop criteria evaluated each generation.
//
// Any combination of the fields may be set; whichever criterion fires first
// wins and is reported on Result.StopReason. A field is "set" when it differs
// from its zero value, except Tol which only matters when Patience > 0.
//
// Patience: stop when the best fitness has not improved by more than Tol for
// Patience consecutive generations. Set to 0 to disable.
// Tol: improvement threshold for Patience. Defaults to 0 (any improvement counts).
// TargetFitness: stop as soon as the current-generation best fitness is
// >= TargetFitness. Set TargetFitnessSet=true to enable (so that 0.0 is a
// legitimate target).
// TimeLimit: maximum wall-clock duration (in nanoseconds; use time.Duration).
// 0 disables.
type EarlyStopping struct {
	Patience         int
	Tol              float64
	TargetFitness    float64
	TargetFitnessSet bool
	TimeLimit        int64 // nanoseconds; use int64(time.Duration)
}
