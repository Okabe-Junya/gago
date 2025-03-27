// Package ga provides functionalities for implementing genetic algorithms,
// including mutation operations for introducing genetic diversity in the population.
package ga

import (
	"math/rand"
)

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
		genomeLen := len(ind.Genotype.Genome)
		if genomeLen <= 1 {
			continue
		}

		for i := range ind.Genotype.Genome {
			if rand.Float64() < mutationRate {
				j := rand.Intn(genomeLen - 1)
				if j >= i {
					j++
				}
				ind.Genotype.Genome[i], ind.Genotype.Genome[j] = ind.Genotype.Genome[j], ind.Genotype.Genome[i]
			}
		}
	}
}

// GaussianMutation performs Gaussian mutation on the given population.
// It adds normally distributed random values to each gene.
// This is useful when working with real-valued genomes.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each gene will be mutated.
// - sigma: the standard deviation of the normal distribution.
//
// This function modifies the input population in place.
func GaussianMutation(population []*Individual, mutationRate float64, sigma float64) {
	for _, ind := range population {
		for i := range ind.Genotype.Genome {
			if rand.Float64() < mutationRate {
				// Add Gaussian noise to the gene
				delta := rand.NormFloat64() * sigma

				// Convert to byte with bounds checking
				result := float64(ind.Genotype.Genome[i]) + delta
				if result < 0 {
					result = 0
				} else if result > 255 {
					result = 255
				}

				ind.Genotype.Genome[i] = byte(result)
			}
		}
	}
}

// InversionMutation performs inversion mutation on the given population.
// It selects a random segment of the genome and reverses the order of genes.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each individual will be mutated.
//
// This function modifies the input population in place.
func InversionMutation(population []*Individual, mutationRate float64) {
	for _, ind := range population {
		if rand.Float64() < mutationRate {
			genomeLen := len(ind.Genotype.Genome)
			if genomeLen <= 1 {
				continue
			}

			// Select two random points
			point1 := rand.Intn(genomeLen)
			point2 := rand.Intn(genomeLen)

			// Ensure point1 < point2
			if point1 > point2 {
				point1, point2 = point2, point1
			}

			// Reverse the segment
			for i, j := point1, point2; i < j; i, j = i+1, j-1 {
				ind.Genotype.Genome[i], ind.Genotype.Genome[j] = ind.Genotype.Genome[j], ind.Genotype.Genome[i]
			}
		}
	}
}

// ScrambleMutation performs scramble mutation on the given population.
// It randomly selects a segment of the genome and shuffles the genes within that segment.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each individual will be mutated.
//
// This function modifies the input population in place.
func ScrambleMutation(population []*Individual, mutationRate float64) {
	for _, ind := range population {
		if rand.Float64() < mutationRate {
			genomeLen := len(ind.Genotype.Genome)
			if genomeLen <= 1 {
				continue
			}

			// Select two random points
			point1 := rand.Intn(genomeLen)
			point2 := rand.Intn(genomeLen)

			// Ensure point1 < point2
			if point1 > point2 {
				point1, point2 = point2, point1
			}

			// Create a temporary slice with the segment to shuffle
			segment := make([]byte, point2-point1+1)
			copy(segment, ind.Genotype.Genome[point1:point2+1])

			// Shuffle the segment
			rand.Shuffle(len(segment), func(i, j int) {
				segment[i], segment[j] = segment[j], segment[i]
			})

			// Copy the shuffled segment back
			copy(ind.Genotype.Genome[point1:point2+1], segment)
		}
	}
}

// UniformMutation performs uniform mutation on the given population.
// For each gene that is selected for mutation, it replaces the gene with a random
// value within the specified range.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - mutationRate: the probability with which each gene will be mutated.
// - min: the minimum value for the random replacement.
// - max: the maximum value for the random replacement.
//
// This function modifies the input population in place.
func UniformMutation(population []*Individual, mutationRate float64, min, max byte) {
	for _, ind := range population {
		for i := range ind.Genotype.Genome {
			if rand.Float64() < mutationRate {
				// Replace with a random value in the range [min, max]
				rangeValue := int(max) - int(min) + 1
				ind.Genotype.Genome[i] = min + byte(rand.Intn(rangeValue))
			}
		}
	}
}

// AdaptiveMutation performs mutation with a rate that adapts based on the individual's fitness.
// Individuals with higher fitness have a lower mutation rate, while those with lower fitness
// have a higher mutation rate. This approach balances exploration and exploitation.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - baseMutationRate: the base mutation rate from which to calculate adaptive rates.
// - mutationFunc: the mutation function to apply with the adaptive rates.
// - bestFitness: the fitness of the best individual in the population.
// - worstFitness: the fitness of the worst individual in the population.
//
// This function modifies the input population in place.
func AdaptiveMutation(
	population []*Individual,
	baseMutationRate float64,
	mutationFunc func([]*Individual, float64),
	bestFitness, worstFitness float64,
) {
	fitnessDiff := worstFitness - bestFitness

	for i, ind := range population {
		// Calculate adaptive mutation rate
		adaptiveRate := baseMutationRate

		if fitnessDiff > 0 {
			// Normalize fitness to [0, 1]
			normalizedFitness := (ind.Phenotype.Fitness - bestFitness) / fitnessDiff

			// Adjust mutation rate based on fitness
			// Higher fitness individuals have lower mutation rates
			adaptiveRate = baseMutationRate * (1.0 - 0.5*normalizedFitness)
		}

		// Apply mutation with adaptive rate
		singleIndividual := []*Individual{ind}
		mutationFunc(singleIndividual, adaptiveRate)
		population[i] = singleIndividual[0]
	}
}
