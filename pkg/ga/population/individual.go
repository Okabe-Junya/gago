// Package population provides types and operations for managing individuals and populations in genetic algorithms.
package population

import (
	"fmt"

	"github.com/Okabe-Junya/gago/pkg/ga/encoding"
)

// Phenotype represents the observable traits of an individual, including its fitness value.
type Phenotype struct {
	Features map[string]interface{}
	Fitness  float64
}

// Individual represents an individual in the population, consisting of its genotype and phenotype.
type Individual struct {
	Genotype  *encoding.Genotype
	Phenotype *Phenotype
}

// NewPhenotype creates a new Phenotype with the specified fitness.
func NewPhenotype(fitness float64) *Phenotype {
	return &Phenotype{
		Fitness:  fitness,
		Features: make(map[string]interface{}),
	}
}

// Clone creates a deep copy of the Individual.
func (ind *Individual) Clone() *Individual {
	// Create a clone of the genotype
	genotypeClone := ind.Genotype.Clone()

	// Create a clone of the features map
	featuresClone := make(map[string]interface{})
	for k, v := range ind.Phenotype.Features {
		featuresClone[k] = v
	}

	// Create and return a new Individual with the cloned data
	return &Individual{
		Genotype: genotypeClone,
		Phenotype: &Phenotype{
			Fitness:  ind.Phenotype.Fitness,
			Features: featuresClone,
		},
	}
}

// String returns a string representation of the Individual.
func (ind *Individual) String() string {
	return fmt.Sprintf("Individual{Fitness: %f}", ind.Phenotype.Fitness)
}

// FindBestIndividual finds the individual with the highest fitness in the given population.
func FindBestIndividual(population []*Individual) *Individual {
	if len(population) == 0 {
		return nil
	}

	best := population[0]
	for _, ind := range population {
		if ind.Phenotype.Fitness > best.Phenotype.Fitness {
			best = ind
		}
	}
	return best
}
