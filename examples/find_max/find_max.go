package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/internal/logger"
	"github.com/Okabe-Junya/gago/pkg/ga"
)

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
		Selection:     func(population []*ga.Individual) []*ga.Individual { return ga.TournamentSelection(population, 3) },
		Crossover:     ga.SinglePointCrossover,
		Mutation:      ga.BitFlipMutation,
		CrossoverRate: crossoverRate,
		MutationRate:  mutationRate,
		Generations:   generations,
		EnableLogger:  true,
		LogLevel:      logger.LevelInfo, // Set appropriate log level
	}

	// Initialize and run the GA
	gaInstance.Initialize(populationSize, initializeGenotype, evaluatePhenotype)

	// Set a termination condition based on fitness threshold
	gaInstance.TermCondition = ga.FitnessThresholdTermination(2.0) // f(x) = x*sin(x) has a maximum of ~2.0 around x = π/2

	bestIndividual, err := gaInstance.Evolve(evaluatePhenotype)
	if err != nil {
		fmt.Printf("Error during evolution: %v\n", err)
		os.Exit(1)
	}

	// Get the results
	if bestIndividual == nil {
		fmt.Println("Error: Failed to find a solution")
		os.Exit(1)
	}

	bestX := decodeGenotype(bestIndividual.Genotype)

	fmt.Printf("Best x: %f, Fitness: %f\n", bestX, bestIndividual.Phenotype.Fitness)
	fmt.Printf("Total generations: %d\n", len(gaInstance.History)-1)
	fmt.Printf("Total runtime: %v\n", gaInstance.GetRuntime())
}

// initializeGenotype creates a new binary genotype for the problem.
func initializeGenotype() *ga.Genotype {
	genotype := ga.NewBinaryGenotype(genomeLength)
	for i := range genotype.Genome {
		genotype.Genome[i] = byte(rand.Intn(2))
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
