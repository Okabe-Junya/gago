// Package ga provides functionalities for implementing genetic algorithms,
// including crossover operations for generating offspring from parent individuals.
package ga

import "math/rand"

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
