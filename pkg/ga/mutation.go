// Package ga provides functionalities for implementing genetic algorithms,
// including mutation operations for introducing genetic diversity in the population.
package ga

import "math/rand"

// BitFlipMutation performs bit-flip mutation on the given population.
//
// In bit-flip mutation, each bit (or gene) in the individual's genome is
// independently flipped with a certain probability, known as the mutation rate.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each gene will be mutated.
//
// This function modifies the input population in place.
func BitFlipMutation(population []*Individual, mutationRate float64) {
	for _, ind := range population {
		for i := range ind.Genotype.Genome {
			if rand.Float64() < mutationRate {
				ind.Genotype.Genome[i] = 1 - ind.Genotype.Genome[i]
			}
		}
	}
}

// SwapMutation performs swap mutation on the given population.
//
// In swap mutation, two genes in the individual's genome are randomly selected
// and swapped with a certain probability, known as the mutation rate.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each gene will be considered for swapping.
//
// This function modifies the input population in place.
func SwapMutation(population []*Individual, mutationRate float64) {
	for _, ind := range population {
		for i := range ind.Genotype.Genome {
			if rand.Float64() < mutationRate {
				j := rand.Intn(len(ind.Genotype.Genome))
				ind.Genotype.Genome[i], ind.Genotype.Genome[j] = ind.Genotype.Genome[j], ind.Genotype.Genome[i]
			}
		}
	}
}
