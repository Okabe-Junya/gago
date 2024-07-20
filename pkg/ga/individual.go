// Package ga provides functionalities for implementing genetic algorithms,
// including the definitions and operations related to individuals in the population.
package ga

// Genotype represents the genetic makeup of an individual, encoded as a sequence of bytes.
type Genotype struct {
	Genome []byte
}

// Phenotype represents the observable traits of an individual, including its fitness value.
type Phenotype struct {
	Fitness float64
}

// Individual represents an individual in the population, consisting of its genotype and phenotype.
type Individual struct {
	Genotype  *Genotype
	Phenotype *Phenotype
}

// NewGenotype creates a new Genotype with the specified genome length.
//
// Parameters:
// - genomeLength: the length of the genome to be created.
//
// Returns:
// - A pointer to the newly created Genotype.
func NewGenotype(genomeLength int) *Genotype {
	return &Genotype{
		Genome: make([]byte, genomeLength),
	}
}

// findBestIndividual finds the individual with the highest fitness in the given population.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A pointer to the individual with the highest fitness.
func findBestIndividual(population []*Individual) *Individual {
	best := population[0]
	for _, ind := range population {
		if ind.Phenotype.Fitness > best.Phenotype.Fitness {
			best = ind
		}
	}
	return best
}
