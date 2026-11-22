package ga

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/Okabe-Junya/gago/internal/logger"
)

var (
	errInvalidPopulationSize = errors.New("population size must be positive")
	errNilInitializeGenotype = errors.New("initializeGenotype function cannot be nil")
	errNilEvaluatePhenotype  = errors.New("evaluatePhenotype function cannot be nil")
	errNilSelection          = errors.New("selection operator cannot be nil")
	errNilCrossover          = errors.New("crossover operator cannot be nil")
	errNilMutation           = errors.New("mutation operator cannot be nil")
	errNoValidIndividuals    = errors.New("population contains no valid individuals")
	errEmptySelection        = errors.New("selection operator returned empty population")
	errEmptyCrossover        = errors.New("crossover operator returned empty population")
	errFitnessPanic          = errors.New("panic during fitness evaluation")
	errNilPhenotype          = errors.New("evaluatePhenotype returned nil phenotype")
)

const (
	defaultMutationRate  = 0.1
	defaultCrossoverRate = 0.8
	minAdaptiveMutation  = 0.01
	maxAdaptiveMutation  = 0.5
	minAdaptiveCrossover = 0.1
	maxAdaptiveCrossover = 0.95
)

// TerminationCondition defines a condition for terminating the GA evolution process.
type TerminationCondition interface {
	Evaluate(ga *GA) bool
}

// TerminationConditionFunc is a function type that implements TerminationCondition.
type TerminationConditionFunc func(*GA) bool

// Evaluate implements the TerminationCondition interface.
func (f TerminationConditionFunc) Evaluate(ga *GA) bool {
	return f(ga)
}

// GA represents the genetic algorithm, including its population, genetic operators,
// and parameters for crossover and mutation rates, and the number of generations to evolve.
//
// Reproducibility: when Seed is non-zero, the GA owns a single *rand.Rand
// seeded with that value; every random decision in the main evolution loop
// (selection, crossover, mutation, initialization) draws from it, so runs
// are bit-for-bit reproducible. When Seed is 0, the rng is seeded from
// time.Now().UnixNano(). The rng is single-goroutine only; parallel fitness
// evaluation must not touch it.
type GA struct {
	StartTime        time.Time
	Logger           *logger.Logger
	rng              *rand.Rand
	Selection        func([]*Individual, *rand.Rand) []*Individual
	Crossover        func([]*Individual, float64, *rand.Rand) []*Individual
	Mutation         func([]*Individual, float64, *rand.Rand)
	TermCondition    TerminationCondition
	EarlyStopping    *EarlyStopping
	OnGeneration     func(gen int, stats *Statistics)
	Population       *Population
	History          []*Statistics
	Generations      int
	ElitismCount     int
	NumParallelEvals int
	Seed             int64
	MutationRate     float64
	CrossoverRate    float64
	LogLevel         logger.LogLevel
	AdaptiveParams   bool
	EnableLogger     bool
	LogJSON          bool
}

// timeNow is swapped by tests to control the clock seen by time-based termination.
var timeNow = time.Now //nolint:gochecknoglobals // test seam for the clock

// RNG returns the GA's random source. May be nil if Initialize has not run yet.
// Callers that need to generate random numbers consistent with the GA's seed
// (e.g., custom genotype initializers) should draw from this.
func (ga *GA) RNG() *rand.Rand {
	return ga.rng
}

// Initialize initializes the population with the specified size and configuration.
// It creates and evaluates the initial population using the provided functions.
//
// Parameters:
//   - populationSize: The size of the population to initialize.
//   - initializeGenotype: A function that creates a new genotype for an individual.
//     It receives the GA's seeded *rand.Rand so initialization is reproducible.
//   - evaluatePhenotype: A function that evaluates a genotype and returns its phenotype.
//
// Returns an error if the input parameters are invalid or if any of the required
// genetic operators are not provided.
func (ga *GA) Initialize(
	populationSize int,
	initializeGenotype func(*rand.Rand) *Genotype,
	evaluatePhenotype func(*Genotype) *Phenotype,
) error {
	if populationSize <= 0 {
		return fmt.Errorf("%w, got %d", errInvalidPopulationSize, populationSize)
	}

	if initializeGenotype == nil {
		return errNilInitializeGenotype
	}

	if evaluatePhenotype == nil {
		return errNilEvaluatePhenotype
	}

	ga.seedRNG()

	initFunc := func() *Individual {
		genotype := initializeGenotype(ga.rng)
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

	ga.History = make([]*Statistics, 0, ga.Generations+1)
	ga.History = append(ga.History, ga.Population.Statistics)

	ga.applyDefaults(populationSize)

	err := ga.validateOperators()
	if err != nil {
		return err
	}

	ga.applyDefaultRates()

	if ga.EnableLogger {
		ga.initializeLogger()
	}

	ga.StartTime = timeNow()

	return nil
}

// Evolve runs the genetic algorithm for the specified number of generations.
// It applies the genetic operators (selection, crossover, mutation) to evolve
// the population and evaluates each new individual.
//
// The evolution process continues until either:
//   - The maximum number of generations is reached
//   - GA.TermCondition.Evaluate returns true
//   - Any criterion configured on GA.EarlyStopping fires (patience exhausted,
//     target fitness reached, or wall-clock time limit hit)
//
// Parameters:
//   - evaluatePhenotype: A function that evaluates a genotype and returns its phenotype.
//
// Returns:
//   - A *Result containing the all-time best individual, the final population,
//     history, the StopReason that ended the run, and the index of the last
//     generation that ran.
//   - An error if any step of the evolution process fails; in that case
//     Result.StopReason is StopError and the returned Result still contains
//     whatever progress had been made before the failure.
func (ga *GA) Evolve(evaluatePhenotype func(*Genotype) *Phenotype) (*Result, error) {
	if evaluatePhenotype == nil {
		return nil, errNilEvaluatePhenotype
	}

	ga.StartTime = timeNow()

	bestIndividual := ga.Population.GetBestIndividual()
	if bestIndividual == nil {
		return nil, fmt.Errorf("initial %w", errNoValidIndividuals)
	}
	// Clone so the all-time best is decoupled from the live population. Without
	// this, gen-0's best can be selected, aliased into offspring via crossover's
	// pass-through branch, mutated in place, and returned as Result.Best with a
	// Genome that no longer matches its reported Fitness. The improvement path
	// below clones for the same reason.
	best := bestTracker{individual: ga.cloneIndividual(bestIndividual)}
	best.fitness = best.individual.Phenotype.Fitness
	stopReason := StopMaxGenerations
	lastGen := 0
	elites := make([]*Individual, 0, ga.ElitismCount)

	for gen := range ga.Generations {
		ga.logGenerationStats(gen, best.fitness, best.noImprovementCount)

		var err error

		elites, err = ga.advanceGeneration(gen, evaluatePhenotype, elites)
		if err != nil {
			return nil, err
		}

		currentBest := ga.Population.GetBestIndividual()
		if currentBest == nil {
			return nil, fmt.Errorf("%w at generation %d", errNoValidIndividuals, gen)
		}

		best.observe(ga, currentBest)

		ga.History = append(ga.History, ga.Population.Statistics)
		lastGen = gen + 1

		if ga.OnGeneration != nil {
			ga.OnGeneration(lastGen, ga.Population.Statistics)
		}

		if reason, stop := ga.shouldStop(gen, best.noImprovementCount); stop {
			stopReason = reason

			break
		}
	}

	return &Result{
		Best:                best.individual,
		Population:          ga.Population,
		History:             ga.History,
		StopReason:          stopReason,
		StoppedAtGeneration: lastGen,
	}, nil
}

// GetStatistics returns the current statistics of the population.
func (ga *GA) GetStatistics() *Statistics {
	return ga.Population.Statistics
}

// GetRuntime returns the elapsed time since evolution started.
func (ga *GA) GetRuntime() time.Duration {
	return time.Since(ga.StartTime)
}

func (ga *GA) seedRNG() {
	// Seed == 0 means "auto-seed from time".
	if ga.Seed == 0 {
		ga.rng = rand.New(rand.NewSource(timeNow().UnixNano()))
	} else {
		ga.rng = rand.New(rand.NewSource(ga.Seed))
	}
}

func (ga *GA) applyDefaults(populationSize int) {
	if ga.NumParallelEvals <= 0 {
		ga.NumParallelEvals = runtime.NumCPU()
	}

	ga.ElitismCount = max(0, min(ga.ElitismCount, populationSize))

	if ga.TermCondition == nil {
		ga.TermCondition = TerminationConditionFunc(func(_ *GA) bool {
			return false
		})
	}
}

func (ga *GA) validateOperators() error {
	switch {
	case ga.Selection == nil:
		return errNilSelection
	case ga.Crossover == nil:
		return errNilCrossover
	case ga.Mutation == nil:
		return errNilMutation
	default:
		return nil
	}
}

func (ga *GA) applyDefaultRates() {
	if ga.MutationRate <= 0 || ga.MutationRate > 1 {
		ga.MutationRate = defaultMutationRate
	}

	if ga.CrossoverRate <= 0 || ga.CrossoverRate > 1 {
		ga.CrossoverRate = defaultCrossoverRate
	}
}

func (ga *GA) logGenerationStats(gen int, bestFitness float64, noImprovementCount int) {
	if ga.Logger == nil {
		return
	}

	genStartTime := time.Now()
	stats := map[string]any{
		"generation":     gen,
		"bestFitness":    bestFitness,
		"averageFitness": ga.Population.Statistics.AverageFitness,
		"diversity":      ga.Population.Statistics.Diversity,
		"noImprovement":  noImprovementCount,
	}
	ga.Logger.LogGenerationStats(gen, stats, time.Since(genStartTime))
}

// advanceGeneration replaces the population with the next evaluated generation,
// carrying elites over. It returns the elites slice for reuse in the next call.
func (ga *GA) advanceGeneration(
	gen int,
	evaluatePhenotype func(*Genotype) *Phenotype,
	elites []*Individual,
) ([]*Individual, error) {
	selectedIndividuals := ga.Selection(ga.Population.Individuals, ga.rng)
	if len(selectedIndividuals) == 0 {
		return elites, fmt.Errorf("%w at generation %d", errEmptySelection, gen)
	}

	offspring := ga.Crossover(selectedIndividuals, ga.CrossoverRate, ga.rng)
	if len(offspring) == 0 {
		return elites, fmt.Errorf("%w at generation %d", errEmptyCrossover, gen)
	}

	ga.Mutation(offspring, ga.MutationRate, ga.rng)

	elites = ga.snapshotElites(elites)

	if ga.AdaptiveParams {
		ga.adaptParams()
	}

	evalStartTime := time.Now()
	ga.evaluatePopulationInParallel(offspring, evaluatePhenotype)

	if ga.Logger != nil {
		ga.Logger.Debug("Population evaluation completed",
			"time", time.Since(evalStartTime),
			"parallelWorkers", ga.NumParallelEvals)
	}

	ga.Population.Individuals = offspring
	for i, elite := range elites {
		ga.Population.Replace(i, elite)
	}

	ga.Population.CalculateStatistics()

	return elites, nil
}

func (ga *GA) snapshotElites(elites []*Individual) []*Individual {
	if ga.ElitismCount <= 0 {
		return elites
	}

	ga.Population.SortByFitness()

	elites = elites[:0]
	for i := 0; i < ga.ElitismCount && i < len(ga.Population.Individuals); i++ {
		elites = append(elites, ga.cloneIndividual(ga.Population.Individuals[i]))
	}

	return elites
}

func (ga *GA) adaptParams() {
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

// bestTracker holds the all-time best individual of a run.
type bestTracker struct {
	individual         *Individual
	fitness            float64
	noImprovementCount int
}

func (b *bestTracker) observe(ga *GA, current *Individual) {
	// EarlyStopping.Tol turns "any improvement" into "improvement > Tol";
	// reuse the same threshold for noImprovementCount.
	tol := 0.0
	if ga.EarlyStopping != nil {
		tol = ga.EarlyStopping.Tol
	}

	if current.Phenotype.Fitness <= b.fitness+tol {
		b.noImprovementCount++

		return
	}

	// Clone so the all-time best is immune to in-place mutation if this
	// pointer survives into the next generation as a parent and gets mutated
	// by the mutation operator (crossover skips alias).
	b.individual = ga.cloneIndividual(current)
	b.fitness = current.Phenotype.Fitness
	b.noImprovementCount = 0
}

// shouldStop checks EarlyStopping before TermCondition so that StopReason
// reports the more specific reason when both would fire on the same generation.
func (ga *GA) shouldStop(gen, noImprovementCount int) (StopReason, bool) {
	if ga.EarlyStopping != nil {
		reason, stop := ga.EarlyStopping.check(
			ga.Population.Statistics.BestFitness, time.Since(ga.StartTime), noImprovementCount)
		if stop {
			return reason, true
		}
	}

	if ga.TermCondition != nil && ga.TermCondition.Evaluate(ga) {
		if ga.Logger != nil {
			ga.Logger.Info("Evolution terminated",
				"reason", "Termination condition met",
				"generation", gen,
				"totalRuntime", time.Since(ga.StartTime))
		}

		return StopTermCondition, true
	}

	return "", false
}

type evalJob struct {
	individual *Individual
	index      int
}

type evalResult struct {
	err       error
	phenotype *Phenotype
	index     int
}

// evaluatePopulationInParallel evaluates the fitness of individuals with a
// worker pool. A panicking or nil-returning fitness function does not abort the
// run: that individual gets -math.MaxFloat64 fitness and the error is logged.
func (ga *GA) evaluatePopulationInParallel(population []*Individual, evaluatePhenotype func(*Genotype) *Phenotype) {
	if len(population) == 0 {
		return
	}

	if ga.NumParallelEvals <= 1 {
		evaluateSequentially(population, evaluatePhenotype)

		return
	}

	jobs := make(chan evalJob, len(population))
	results := make(chan evalResult, len(population))

	var wg sync.WaitGroup
	for workerID := range min(ga.NumParallelEvals, len(population)) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := range jobs {
				phenotype, err := ga.safeEvaluate(workerID, j, evaluatePhenotype)
				results <- evalResult{index: j.index, phenotype: phenotype, err: err}
			}
		}()
	}

	for i, ind := range population {
		if ind == nil {
			continue
		}

		jobs <- evalJob{index: i, individual: ind}
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	ga.collectEvaluations(population, results)
}

func evaluateSequentially(population []*Individual, evaluatePhenotype func(*Genotype) *Phenotype) {
	for _, ind := range population {
		if ind == nil {
			continue
		}

		ind.Phenotype = evaluatePhenotype(ind.Genotype)
	}
}

func (ga *GA) collectEvaluations(population []*Individual, results <-chan evalResult) {
	var evalErrors []error

	for res := range results {
		if res.err != nil {
			evalErrors = append(evalErrors, res.err)
			population[res.index].Phenotype = &Phenotype{Fitness: -math.MaxFloat64}
		} else {
			population[res.index].Phenotype = res.phenotype
		}
	}

	if len(evalErrors) > 0 && ga.Logger != nil {
		ga.Logger.Error("Errors during parallel evaluation",
			"errorCount", len(evalErrors),
			"errors", evalErrors)
	}
}

func (ga *GA) safeEvaluate(workerID int, j evalJob, evaluatePhenotype func(*Genotype) *Phenotype) (*Phenotype, error) {
	var (
		phenotype *Phenotype
		err       error
	)

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%w: %v", errFitnessPanic, r)

				if ga.Logger != nil {
					ga.Logger.Error("Panic in worker goroutine",
						"workerID", workerID,
						"individualIndex", j.index,
						"panic", fmt.Sprintf("%v", r))
				}
			}
		}()

		phenotype = evaluatePhenotype(j.individual.Genotype)
	}()

	if err != nil {
		return nil, err
	}

	if phenotype == nil {
		return nil, errNilPhenotype
	}

	return phenotype, nil
}

// updateAdaptiveParams raises the mutation rate and lowers the crossover rate
// as population diversity falls, and the reverse as it rises.
func (ga *GA) updateAdaptiveParams() {
	normalizedDiversity := math.Min(ga.Population.Statistics.Diversity, 1.0)

	ga.MutationRate = defaultMutationRate * (1.0 + (1.0 - normalizedDiversity))
	ga.CrossoverRate = defaultCrossoverRate * normalizedDiversity

	ga.MutationRate = math.Max(minAdaptiveMutation, math.Min(maxAdaptiveMutation, ga.MutationRate))
	ga.CrossoverRate = math.Max(minAdaptiveCrossover, math.Min(maxAdaptiveCrossover, ga.CrossoverRate))
}

// cloneIndividual creates a deep copy of an individual.
func (ga *GA) cloneIndividual(ind *Individual) *Individual {
	if ind == nil {
		return nil
	}

	var featuresClone []float64
	if ind.Phenotype != nil {
		featuresClone = cloneNonEmpty(ind.Phenotype.Features)
	}

	return &Individual{
		Genotype: &Genotype{
			Genome:     append([]byte(nil), ind.Genotype.Genome...),
			MinValues:  cloneNonEmpty(ind.Genotype.MinValues),
			MaxValues:  cloneNonEmpty(ind.Genotype.MaxValues),
			GenomeType: ind.Genotype.GenomeType,
		},
		Phenotype: &Phenotype{
			Fitness:  ind.Phenotype.Fitness,
			Features: featuresClone,
		},
	}
}

func cloneNonEmpty(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	return append([]float64(nil), values...)
}

func (ga *GA) initializeLogger() {
	options := []logger.Option{}

	if ga.LogLevel != 0 {
		options = append(options, logger.WithLevel(ga.LogLevel))
	}

	if ga.LogJSON {
		options = append(options, logger.WithJSON())
	}

	ga.Logger = logger.NewLogger(ga.EnableLogger, options...)

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
