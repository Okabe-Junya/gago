// Package population provides types and operations for managing individuals and populations in genetic algorithms.
package population

import (
	"math"
	"sort"
)

// Population represents a collection of individuals in a genetic algorithm.
type Population struct {
	Statistics  *Statistics
	Individuals []*Individual
}

// Statistics stores statistical information about a population.
type Statistics struct {
	BestFitness    float64
	WorstFitness   float64
	AverageFitness float64
	Diversity      float64
}

// NewPopulation creates a new population with the given size using the initialization function.
func NewPopulation(size int, initFunc func() *Individual) *Population {
	pop := &Population{
		Individuals: make([]*Individual, size),
		Statistics:  &Statistics{},
	}
	for i := 0; i < size; i++ {
		pop.Individuals[i] = initFunc()
	}
	return pop
}

// CalculateStatistics calculates statistical information about the population.
func (p *Population) CalculateStatistics() {
	if len(p.Individuals) == 0 {
		return
	}
	bestFitness := p.Individuals[0].Phenotype.Fitness
	worstFitness := p.Individuals[0].Phenotype.Fitness
	totalFitness := 0.0
	for _, ind := range p.Individuals {
		fitness := ind.Phenotype.Fitness
		if fitness > bestFitness {
			bestFitness = fitness
		}
		if fitness < worstFitness {
			worstFitness = fitness
		}
		totalFitness += fitness
	}
	averageFitness := totalFitness / float64(len(p.Individuals))

	// Calculate genetic diversity as standard deviation of fitness values
	sumSquaredDiffs := 0.0
	for _, ind := range p.Individuals {
		diff := ind.Phenotype.Fitness - averageFitness
		sumSquaredDiffs += diff * diff
	}
	diversity := math.Sqrt(sumSquaredDiffs / float64(len(p.Individuals)))

	p.Statistics = &Statistics{
		BestFitness:    bestFitness,
		WorstFitness:   worstFitness,
		AverageFitness: averageFitness,
		Diversity:      diversity,
	}
}

// SortByFitness sorts the population by fitness in descending order.
func (p *Population) SortByFitness() {
	sort.Slice(p.Individuals, func(i, j int) bool {
		return p.Individuals[i].Phenotype.Fitness > p.Individuals[j].Phenotype.Fitness
	})
}

// GetBestIndividual returns the individual with the highest fitness.
func (p *Population) GetBestIndividual() *Individual {
	return FindBestIndividual(p.Individuals)
}

// GetWorstIndividual returns the individual with the lowest fitness.
func (p *Population) GetWorstIndividual() *Individual {
	if len(p.Individuals) == 0 {
		return nil
	}
	worst := p.Individuals[0]
	for _, ind := range p.Individuals {
		if ind.Phenotype.Fitness < worst.Phenotype.Fitness {
			worst = ind
		}
	}
	return worst
}

// Replace replaces an individual at the specified index with a new individual.
func (p *Population) Replace(index int, individual *Individual) {
	if index >= 0 && index < len(p.Individuals) {
		p.Individuals[index] = individual
	}
}

// Size returns the number of individuals in the population.
func (p *Population) Size() int {
	return len(p.Individuals)
}
