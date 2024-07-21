package main

import (
	"fmt"
	"math"
	"math/rand"

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
	gaInstance := &ga.GA{
		Selection:     func(population []*ga.Individual) []*ga.Individual { return ga.TournamentSelection(population, 3) },
		Crossover:     ga.SinglePointCrossover,
		Mutation:      ga.BitFlipMutation,
		CrossoverRate: crossoverRate,
		MutationRate:  mutationRate,
		Generations:   generations,
		EnableLogger:  true,
	}

	gaInstance.Initialize(populationSize, initializeGenotype, evaluatePhenotype)
	gaInstance.Evolve(evaluatePhenotype)

	bestIndividual := findBestIndividual(gaInstance.Population)
	bestX := decodeGenotype(bestIndividual.Genotype)

	fmt.Printf("Best x: %f, Fitness: %f\n", bestX, bestIndividual.Phenotype.Fitness)
}

func initializeGenotype() *ga.Genotype {
	genotype := ga.NewGenotype(genomeLength)
	for i := range genotype.Genome {
		genotype.Genome[i] = byte(rand.Intn(2))
	}
	return genotype
}

func evaluatePhenotype(genotype *ga.Genotype) *ga.Phenotype {
	x := decodeGenotype(genotype)
	fitness := x * math.Sin(x)
	return &ga.Phenotype{Fitness: fitness}
}

func decodeGenotype(genotype *ga.Genotype) float64 {
	var value int64
	for _, bit := range genotype.Genome {
		value = (value << 1) | int64(bit)
	}
	return lowerBound + (upperBound-lowerBound)*float64(value)/float64((1<<genomeLength)-1)
}

func findBestIndividual(population []*ga.Individual) *ga.Individual {
	best := population[0]
	for _, ind := range population {
		if ind.Phenotype.Fitness > best.Phenotype.Fitness {
			best = ind
		}
	}
	return best
}
