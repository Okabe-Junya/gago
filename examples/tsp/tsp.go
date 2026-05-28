// Travelling Salesman Problem with a permutation encoding.
//
// Each chromosome is a permutation of city indices representing the visit
// order. The tour returns to the starting city after the last visit. We
// minimize total tour length, expressed as a fitness function that returns
// -length (since gago, like most GAs, treats higher fitness as better).
//
// Demonstrates the permutation operators: OrderBasedCrossover (OX1) for
// recombination and SwapMutation for variation.
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

const (
	nCities        = 20
	populationSize = 100
	generations    = 500
	seed           = 42
)

// cities is a fixed set of points in the unit square, seeded for reproducibility.
var cities = func() [][2]float64 {
	r := rand.New(rand.NewSource(2026))
	out := make([][2]float64, nCities)
	for i := range out {
		out[i] = [2]float64{r.Float64(), r.Float64()}
	}
	return out
}()

func tourLength(order []byte) float64 {
	total := 0.0
	n := len(order)
	for i := 0; i < n; i++ {
		x1, y1 := cities[order[i]][0], cities[order[i]][1]
		x2, y2 := cities[order[(i+1)%n]][0], cities[order[(i+1)%n]][1]
		total += math.Hypot(x2-x1, y2-y1)
	}
	return total
}

func main() {
	gaInstance := &ga.GA{
		Selection: func(p []*ga.Individual, rng *rand.Rand) []*ga.Individual {
			return ga.TournamentSelection(p, 4, rng)
		},
		Crossover:     ga.OrderBasedCrossover,
		Mutation:      ga.SwapMutation,
		CrossoverRate: 0.9,
		MutationRate:  0.05,
		Generations:   generations,
		ElitismCount:  2,
		Seed:          seed,
	}

	init := func(rng *rand.Rand) *ga.Genotype {
		return ga.NewPermutationGenotype(nCities, rng)
	}
	eval := func(g *ga.Genotype) *ga.Phenotype {
		return &ga.Phenotype{Fitness: -tourLength(g.Genome)}
	}

	if err := gaInstance.Initialize(populationSize, init, eval); err != nil {
		fmt.Fprintf(os.Stderr, "Initialize failed: %v\n", err)
		os.Exit(1)
	}
	initialBest := -gaInstance.Population.GetBestIndividual().Phenotype.Fitness

	result, err := gaInstance.Evolve(eval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Evolve failed: %v\n", err)
		os.Exit(1)
	}

	finalLength := -result.Best.Phenotype.Fitness
	improvement := (1 - finalLength/initialBest) * 100
	fmt.Printf("Initial best tour: %.4f\n", initialBest)
	fmt.Printf("Final   best tour: %.4f\n", finalLength)
	fmt.Printf("Improvement:       %.1f%%\n", improvement)
	fmt.Printf("Route: %v\n", result.Best.Genotype.Genome)
}
