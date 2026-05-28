// Minimize the 2D Rastrigin function on the real numbers.
//
// The Rastrigin function f(x) = 10n + sum(xi^2 - 10*cos(2*pi*xi)) has a global
// minimum of 0 at the origin and many local minima, so it is a good stress test
// for the real-valued GA path (GaussianMutation + UniformCrossover).
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Okabe-Junya/gago/pkg/ga"
)

const (
	dim            = 2
	lower          = -5.12
	upper          = 5.12
	populationSize = 80
	generations    = 400
	sigma          = 0.3
	seed           = 42
)

func rastrigin(x []float64) float64 {
	total := 10.0 * float64(len(x))
	for _, xi := range x {
		total += xi*xi - 10.0*math.Cos(2.0*math.Pi*xi)
	}
	return total
}

func main() {
	mins := make([]float64, dim)
	maxs := make([]float64, dim)
	for i := range mins {
		mins[i] = lower
		maxs[i] = upper
	}

	gaInstance := &ga.GA{
		Selection: func(p []*ga.Individual, rng *rand.Rand) []*ga.Individual {
			return ga.TournamentSelection(p, 4, rng)
		},
		Crossover: ga.UniformCrossover,
		Mutation: func(p []*ga.Individual, rate float64, rng *rand.Rand) {
			ga.GaussianMutation(p, rate, sigma, rng)
		},
		CrossoverRate: 0.9,
		MutationRate:  0.3,
		Generations:   generations,
		ElitismCount:  2,
		Seed:          seed,
	}

	init := func(rng *rand.Rand) *ga.Genotype {
		return ga.NewRealGenotype(dim, mins, maxs, rng)
	}
	decode := func(g *ga.Genotype) []float64 {
		out := make([]float64, dim)
		for i := range out {
			v, _ := g.GetRealValue(i)
			out[i] = v
		}
		return out
	}
	eval := func(g *ga.Genotype) *ga.Phenotype {
		return &ga.Phenotype{Fitness: -rastrigin(decode(g))}
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

	point := decode(result.Best.Genotype)
	fmt.Printf("Best point:  %v\n", point)
	fmt.Printf("Rastrigin:   %.6f (global min = 0)\n", -result.Best.Phenotype.Fitness)
}
