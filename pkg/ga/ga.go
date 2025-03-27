// Package ga provides functionalities for implementing genetic algorithms,
// including the main GA struct and its methods for initialization and evolution.
package ga

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/Okabe-Junya/gago/internal/logger"
)

// TerminationCondition defines a condition for terminating the GA evolution process.
type TerminationCondition interface {
	Evaluate(*GA) bool
}

// TerminationConditionFunc is a function type that implements TerminationCondition.
type TerminationConditionFunc func(*GA) bool

// Evaluate implements the TerminationCondition interface.
func (f TerminationConditionFunc) Evaluate(ga *GA) bool {
	return f(ga)
}

// GA represents the genetic algorithm, including its population, genetic operators,
// and parameters for crossover and mutation rates, and the number of generations to evolve.
type GA struct {
	StartTime        time.Time
	Logger           *logger.Logger
	Selection        func([]*Individual) []*Individual
	Crossover        func([]*Individual, float64) []*Individual
	Mutation         func([]*Individual, float64)
	TermCondition    TerminationCondition
	Population       *Population
	History          []*Statistics
	Generations      int
	ElitismCount     int
	NumParallelEvals int
	MutationRate     float64
	CrossoverRate    float64
	LogLevel         logger.LogLevel
	AdaptiveParams   bool
	EnableLogger     bool
	LogJSON          bool
}

// Function variable for time operations, allows for test mocking
var timeNow = time.Now

// Initialize initializes the population with the specified size and configuration.
// It creates and evaluates the initial population using the provided functions.
//
// Parameters:
//   - populationSize: The size of the population to initialize.
//   - initializeGenotype: A function that creates a new genotype for an individual.
//   - evaluatePhenotype: A function that evaluates a genotype and returns its phenotype.
//
// Returns an error if the input parameters are invalid or if any of the required
// genetic operators are not provided.
func (ga *GA) Initialize(populationSize int, initializeGenotype func() *Genotype, evaluatePhenotype func(*Genotype) *Phenotype) error {
	// Validate input parameters
	if populationSize <= 0 {
		return fmt.Errorf("population size must be positive, got %d", populationSize)
	}
	if initializeGenotype == nil {
		return fmt.Errorf("initializeGenotype function cannot be nil")
	}
	if evaluatePhenotype == nil {
		return fmt.Errorf("evaluatePhenotype function cannot be nil")
	}

	// Create individuals with the initialization function
	initFunc := func() *Individual {
		genotype := initializeGenotype()
		if genotype == nil {
			panic("initializeGenotype returned nil genotype")
		}
		phenotype := evaluatePhenotype(genotype)
		if phenotype == nil {
			panic("evaluatePhenotype returned nil phenotype")
		}
		return &Individual{Genotype: genotype, Phenotype: phenotype}
	}

	ga.Population = NewPopulation(populationSize, initFunc)
	ga.Population.CalculateStatistics()

	// Initialize history
	ga.History = make([]*Statistics, 0, ga.Generations+1) // Pre-allocate capacity for efficiency
	ga.History = append(ga.History, ga.Population.Statistics)

	// Set default values for optional parameters
	if ga.NumParallelEvals <= 0 {
		ga.NumParallelEvals = runtime.NumCPU()
	}

	if ga.ElitismCount < 0 {
		ga.ElitismCount = 0
	} else if ga.ElitismCount > populationSize {
		ga.ElitismCount = populationSize
	}

	// Default termination condition if none provided
	if ga.TermCondition == nil {
		ga.TermCondition = TerminationConditionFunc(func(ga *GA) bool {
			return false
		})
	}

	// Validate genetic operators
	if ga.Selection == nil {
		return fmt.Errorf("selection operator cannot be nil")
	}
	if ga.Crossover == nil {
		return fmt.Errorf("crossover operator cannot be nil")
	}
	if ga.Mutation == nil {
		return fmt.Errorf("mutation operator cannot be nil")
	}

	// Validate rates
	if ga.MutationRate <= 0 || ga.MutationRate > 1 {
		ga.MutationRate = 0.1 // Default mutation rate
	}
	if ga.CrossoverRate <= 0 || ga.CrossoverRate > 1 {
		ga.CrossoverRate = 0.8 // Default crossover rate
	}

	if ga.EnableLogger {
		ga.initializeLogger()
	}

	// Initialize runtime tracking
	ga.StartTime = timeNow()
	return nil
}

// Evolve runs the genetic algorithm for the specified number of generations.
// It applies the genetic operators (selection, crossover, mutation) to evolve
// the population and evaluates each new individual.
//
// The evolution process continues until either:
// - The maximum number of generations is reached
// - A termination condition is met
//
// Parameters:
//   - evaluatePhenotype: A function that evaluates a genotype and returns its phenotype.
//
// Returns:
//   - The best individual found during the evolution process.
//   - An error if any step of the evolution process fails.
func (ga *GA) Evolve(evaluatePhenotype func(*Genotype) *Phenotype) (*Individual, error) {
	if evaluatePhenotype == nil {
		return nil, fmt.Errorf("evaluatePhenotype function cannot be nil")
	}

	// Reset start time for this evolution
	ga.StartTime = timeNow()
	bestIndividual := ga.Population.GetBestIndividual()
	if bestIndividual == nil {
		return nil, fmt.Errorf("initial population contains no valid individuals")
	}

	bestFitness := bestIndividual.Phenotype.Fitness
	noImprovementCount := 0

	// Pre-allocate slices for better performance
	var selectedIndividuals []*Individual
	var offspring []*Individual
	elites := make([]*Individual, 0, ga.ElitismCount)

	for gen := 0; gen < ga.Generations; gen++ {
		genStartTime := time.Now()

		// Log generation stats
		if ga.Logger != nil {
			stats := map[string]interface{}{
				"generation":     gen,
				"bestFitness":    bestIndividual.Phenotype.Fitness,
				"averageFitness": ga.Population.Statistics.AverageFitness,
				"diversity":      ga.Population.Statistics.Diversity,
				"noImprovement":  noImprovementCount,
			}
			ga.Logger.LogGenerationStats(gen, stats, time.Since(genStartTime))
		}

		// Apply genetic operators
		selectedIndividuals = ga.Selection(ga.Population.Individuals)
		if len(selectedIndividuals) == 0 {
			return nil, fmt.Errorf("selection operator returned empty population at generation %d", gen)
		}

		offspring = ga.Crossover(selectedIndividuals, ga.CrossoverRate)
		if len(offspring) == 0 {
			return nil, fmt.Errorf("crossover operator returned empty population at generation %d", gen)
		}

		ga.Mutation(offspring, ga.MutationRate)

		// Store elite individuals if elitism is enabled
		if ga.ElitismCount > 0 {
			ga.Population.SortByFitness()
			elites = elites[:0] // Reuse slice
			for i := 0; i < ga.ElitismCount && i < len(ga.Population.Individuals); i++ {
				elites = append(elites, ga.cloneIndividual(ga.Population.Individuals[i]))
			}
		}

		// Update mutation and crossover rates if adaptive parameters are enabled
		if ga.AdaptiveParams {
			oldMutationRate := ga.MutationRate
			oldCrossoverRate := ga.CrossoverRate
			ga.updateAdaptiveParams()
			if ga.Logger != nil {
				ga.Logger.Debug("Adaptive parameters updated",
					"oldMutationRate", oldMutationRate,
					"newMutationRate", ga.MutationRate,
					"oldCrossoverRate", oldCrossoverRate,
					"newCrossoverRate", ga.CrossoverRate)
			}
		}

		// Evaluate new population in parallel
		evalStartTime := time.Now()
		ga.evaluatePopulationInParallel(offspring, evaluatePhenotype)
		if ga.Logger != nil {
			ga.Logger.Debug("Population evaluation completed",
				"time", time.Since(evalStartTime),
				"parallelWorkers", ga.NumParallelEvals)
		}

		// Create new population
		ga.Population.Individuals = offspring

		// Reinsert elite individuals if elitism is enabled
		if ga.ElitismCount > 0 {
			for i, elite := range elites {
				if i < len(ga.Population.Individuals) {
					ga.Population.Replace(i, elite)
				}
			}
		}

		// Calculate statistics for the new population
		ga.Population.CalculateStatistics()

		// Update best individual and no improvement counter
		currentBest := ga.Population.GetBestIndividual()
		if currentBest == nil {
			return nil, fmt.Errorf("population contains no valid individuals at generation %d", gen)
		}

		if currentBest.Phenotype.Fitness > bestFitness {
			bestIndividual = currentBest
			bestFitness = currentBest.Phenotype.Fitness
			noImprovementCount = 0
		} else {
			noImprovementCount++
		}

		// Add current statistics to history
		ga.History = append(ga.History, ga.Population.Statistics)

		// Check termination condition after recording statistics
		if ga.TermCondition != nil && ga.TermCondition.Evaluate(ga) {
			if ga.Logger != nil {
				ga.Logger.Info("Evolution terminated",
					"reason", "Termination condition met",
					"generation", gen,
					"totalRuntime", time.Since(ga.StartTime))
			}
			break
		}
	}

	return bestIndividual, nil
}

// evaluatePopulationInParallel evaluates the fitness of individuals in parallel.
// It uses a worker pool pattern to process individuals efficiently and safely handles panics.
// The context allows for graceful cancellation of the evaluation process.
func (ga *GA) evaluatePopulationInParallel(population []*Individual, evaluatePhenotype func(*Genotype) *Phenotype) {
	if len(population) == 0 {
		return
	}

	// Use sequential evaluation if parallel processing is disabled
	if ga.NumParallelEvals <= 1 {
		for _, ind := range population {
			if ind == nil {
				continue
			}
			ind.Phenotype = evaluatePhenotype(ind.Genotype)
		}
		return
	}

	// Create a context that can be used to signal cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all resources are released when we're done

	// Create channels for work distribution and result collection
	type job struct {
		individual *Individual
		index      int
	}
	type result struct {
		err       error
		phenotype *Phenotype
		index     int
	}

	jobs := make(chan job, len(population))
	results := make(chan result, len(population))
	errorChan := make(chan error, len(population))

	// Start worker goroutines
	var wg sync.WaitGroup
	numWorkers := min(ga.NumParallelEvals, len(population))

	// Create a pool of worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Process jobs until the channel is closed or context is cancelled
			for {
				select {
				case <-ctx.Done():
					// Context was cancelled, stop processing
					return
				case j, ok := <-jobs:
					if !ok {
						// Channel closed, no more jobs
						return
					}

					// Evaluate the individual with panic recovery
					var phenotype *Phenotype
					var err error

					func() {
						defer func() {
							if r := recover(); r != nil {
								err = fmt.Errorf("panic during fitness evaluation: %v", r)
								if ga.Logger != nil {
									ga.Logger.Error("Panic in worker goroutine",
										"workerID", workerID,
										"individualIndex", j.index,
										"panic", fmt.Sprintf("%v", r))
								}
							}
						}()

						phenotype = evaluatePhenotype(j.individual.Genotype)
						if phenotype == nil {
							err = fmt.Errorf("evaluatePhenotype returned nil phenotype")
						}
					}()

					// Send the result back
					results <- result{
						index:     j.index,
						phenotype: phenotype,
						err:       err,
					}
				}
			}
		}(i)
	}

	// Send all individuals to the worker pool
	for i, ind := range population {
		if ind == nil {
			continue
		}
		jobs <- job{index: i, individual: ind}
	}

	// Close the jobs channel to signal no more work
	close(jobs)

	// Collect results in a separate goroutine to avoid deadlock
	// if a worker panics and can't send to the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process the results
	for res := range results {
		if res.err != nil {
			errorChan <- res.err
			// Assign a very low fitness to avoid further issues
			population[res.index].Phenotype = &Phenotype{Fitness: -math.MaxFloat64}
		} else {
			population[res.index].Phenotype = res.phenotype
		}
	}

	// Collect and log errors
	close(errorChan)
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 && ga.Logger != nil {
		ga.Logger.Error("Errors during parallel evaluation",
			"errorCount", len(errors),
			"errors", errors)
	}
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// updateAdaptiveParams updates mutation and crossover rates based on population diversity.
func (ga *GA) updateAdaptiveParams() {
	// Decrease mutation rate and increase crossover rate when diversity is high
	// Increase mutation rate and decrease crossover rate when diversity is low
	diversity := ga.Population.Statistics.Diversity

	// Normalize diversity to a value between 0 and 1
	normalizedDiversity := math.Min(diversity, 1.0)

	// Base mutation rate is higher when diversity is low
	baseMutationRate := 0.1
	baseCrossoverRate := 0.8

	// Adjust rates based on diversity
	// When diversity is low, increase mutation rate and decrease crossover rate
	// When diversity is high, decrease mutation rate and increase crossover rate
	ga.MutationRate = baseMutationRate * (1.0 + (1.0 - normalizedDiversity))
	ga.CrossoverRate = baseCrossoverRate * normalizedDiversity

	// Ensure rates stay within reasonable bounds
	ga.MutationRate = math.Max(0.01, math.Min(0.5, ga.MutationRate))
	ga.CrossoverRate = math.Max(0.1, math.Min(0.95, ga.CrossoverRate))
}

// cloneIndividual creates a deep copy of an individual.
func (ga *GA) cloneIndividual(ind *Individual) *Individual {
	if ind == nil {
		return nil
	}

	genomeClone := make([]byte, len(ind.Genotype.Genome))
	copy(genomeClone, ind.Genotype.Genome)

	// Clone MinValues if they exist
	var minValuesClone []float64
	if len(ind.Genotype.MinValues) > 0 {
		minValuesClone = make([]float64, len(ind.Genotype.MinValues))
		copy(minValuesClone, ind.Genotype.MinValues)
	}

	// Clone MaxValues if they exist
	var maxValuesClone []float64
	if len(ind.Genotype.MaxValues) > 0 {
		maxValuesClone = make([]float64, len(ind.Genotype.MaxValues))
		copy(maxValuesClone, ind.Genotype.MaxValues)
	}

	// Clone Features
	var featuresClone []float64
	if ind.Phenotype != nil && len(ind.Phenotype.Features) > 0 {
		featuresClone = make([]float64, len(ind.Phenotype.Features))
		copy(featuresClone, ind.Phenotype.Features)
	}

	return &Individual{
		Genotype: &Genotype{
			Genome:     genomeClone,
			MinValues:  minValuesClone,
			MaxValues:  maxValuesClone,
			GenomeType: ind.Genotype.GenomeType,
		},
		Phenotype: &Phenotype{
			Fitness:  ind.Phenotype.Fitness,
			Features: featuresClone,
		},
	}
}

// GetStatistics returns the current statistics of the population.
func (ga *GA) GetStatistics() *Statistics {
	return ga.Population.Statistics
}

// GetRuntime returns the elapsed time since evolution started.
func (ga *GA) GetRuntime() time.Duration {
	return time.Since(ga.StartTime)
}

// initializeLogger initializes the logger with the specified configuration.
func (ga *GA) initializeLogger() {
	options := []logger.LoggerOption{}

	// Set log level if specified
	if ga.LogLevel != 0 { // Assuming 0 is the default/unspecified value
		options = append(options, logger.WithLevel(ga.LogLevel))
	}

	// Set JSON format if requested
	if ga.LogJSON {
		options = append(options, logger.WithJSON())
	}

	ga.Logger = logger.NewLogger(ga.EnableLogger, options...)

	// Log initial configuration
	if ga.Logger != nil {
		ga.Logger.Info("Genetic algorithm initialized",
			"generations", ga.Generations,
			"populationSize", ga.Population.Size(),
			"crossoverRate", ga.CrossoverRate,
			"mutationRate", ga.MutationRate,
			"elitismCount", ga.ElitismCount,
			"adaptiveParams", ga.AdaptiveParams,
			"parallelEvals", ga.NumParallelEvals)
	}
}

// GenerationCountTermination returns a termination condition that terminates after a specified number of generations.
func GenerationCountTermination(maxGenerations int) TerminationCondition {
	return TerminationConditionFunc(func(ga *GA) bool {
		return len(ga.History) >= maxGenerations
	})
}

// ConvergenceTermination returns a termination condition that terminates when
// the best fitness hasn't improved by the specified threshold over the specified number of generations.
func ConvergenceTermination(noImprovementGens int, improvementThreshold float64) TerminationCondition {
	return TerminationConditionFunc(func(ga *GA) bool {
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
func TimeBasedTermination(duration time.Duration) TerminationCondition {
	return TerminationConditionFunc(func(ga *GA) bool {
		return timeNow().Sub(ga.StartTime) >= duration
	})
}

// FitnessThresholdTermination returns a termination condition that terminates when
// the best fitness reaches or exceeds the specified threshold.
func FitnessThresholdTermination(threshold float64) TerminationCondition {
	return TerminationConditionFunc(func(ga *GA) bool {
		return ga.Population.Statistics.BestFitness >= threshold
	})
}
