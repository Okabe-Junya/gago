// Package ga provides functionalities for implementing genetic algorithms,
// including selection operations for choosing individuals from the population
// to create the next generation.
package ga

import "math/rand"

// TournamentSelection performs tournament selection on the given population.
//
// In tournament selection, a subset of individuals is randomly chosen from the population,
// and the individual with the highest fitness in this subset is selected. This process is repeated
// until the desired number of individuals is selected.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - tournamentSize: the number of individuals to be chosen randomly for each tournament.
//
// Returns:
// - A new population of selected individuals.
func TournamentSelection(population []*Individual, tournamentSize int) []*Individual {
	selected := make([]*Individual, len(population))
	for i := range selected {
		best := population[rand.Intn(len(population))]
		for j := 0; j < tournamentSize-1; j++ {
			contender := population[rand.Intn(len(population))]
			if contender.Phenotype.Fitness > best.Phenotype.Fitness {
				best = contender
			}
		}
		selected[i] = best
	}
	return selected
}

// RouletteWheelSelection performs roulette wheel selection on the given population.
//
// In roulette wheel selection, individuals are selected based on their fitness proportionate to
// the total fitness of the population. This method ensures that individuals with higher fitness
// have a higher chance of being selected.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A new population of selected individuals.
func RouletteWheelSelection(population []*Individual) []*Individual {
	totalFitness := 0.0
	for _, ind := range population {
		totalFitness += ind.Phenotype.Fitness
	}

	selected := make([]*Individual, len(population))
	for i := range selected {
		pick := rand.Float64() * totalFitness
		current := 0.0
		for _, ind := range population {
			current += ind.Phenotype.Fitness
			if current > pick {
				selected[i] = ind
				break
			}
		}
	}
	return selected
}
