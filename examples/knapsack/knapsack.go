// 0/1 Knapsack via GA with a penalty term for capacity violations.
//
// Each gene is 1 if the item is taken, 0 otherwise. The fitness is the total
// value of the chosen items minus a large penalty proportional to how much the
// chosen weight exceeds the capacity.
package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

type item struct {
	name   string
	weight float64
	value  float64
}

var items = []item{
	{"camera", 3.0, 6.0},
	{"tent", 5.0, 9.0},
	{"stove", 2.0, 4.0},
	{"food", 6.0, 12.0},
	{"water", 4.0, 7.0},
	{"map", 1.0, 2.0},
	{"rope", 2.5, 3.5},
	{"first-aid", 1.5, 5.0},
	{"flashlight", 1.0, 1.5},
	{"compass", 0.5, 1.5},
	{"sleeping-bag", 4.5, 8.0},
	{"snacks", 2.0, 5.5},
}

const (
	capacity       = 15.0
	populationSize = 50
	generations    = 300
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
		MutationRate:  1.0 / float64(len(items)),
		Generations:   generations,
		ElitismCount:  2,
		Seed:          seed,
	}

	init := func(rng *rand.Rand) *ga.Genotype {
		g := ga.NewBinaryGenotype(len(items))
		for i := range g.Genome {
			g.Genome[i] = byte(rng.Intn(2))
		}
		return g
	}
	eval := func(g *ga.Genotype) *ga.Phenotype {
		totalWeight, totalValue := 0.0, 0.0
		for i, gene := range g.Genome {
			if gene == 1 {
				totalWeight += items[i].weight
				totalValue += items[i].value
			}
		}
		overflow := 0.0
		if totalWeight > capacity {
			overflow = totalWeight - capacity
		}
		return &ga.Phenotype{Fitness: totalValue - 100.0*overflow}
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

	weight, value := 0.0, 0.0
	picked := []string{}
	for i, gene := range result.Best.Genotype.Genome {
		if gene == 1 {
			picked = append(picked, items[i].name)
			weight += items[i].weight
			value += items[i].value
		}
	}
	fmt.Printf("Picked %d items: %v\n", len(picked), picked)
	fmt.Printf("Total weight: %.1f / %.1f\n", weight, capacity)
	fmt.Printf("Total value:  %.1f\n", value)
}
