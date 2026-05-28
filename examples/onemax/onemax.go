// OneMax: maximize the number of 1s in a binary chromosome.
//
// This is the canonical GA baseline problem. With BitFlipMutation and any
// standard crossover, the GA should reach the all-1s optimum in a few dozen
// generations.
package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

const (
	nGenes         = 40
	populationSize = 60
	generations    = 200
	seed           = 42
)

func main() {
	gaInstance := &ga.GA{
		Selection: func(p []*ga.Individual, rng *rand.Rand) []*ga.Individual {
			return ga.TournamentSelection(p, 3, rng)
		},
		Crossover:     ga.UniformCrossover,
		Mutation:      ga.BitFlipMutation,
		CrossoverRate: 0.9,
		MutationRate:  1.0 / nGenes,
		Generations:   generations,
		ElitismCount:  2,
		Seed:          seed,
		EarlyStopping: &ga.EarlyStopping{TargetFitness: float64(nGenes), TargetFitnessSet: true},
	}

	init := func(rng *rand.Rand) *ga.Genotype {
		g := ga.NewBinaryGenotype(nGenes)
		for i := range g.Genome {
			g.Genome[i] = byte(rng.Intn(2))
		}
		return g
	}
	eval := func(g *ga.Genotype) *ga.Phenotype {
		fitness := 0.0
		for _, b := range g.Genome {
			if b == 1 {
				fitness++
			}
		}
		return &ga.Phenotype{Fitness: fitness}
	}

	if err := gaInstance.Initialize(populationSize, init, eval); err != nil {
		fmt.Fprintf(os.Stderr, "Initialize failed: %v\n", err)
		os.Exit(1)
	}
	result, err := gaInstance.Evolve(eval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Evolve failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Best fitness: %d / %d\n", int(result.Best.Phenotype.Fitness), nGenes)
	fmt.Printf("Best chromosome: %v\n", result.Best.Genotype.Genome)
	fmt.Printf("Stop reason: %s at generation %d\n", result.StopReason, result.StoppedAtGeneration)
}
