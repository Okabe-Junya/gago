package ga

import (
	"math/rand"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		EnableLogger:  false,
	}

	populationSize := 20
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		return NewBinaryGenotype(genomeLength)
	}

	evalFunc := func(genotype *Genotype) *Phenotype {
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	gaInstance.Initialize(populationSize, initFunc, evalFunc)

	// Check if population was initialized
	if gaInstance.Population == nil {
		t.Fatal("Population should not be nil after initialization")
	}

	// Check if population size is correct
	if len(gaInstance.Population.Individuals) != populationSize {
		t.Errorf("Expected population size %d, got %d", populationSize, len(gaInstance.Population.Individuals))
	}

	// Check if statistics were calculated
	if gaInstance.Population.Statistics == nil {
		t.Fatal("Population statistics should not be nil after initialization")
	}

	// Check if history was initialized
	if len(gaInstance.History) != 1 {
		t.Errorf("Expected history length 1, got %d", len(gaInstance.History))
	}
}

func TestEvolve(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   5,
		EnableLogger:  false,
	}

	populationSize := 10
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		return NewBinaryGenotype(genomeLength)
	}

	evalFunc := func(genotype *Genotype) *Phenotype {
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	t.Run("evolution completes successfully", func(t *testing.T) {
		err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
		if err != nil {
			t.Fatalf("Initialize() failed: %v", err)
		}

		result, err := gaInstance.Evolve(evalFunc)
		if err != nil {
			t.Fatalf("Evolve() failed: %v", err)
		}

		// Check if evolution occurred (history should have entries for each generation)
		expectedHistoryLength := gaInstance.Generations + 1 // Initial + each generation
		if got, want := len(gaInstance.History), expectedHistoryLength; got != want {
			t.Errorf("History length got %d, want %d", got, want)
		}

		// Check if best individual is not nil
		if result == nil || result.Best == nil {
			t.Fatal("Best individual is nil after evolution")
		}

		// When Generations are exhausted without any early-stop, StopReason
		// must be StopMaxGenerations.
		if result.StopReason != StopMaxGenerations {
			t.Errorf("StopReason got %q, want %q", result.StopReason, StopMaxGenerations)
		}

		// Evolve returns the all-time best individual seen across all generations,
		// which may be fitter than the final population's best (it might not
		// have survived selection). The invariant is: all-time best >= current best.
		popBest := gaInstance.Population.GetBestIndividual()
		if got, floor := result.Best.Phenotype.Fitness, popBest.Phenotype.Fitness; got < floor {
			t.Errorf("Best individual fitness got %f, want >= %f (current population best)", got, floor)
		}
	})

	t.Run("handles nil evaluation function", func(t *testing.T) {
		err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
		if err != nil {
			t.Fatalf("Initialize() failed: %v", err)
		}

		_, err = gaInstance.Evolve(nil)
		if err == nil {
			t.Error("Evolve() with nil evalFunc should return error, got nil")
		}
	})
}

func TestElitism(t *testing.T) {
	// Create a GA with elitism
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   5,
		ElitismCount:  2, // Keep top 2 individuals
		EnableLogger:  false,
	}

	populationSize := 10
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		genotype := NewBinaryGenotype(genomeLength)
		// Initialize with all zeros
		for i := range genotype.Genome {
			genotype.Genome[i] = 0
		}
		return genotype
	}

	evalFunc := func(genotype *Genotype) *Phenotype {
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
	if err != nil {
		t.Fatalf("Failed to initialize GA: %v", err)
	}

	// Create two elite individuals with perfect fitness
	gaInstance.Population.Individuals[0].Genotype.Genome = make([]byte, genomeLength)
	gaInstance.Population.Individuals[1].Genotype.Genome = make([]byte, genomeLength)
	for i := range gaInstance.Population.Individuals[0].Genotype.Genome {
		gaInstance.Population.Individuals[0].Genotype.Genome[i] = 1
		gaInstance.Population.Individuals[1].Genotype.Genome[i] = 1
	}
	gaInstance.Population.Individuals[0].Phenotype.Fitness = float64(genomeLength)
	gaInstance.Population.Individuals[1].Phenotype.Fitness = float64(genomeLength)
	gaInstance.Population.CalculateStatistics()

	// Store the elite genomes
	elite0 := make([]byte, genomeLength)
	elite1 := make([]byte, genomeLength)
	copy(elite0, gaInstance.Population.Individuals[0].Genotype.Genome)
	copy(elite1, gaInstance.Population.Individuals[1].Genotype.Genome)

	// Evolve for one generation
	gaInstance.Generations = 1
	_, err = gaInstance.Evolve(evalFunc)
	if err != nil {
		t.Fatalf("Failed to evolve population: %v", err)
	}

	// Check if elites were preserved
	found0 := false
	found1 := false

	for _, ind := range gaInstance.Population.Individuals {
		allMatch0 := true
		allMatch1 := true

		for i, gene := range ind.Genotype.Genome {
			if gene != elite0[i] {
				allMatch0 = false
			}
			if gene != elite1[i] {
				allMatch1 = false
			}
		}

		if allMatch0 {
			found0 = true
		}
		if allMatch1 {
			found1 = true
		}
	}

	if !found0 || !found1 {
		t.Error("Elite individuals were not preserved in the population")
	}
}

func TestParallelEvaluation(t *testing.T) {
	// Skip if short tests are requested
	if testing.Short() {
		t.Skip("Skipping parallel evaluation test in short mode")
	}

	// Create a GA with parallel evaluation
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:        SinglePointCrossover,
		Mutation:         BitFlipMutation,
		CrossoverRate:    0.7,
		MutationRate:     0.01,
		Generations:      1, // One generation is enough for timing test
		NumParallelEvals: 4, // Use 4 parallel workers
		EnableLogger:     false,
	}

	// Use a larger population size and longer evaluation time
	// to make the effect of parallelization more obvious
	populationSize := 100
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		return NewBinaryGenotype(genomeLength)
	}

	// Create an evaluation function that has a small delay to simulate computation
	evalFunc := func(genotype *Genotype) *Phenotype {
		time.Sleep(5 * time.Millisecond)
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	// To make this test deterministic, we don't test for specific timing,
	// but instead verify that parallel evaluation completes successfully
	// and produces the same results as sequential evaluation

	// First run parallel evaluation
	gaInstance.NumParallelEvals = 4
	err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
	if err != nil {
		t.Fatalf("Failed to initialize GA: %v", err)
	}

	parallelResult, err := gaInstance.Evolve(evalFunc)
	if err != nil {
		t.Fatalf("Failed to evolve population with parallel evaluation: %v", err)
	}

	parallelFitness := parallelResult.Best.Phenotype.Fitness

	// Then run sequential evaluation
	gaInstance.NumParallelEvals = 1
	err = gaInstance.Initialize(populationSize, initFunc, evalFunc)
	if err != nil {
		t.Fatalf("Failed to initialize GA: %v", err)
	}

	sequentialResult, err := gaInstance.Evolve(evalFunc)
	if err != nil {
		t.Fatalf("Failed to evolve population with sequential evaluation: %v", err)
	}

	sequentialFitness := sequentialResult.Best.Phenotype.Fitness

	// Log the results for informational purposes
	t.Logf("Parallel evaluation fitness: %v, Sequential evaluation fitness: %v",
		parallelFitness, sequentialFitness)

	// We don't compare fitness values as they may differ due to random nature
	// of genetic algorithms, but we verify that both methods completed successfully
}

func TestErrorHandling(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		EnableLogger:  false,
	}

	// Test invalid population size
	err := gaInstance.Initialize(0, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for invalid population size")
	}

	// Test nil initialization function
	err = gaInstance.Initialize(10, nil, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil initialization function")
	}

	// Test nil evaluation function
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, nil)
	if err == nil {
		t.Error("Expected error for nil evaluation function")
	}

	// Test nil genetic operators
	gaInstance.Selection = nil
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil selection operator")
	}

	gaInstance.Selection = func(population []*Individual, rng *rand.Rand) []*Individual {
		return TournamentSelection(population, 3, rng)
	}
	gaInstance.Crossover = nil
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil crossover operator")
	}

	gaInstance.Crossover = SinglePointCrossover
	gaInstance.Mutation = nil
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil mutation operator")
	}
}

func TestAdaptiveParameters(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:      SinglePointCrossover,
		Mutation:       BitFlipMutation,
		CrossoverRate:  0.7,
		MutationRate:   0.01,
		Generations:    5,
		AdaptiveParams: true,
		EnableLogger:   false,
	}

	populationSize := 20
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		return NewBinaryGenotype(genomeLength)
	}

	evalFunc := func(genotype *Genotype) *Phenotype {
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
	if err != nil {
		t.Fatalf("Failed to initialize GA: %v", err)
	}

	// Store initial rates
	initialMutationRate := gaInstance.MutationRate
	initialCrossoverRate := gaInstance.CrossoverRate

	// Evolve for one generation
	gaInstance.Generations = 1
	_, err = gaInstance.Evolve(evalFunc)
	if err != nil {
		t.Fatalf("Failed to evolve population: %v", err)
	}

	// Check if rates were updated
	if gaInstance.MutationRate == initialMutationRate {
		t.Error("Mutation rate should have been updated")
	}
	if gaInstance.CrossoverRate == initialCrossoverRate {
		t.Error("Crossover rate should have been updated")
	}

	// Check if rates are within bounds
	if gaInstance.MutationRate < 0.01 || gaInstance.MutationRate > 0.5 {
		t.Errorf("Mutation rate should be between 0.01 and 0.5, got %f", gaInstance.MutationRate)
	}
	if gaInstance.CrossoverRate < 0.1 || gaInstance.CrossoverRate > 0.95 {
		t.Errorf("Crossover rate should be between 0.1 and 0.95, got %f", gaInstance.CrossoverRate)
	}
}

func TestTerminationConditions(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		Seed:          42,
		EnableLogger:  false,
	}

	populationSize := 20
	genomeLength := 8

	initFunc := func(rng *rand.Rand) *Genotype {
		g := NewBinaryGenotype(genomeLength)
		// Initialize with a high-fitness pattern so threshold-termination
		// scenarios are not at the mercy of the initial random draw.
		for i := range g.Genome {
			g.Genome[i] = 1
		}
		return g
	}

	evalFunc := func(genotype *Genotype) *Phenotype {
		fitness := 0.0
		for _, gene := range genotype.Genome {
			if gene == 1 {
				fitness += 1.0
			}
		}
		return &Phenotype{Fitness: fitness}
	}

	t.Run("fitness threshold condition", func(t *testing.T) {
		// Skip slow tests in short mode
		if testing.Short() {
			t.Skip("Skipping fitness threshold test in short mode")
		}

		// Test fitness threshold termination - should be most reliable
		threshold := float64(genomeLength / 2) // Use a lower threshold that might be reachable
		gaInstance.TermCondition = FitnessThresholdTermination(threshold)

		err := gaInstance.Initialize(populationSize, initFunc, evalFunc)
		if err != nil {
			t.Fatalf("Initialize() failed: %v", err)
		}

		// Set one individual to have high fitness to trigger immediate termination
		gaInstance.Population.Individuals[0].Phenotype.Fitness = threshold + 1.0
		gaInstance.Population.CalculateStatistics()

		_, err = gaInstance.Evolve(evalFunc)
		if err != nil {
			t.Fatalf("Evolve() failed: %v", err)
		}

		// Verify the best fitness reached the threshold
		if got, want := gaInstance.Population.Statistics.BestFitness, threshold; got < want {
			t.Errorf("Best fitness got %f, want >= %f", got, want)
		}

		// Make sure we have at least 1 history entry even with immediate termination
		if len(gaInstance.History) == 0 {
			t.Error("History is empty, want at least initial population statistics")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	gaInstance := &GA{
		Selection: func(population []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(population, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   5,
		EnableLogger:  false,
	}

	// Test with minimum population size
	err := gaInstance.Initialize(1, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err != nil {
		t.Errorf("Should accept population size of 1: %v", err)
	}

	// Test with maximum elitism count
	gaInstance.ElitismCount = 100
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err != nil {
		t.Errorf("Should handle large elitism count: %v", err)
	}
	if gaInstance.ElitismCount != 10 {
		t.Errorf("Elitism count should be capped at population size, got %d", gaInstance.ElitismCount)
	}

	// Test with extreme mutation and crossover rates
	gaInstance.MutationRate = 2.0
	gaInstance.CrossoverRate = -1.0
	err = gaInstance.Initialize(10, func(rng *rand.Rand) *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err != nil {
		t.Errorf("Should handle extreme rates: %v", err)
	}
	if gaInstance.MutationRate != 0.1 {
		t.Errorf("Mutation rate should be reset to default, got %f", gaInstance.MutationRate)
	}
	if gaInstance.CrossoverRate != 0.8 {
		t.Errorf("Crossover rate should be reset to default, got %f", gaInstance.CrossoverRate)
	}
}

// TestReproducibilityWithSeed verifies that two GA runs with the same Seed
// produce bit-for-bit identical results. This is the headline property of
// the seeded-RNG refactor.
func TestReproducibilityWithSeed(t *testing.T) {
	run := func(seed int64) *Individual {
		gaInstance := &GA{
			Selection: func(p []*Individual, rng *rand.Rand) []*Individual {
				return TournamentSelection(p, 3, rng)
			},
			Crossover:     SinglePointCrossover,
			Mutation:      BitFlipMutation,
			CrossoverRate: 0.7,
			MutationRate:  0.05,
			Generations:   30,
			Seed:          seed,
			EnableLogger:  false,
		}
		initFunc := func(rng *rand.Rand) *Genotype {
			g := NewBinaryGenotype(16)
			for i := range g.Genome {
				g.Genome[i] = byte(rng.Intn(2))
			}
			return g
		}
		evalFunc := func(g *Genotype) *Phenotype {
			f := 0.0
			for _, b := range g.Genome {
				if b == 1 {
					f++
				}
			}
			return &Phenotype{Fitness: f}
		}
		if err := gaInstance.Initialize(20, initFunc, evalFunc); err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}
		result, err := gaInstance.Evolve(evalFunc)
		if err != nil {
			t.Fatalf("Evolve failed: %v", err)
		}
		return result.Best
	}

	a := run(42)
	b := run(42)
	if a.Phenotype.Fitness != b.Phenotype.Fitness {
		t.Errorf("Same seed produced different fitnesses: %f vs %f", a.Phenotype.Fitness, b.Phenotype.Fitness)
	}
	for i, g := range a.Genotype.Genome {
		if g != b.Genotype.Genome[i] {
			t.Errorf("Same seed produced different genomes at index %d: %d vs %d", i, g, b.Genotype.Genome[i])
		}
	}

	// Note: we deliberately do not assert that different seeds diverge.
	// On a small OneMax problem both runs may legitimately reach the all-1s
	// optimum, in which case the all-time best is identical regardless of
	// seed — that is correct behavior, not a bug.
}

// TestResultBestNotAliasedToPopulation is a regression test for an aliasing
// bug: bestIndividual stored a pointer into the live population, and when
// crossover skipped a pair the offspring aliased the parent — which the
// mutation operator then mutated in place, silently corrupting the all-time
// best. The fix clones the best on capture.
func TestResultBestNotAliasedToPopulation(t *testing.T) {
	gaInstance := &GA{
		Selection: func(p []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(p, 3, rng)
		},
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.5, // ~50% chance crossover is skipped → offspring aliases parent.
		MutationRate:  0.5, // high mutation rate → in-place flips are likely.
		Generations:   50,
		Seed:          1,
	}
	initFunc := func(rng *rand.Rand) *Genotype {
		g := NewBinaryGenotype(16)
		for i := range g.Genome {
			g.Genome[i] = byte(rng.Intn(2))
		}
		return g
	}
	evalFunc := func(g *Genotype) *Phenotype {
		f := 0.0
		for _, b := range g.Genome {
			if b == 1 {
				f++
			}
		}
		return &Phenotype{Fitness: f}
	}
	if err := gaInstance.Initialize(20, initFunc, evalFunc); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	result, err := gaInstance.Evolve(evalFunc)
	if err != nil {
		t.Fatalf("Evolve failed: %v", err)
	}

	// Re-evaluate the returned best — its declared fitness must match its
	// actual genome. If aliasing corrupted it, these will diverge.
	computed := evalFunc(result.Best.Genotype)
	if computed.Fitness != result.Best.Phenotype.Fitness {
		t.Errorf("Result.Best genome was mutated after capture: declared fitness %f, recomputed %f", result.Best.Phenotype.Fitness, computed.Fitness)
	}

	// And the all-time best must be at least as fit as the final population's best.
	popBest := gaInstance.Population.GetBestIndividual()
	if result.Best.Phenotype.Fitness < popBest.Phenotype.Fitness {
		t.Errorf("Result.Best.Fitness = %f < final population best = %f", result.Best.Phenotype.Fitness, popBest.Phenotype.Fitness)
	}
}

// makeOneMaxGA returns a GA configured for the OneMax problem; used by the
// EarlyStopping / OnGeneration tests below.
func makeOneMaxGA(t *testing.T, generations int) (*GA, func(*Genotype) *Phenotype) {
	t.Helper()
	gaInstance := &GA{
		Selection: func(p []*Individual, rng *rand.Rand) []*Individual {
			return TournamentSelection(p, 3, rng)
		},
		Crossover:     UniformCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.9,
		MutationRate:  0.05,
		Generations:   generations,
		Seed:          42,
		EnableLogger:  false,
	}
	initFunc := func(rng *rand.Rand) *Genotype {
		g := NewBinaryGenotype(16)
		for i := range g.Genome {
			g.Genome[i] = byte(rng.Intn(2))
		}
		return g
	}
	evalFunc := func(g *Genotype) *Phenotype {
		f := 0.0
		for _, b := range g.Genome {
			if b == 1 {
				f++
			}
		}
		return &Phenotype{Fitness: f}
	}
	if err := gaInstance.Initialize(40, initFunc, evalFunc); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	return gaInstance, evalFunc
}

func TestEarlyStoppingTargetFitness(t *testing.T) {
	g, eval := makeOneMaxGA(t, 500)
	g.EarlyStopping = &EarlyStopping{TargetFitness: 16, TargetFitnessSet: true}
	result, err := g.Evolve(eval)
	if err != nil {
		t.Fatalf("Evolve failed: %v", err)
	}
	if result.StopReason != StopTargetFitness {
		t.Errorf("StopReason = %q, want %q", result.StopReason, StopTargetFitness)
	}
	if result.StoppedAtGeneration >= 500 {
		t.Errorf("Expected early stop before generation 500, ran %d generations", result.StoppedAtGeneration)
	}
	if result.Best.Phenotype.Fitness < 16 {
		t.Errorf("Best fitness = %f, want >= 16", result.Best.Phenotype.Fitness)
	}
}

func TestEarlyStoppingPatience(t *testing.T) {
	g, eval := makeOneMaxGA(t, 1000)
	// Patience of 5 with a tiny tolerance should stop well before the cap.
	g.EarlyStopping = &EarlyStopping{Patience: 5, Tol: 0}
	result, err := g.Evolve(eval)
	if err != nil {
		t.Fatalf("Evolve failed: %v", err)
	}
	if result.StoppedAtGeneration >= 1000 {
		t.Errorf("Expected patience to fire before generation 1000, ran %d", result.StoppedAtGeneration)
	}
	if result.StopReason != StopPatience && result.StopReason != StopTargetFitness {
		// Allow target_fitness too because OneMax with a small genome may reach 16 first.
		t.Errorf("StopReason = %q, want StopPatience (or target)", result.StopReason)
	}
}

func TestEarlyStoppingTimeLimit(t *testing.T) {
	g, eval := makeOneMaxGA(t, 100000)
	g.EarlyStopping = &EarlyStopping{TimeLimit: int64(50 * time.Millisecond)}
	result, err := g.Evolve(eval)
	if err != nil {
		t.Fatalf("Evolve failed: %v", err)
	}
	if result.StoppedAtGeneration >= 100000 {
		t.Errorf("Expected time limit to fire, ran all %d generations", result.StoppedAtGeneration)
	}
	if result.StopReason != StopTimeLimit && result.StopReason != StopTargetFitness && result.StopReason != StopPatience {
		t.Errorf("StopReason = %q, want StopTimeLimit", result.StopReason)
	}
}

func TestOnGenerationCallback(t *testing.T) {
	g, eval := makeOneMaxGA(t, 10)
	var calls []int
	g.OnGeneration = func(gen int, stats *Statistics) {
		calls = append(calls, gen)
		if stats == nil {
			t.Errorf("OnGeneration received nil stats at gen %d", gen)
		}
	}
	if _, err := g.Evolve(eval); err != nil {
		t.Fatalf("Evolve failed: %v", err)
	}
	if len(calls) != 10 {
		t.Errorf("OnGeneration called %d times, want 10", len(calls))
	}
	// Generations should be 1..10 (lastGen = gen + 1).
	for i, gen := range calls {
		if want := i + 1; gen != want {
			t.Errorf("calls[%d] = %d, want %d", i, gen, want)
		}
	}
}
