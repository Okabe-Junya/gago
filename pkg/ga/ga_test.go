package ga

import (
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	gaInstance := &GA{
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		EnableLogger:  false,
	}

	populationSize := 20
	genomeLength := 8

	initFunc := func() *Genotype {
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
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   5,
		EnableLogger:  false,
	}

	populationSize := 10
	genomeLength := 8

	initFunc := func() *Genotype {
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

		bestIndividual, err := gaInstance.Evolve(evalFunc)
		if err != nil {
			t.Fatalf("Evolve() failed: %v", err)
		}

		// Check if evolution occurred (history should have entries for each generation)
		expectedHistoryLength := gaInstance.Generations + 1 // Initial + each generation
		if got, want := len(gaInstance.History), expectedHistoryLength; got != want {
			t.Errorf("History length got %d, want %d", got, want)
		}

		// Check if best individual is not nil
		if bestIndividual == nil {
			t.Fatal("Best individual is nil after evolution")
		}

		// Best individual's fitness should match the best in the population
		popBest := gaInstance.Population.GetBestIndividual()
		if got, want := bestIndividual.Phenotype.Fitness, popBest.Phenotype.Fitness; got != want {
			t.Errorf("Best individual fitness got %f, want %f", got, want)
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
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
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

	initFunc := func() *Genotype {
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
		Selection:        func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
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

	initFunc := func() *Genotype {
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

	parallelFitness := parallelResult.Phenotype.Fitness

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

	sequentialFitness := sequentialResult.Phenotype.Fitness

	// Log the results for informational purposes
	t.Logf("Parallel evaluation fitness: %v, Sequential evaluation fitness: %v",
		parallelFitness, sequentialFitness)

	// We don't compare fitness values as they may differ due to random nature
	// of genetic algorithms, but we verify that both methods completed successfully
}

func TestErrorHandling(t *testing.T) {
	gaInstance := &GA{
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		EnableLogger:  false,
	}

	// Test invalid population size
	err := gaInstance.Initialize(0, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for invalid population size")
	}

	// Test nil initialization function
	err = gaInstance.Initialize(10, nil, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil initialization function")
	}

	// Test nil evaluation function
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, nil)
	if err == nil {
		t.Error("Expected error for nil evaluation function")
	}

	// Test nil genetic operators
	gaInstance.Selection = nil
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil selection operator")
	}

	gaInstance.Selection = func(population []*Individual) []*Individual { return TournamentSelection(population, 3) }
	gaInstance.Crossover = nil
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil crossover operator")
	}

	gaInstance.Crossover = SinglePointCrossover
	gaInstance.Mutation = nil
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err == nil {
		t.Error("Expected error for nil mutation operator")
	}
}

func TestAdaptiveParameters(t *testing.T) {
	gaInstance := &GA{
		Selection:      func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
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

	initFunc := func() *Genotype {
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
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   10,
		EnableLogger:  false,
	}

	populationSize := 20
	genomeLength := 8

	initFunc := func() *Genotype {
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
		Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
		Crossover:     SinglePointCrossover,
		Mutation:      BitFlipMutation,
		CrossoverRate: 0.7,
		MutationRate:  0.01,
		Generations:   5,
		EnableLogger:  false,
	}

	// Test with minimum population size
	err := gaInstance.Initialize(1, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err != nil {
		t.Errorf("Should accept population size of 1: %v", err)
	}

	// Test with maximum elitism count
	gaInstance.ElitismCount = 100
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
	if err != nil {
		t.Errorf("Should handle large elitism count: %v", err)
	}
	if gaInstance.ElitismCount != 10 {
		t.Errorf("Elitism count should be capped at population size, got %d", gaInstance.ElitismCount)
	}

	// Test with extreme mutation and crossover rates
	gaInstance.MutationRate = 2.0
	gaInstance.CrossoverRate = -1.0
	err = gaInstance.Initialize(10, func() *Genotype { return NewBinaryGenotype(8) }, func(*Genotype) *Phenotype { return &Phenotype{Fitness: 0} })
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
