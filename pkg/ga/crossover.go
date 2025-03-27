// Package ga provides functionalities for implementing genetic algorithms,
// including crossover operations for generating offspring from parent individuals.
package ga

import (
	"math/rand"
	"sort"
)

// SinglePointCrossover performs a single-point crossover on the given population.
//
// In single-point crossover, a random crossover point is selected, and the
// offspring are created by exchanging the segments of the parent individuals' genomes
// after this point.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func SinglePointCrossover(population []*Individual, crossoverRate float64) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rand.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			point := rand.Intn(len(parent1.Genome))
			child1 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			child2 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			copy(child1.Genome[:point], parent1.Genome[:point])
			copy(child1.Genome[point:], parent2.Genome[point:])
			copy(child2.Genome[:point], parent2.Genome[:point])
			copy(child2.Genome[point:], parent1.Genome[point:])
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// UniformCrossover performs a uniform crossover on the given population.
//
// In uniform crossover, each gene from the parent individuals is independently
// chosen with a 50% probability to be included in the offspring. This allows
// for more genetic diversity in the offspring compared to single-point crossover.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func UniformCrossover(population []*Individual, crossoverRate float64) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rand.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			child1 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			child2 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			for j := range parent1.Genome {
				if rand.Float64() < 0.5 {
					child1.Genome[j] = parent1.Genome[j]
					child2.Genome[j] = parent2.Genome[j]
				} else {
					child1.Genome[j] = parent2.Genome[j]
					child2.Genome[j] = parent1.Genome[j]
				}
			}
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// MultiPointCrossover performs a multi-point crossover on the given population.
//
// In multi-point crossover, multiple crossover points are selected, and the
// offspring are created by exchanging segments of the parent individuals' genomes
// between these points. This allows for more genetic material exchange compared to
// single-point crossover.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
// - numPoints: the number of crossover points to use.
//
// Returns:
// - A new population of offspring generated from the input population.
func MultiPointCrossover(population []*Individual, crossoverRate float64, numPoints int) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rand.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype

			// Generate crossover points
			genomeLength := len(parent1.Genome)
			if numPoints > genomeLength-1 {
				numPoints = genomeLength - 1
			}

			points := make([]int, numPoints)
			for j := 0; j < numPoints; j++ {
				points[j] = rand.Intn(genomeLength)
			}
			sort.Ints(points)

			// Create children
			child1 := &Genotype{Genome: make([]byte, genomeLength)}
			child2 := &Genotype{Genome: make([]byte, genomeLength)}

			// Start with parent1's genes for child1 and parent2's genes for child2
			swap := false
			start := 0

			for j := 0; j < numPoints; j++ {
				end := points[j]

				if !swap {
					copy(child1.Genome[start:end], parent1.Genome[start:end])
					copy(child2.Genome[start:end], parent2.Genome[start:end])
				} else {
					copy(child1.Genome[start:end], parent2.Genome[start:end])
					copy(child2.Genome[start:end], parent1.Genome[start:end])
				}

				swap = !swap
				start = end
			}

			// Handle the last segment
			if !swap {
				copy(child1.Genome[start:], parent1.Genome[start:])
				copy(child2.Genome[start:], parent2.Genome[start:])
			} else {
				copy(child1.Genome[start:], parent2.Genome[start:])
				copy(child2.Genome[start:], parent1.Genome[start:])
			}

			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// OrderBasedCrossover performs an order-based crossover on the given population.
// This type of crossover is useful for permutation problems (like TSP).
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func OrderBasedCrossover(population []*Individual, crossoverRate float64) []*Individual {
	offspring := make([]*Individual, len(population))

	for i := 0; i < len(population)/2; i++ {
		if rand.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			genomeLength := len(parent1.Genome)

			// Select a random subset of positions
			start := rand.Intn(genomeLength)
			length := rand.Intn(genomeLength - start + 1)
			end := start + length

			// Create children
			child1 := &Genotype{Genome: make([]byte, genomeLength)}
			child2 := &Genotype{Genome: make([]byte, genomeLength)}

			// Initialize with -1 to mark as unfilled
			for j := 0; j < genomeLength; j++ {
				child1.Genome[j] = 255 // Sentinel value
				child2.Genome[j] = 255 // Sentinel value
			}

			// Copy the selected segment from parent to child
			copy(child1.Genome[start:end], parent1.Genome[start:end])
			copy(child2.Genome[start:end], parent2.Genome[start:end])

			// Fill the remaining positions in the order they appear in the other parent
			fillOrderBasedOffspring(parent2.Genome, child1.Genome, start, end)
			fillOrderBasedOffspring(parent1.Genome, child2.Genome, start, end)

			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}

	return offspring
}

// fillOrderBasedOffspring fills the remaining positions in a child genome for order-based crossover.
func fillOrderBasedOffspring(parentGenome, childGenome []byte, start, end int) {
	childIdx := 0

	// Skip positions that are already filled
	if childIdx == start {
		childIdx = end
	}

	for _, gene := range parentGenome {
		// Check if this gene is already in the child
		alreadyExists := false
		for j := start; j < end; j++ {
			if childGenome[j] == gene {
				alreadyExists = true
				break
			}
		}

		if !alreadyExists {
			childGenome[childIdx] = gene
			childIdx++
			if childIdx == start {
				childIdx = end
			}
		}
	}
}
