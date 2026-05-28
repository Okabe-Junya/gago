package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/internal/logger"
	"github.com/Okabe-Junya/gago/pkg/ga"
)

const seed = 42 // deterministic for reproducibility

const (
	populationSize = 50
	genomeLength   = 16 // Length of the binary representation of the genotype
	generations    = 100
	crossoverRate  = 0.7
	mutationRate   = 0.01
	lowerBound     = 0.0
	upperBound     = math.Pi
)

// main runs the genetic algorithm to find the maximum of the function f(x) = x * sin(x).
func main() {
	// Configure the GA with appropriate parameters
	gaInstance := &ga.GA{
		Selection: func(population []*ga.Individual, rng *rand.Rand) []*ga.Individual {
			return ga.TournamentSelection(population, 3, rng)
		},
		Crossover:     ga.SinglePointCrossover,
		Mutation:      ga.BitFlipMutation,
		CrossoverRate: crossoverRate,
		MutationRate:  mutationRate,
		Generations:   generations,
		Seed:          seed,
		EnableLogger:  true,
		LogLevel:      logger.LevelInfo, // Set appropriate log level
	}

	// Initialize and run the GA
	gaInstance.Initialize(populationSize, initializeGenotype, evaluatePhenotype)

	// Set a termination condition based on fitness threshold
	gaInstance.TermCondition = ga.FitnessThresholdTermination(2.0) // f(x) = x*sin(x) has a maximum of ~2.0 around x = π/2

	result, err := gaInstance.Evolve(evaluatePhenotype)
	if err != nil {
		fmt.Printf("Error during evolution: %v\n", err)
		os.Exit(1)
	}

	// Get the results
	if result == nil || result.Best == nil {
		fmt.Println("Error: Failed to find a solution")
		os.Exit(1)
	}

	bestX := decodeGenotype(result.Best.Genotype)

	fmt.Printf("Best x: %f, Fitness: %f\n", bestX, result.Best.Phenotype.Fitness)
	fmt.Printf("Stop reason: %s\n", result.StopReason)
	fmt.Printf("Stopped at generation: %d\n", result.StoppedAtGeneration)
	fmt.Printf("Total runtime: %v\n", gaInstance.GetRuntime())
}

// initializeGenotype creates a new binary genotype for the problem.
func initializeGenotype(rng *rand.Rand) *ga.Genotype {
	genotype := ga.NewBinaryGenotype(genomeLength)
	for i := range genotype.Genome {
		genotype.Genome[i] = byte(rng.Intn(2))
	}
	return genotype
}

// evaluatePhenotype evaluates the fitness of a genotype.
func evaluatePhenotype(genotype *ga.Genotype) *ga.Phenotype {
	x := decodeGenotype(genotype)
	fitness := x * math.Sin(x) // Objective function f(x) = x*sin(x)
	return &ga.Phenotype{Fitness: fitness}
}

// decodeGenotype converts binary genotype to real value in [lowerBound, upperBound].
func decodeGenotype(genotype *ga.Genotype) float64 {
	var value int64
	for _, bit := range genotype.Genome {
		value = (value << 1) | int64(bit)
	}
	return lowerBound + (upperBound-lowerBound)*float64(value)/float64((1<<genomeLength)-1)
}
