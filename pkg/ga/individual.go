package ga

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

var (
	errNotBinaryGenome      = errors.New("GetBinaryValue called on non-binary encoded genome")
	errNotIntegerGenome     = errors.New("GetIntegerValue called on non-integer encoded genome")
	errNotRealGenome        = errors.New("GetRealValue called on non-real encoded genome")
	errNotPermutationGenome = errors.New("GetPermutation called on non-permutation encoded genome")
	errPositionOutOfBounds  = errors.New("position out of bounds")
)

// GenomeType represents the type of genome encoding.
type GenomeType int

const (
	// BinaryEncoding represents a binary encoding of the genome.
	BinaryEncoding GenomeType = iota
	// IntegerEncoding represents an integer encoding of the genome.
	IntegerEncoding
	// RealEncoding represents a real-valued encoding of the genome.
	RealEncoding
	// PermutationEncoding represents a permutation encoding of the genome.
	PermutationEncoding
)

// Genotype represents the genetic makeup of an individual.
type Genotype struct {
	Genome     []byte
	MinValues  []float64
	MaxValues  []float64
	GenomeType GenomeType
}

// Phenotype represents the expressed traits of an individual.
type Phenotype struct {
	Features []float64
	Fitness  float64
}

// Individual represents a solution in the population.
type Individual struct {
	Genotype  *Genotype
	Phenotype *Phenotype
}

// NewBinaryGenotype creates a new binary genotype with the specified length.
func NewBinaryGenotype(genomeLength int) *Genotype {
	return &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: BinaryEncoding,
	}
}

// NewIntegerGenotype creates a new integer genotype with the specified length,
// and values between minValue and maxValue.
func NewIntegerGenotype(genomeLength int, minValue, maxValue int, rng *rand.Rand) *Genotype {
	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: IntegerEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	// Genes hold a normalized 0-255 slot that GetIntegerValue maps onto [minValue, maxValue].
	rng.Read(genotype.Genome)

	for i := range genotype.Genome {
		genotype.MinValues[i] = float64(minValue)
		genotype.MaxValues[i] = float64(maxValue)
	}

	return genotype
}

// NewRealGenotype creates a new real-valued genotype with the specified length,
// and values between minValues and maxValues.
func NewRealGenotype(genomeLength int, minValues, maxValues []float64, rng *rand.Rand) *Genotype {
	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: RealEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	for i := range genotype.Genome {
		genotype.MinValues[i] = minValues[i%len(minValues)]
		genotype.MaxValues[i] = maxValues[i%len(maxValues)]
		genotype.Genome[i] = byte(math.MaxUint8 * rng.Float64())
	}

	return genotype
}

// NewPermutationGenotype creates a new permutation genotype with values [0, 1, ..., size-1].
func NewPermutationGenotype(size int, rng *rand.Rand) *Genotype {
	genotype := &Genotype{
		Genome:     make([]byte, size),
		GenomeType: PermutationEncoding,
	}

	// Initialize with sequential values
	for i := range genotype.Genome {
		genotype.Genome[i] = byte(i)
	}

	// Shuffle the genotype to create a random permutation
	for i := size - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		genotype.Genome[i], genotype.Genome[j] = genotype.Genome[j], genotype.Genome[i]
	}

	return genotype
}

// NewPhenotype creates a new phenotype with the specified fitness.
func NewPhenotype(fitness float64) *Phenotype {
	return &Phenotype{
		Fitness: fitness,
	}
}

// MutateReal mutates a real-valued genotype by adding Gaussian noise.
func MutateReal(
	genotype *Genotype,
	minValues, maxValues []float64,
	mutationRate float64,
	sigma float64,
	rng *rand.Rand,
) {
	if genotype == nil || len(genotype.Genome) == 0 {
		return
	}

	for i := range genotype.Genome {
		if rng.Float64() < mutationRate {
			// Calculate the valid range for this gene
			rangeValue := maxValues[i%len(maxValues)] - minValues[i%len(minValues)]

			// Add Gaussian noise scaled by sigma and the range
			delta := rng.NormFloat64() * sigma * rangeValue

			// Apply the mutation and clamp to valid range
			newValue := float64(genotype.Genome[i]) + delta
			if newValue < float64(minValues[i%len(minValues)]) {
				newValue = float64(minValues[i%len(minValues)])
			} else if newValue > float64(maxValues[i%len(maxValues)]) {
				newValue = float64(maxValues[i%len(maxValues)])
			}

			genotype.Genome[i] = byte(newValue)
		}
	}
}

// Clone creates a deep copy of the Individual.
func (ind *Individual) Clone() *Individual {
	if ind == nil || ind.Genotype == nil || ind.Phenotype == nil {
		return nil
	}

	return &Individual{
		Genotype: &Genotype{
			Genome:     append([]byte(nil), ind.Genotype.Genome...),
			MinValues:  cloneNonEmpty(ind.Genotype.MinValues),
			MaxValues:  cloneNonEmpty(ind.Genotype.MaxValues),
			GenomeType: ind.Genotype.GenomeType,
		},
		Phenotype: &Phenotype{
			Fitness:  ind.Phenotype.Fitness,
			Features: cloneNonEmpty(ind.Phenotype.Features),
		},
	}
}

// String returns a string representation of the Individual.
func (ind *Individual) String() string {
	return fmt.Sprintf("Individual{Fitness: %f}", ind.Phenotype.Fitness)
}

// findBestIndividual finds the individual with the highest fitness in the given population.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A pointer to the individual with the highest fitness.
func findBestIndividual(population []*Individual) *Individual {
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

// GetBinaryValue returns the binary value (0 or 1) at the specified position.
func (g *Genotype) GetBinaryValue(position int) (int, error) {
	if g.GenomeType != BinaryEncoding {
		return 0, errNotBinaryGenome
	}

	if err := g.checkPosition(position); err != nil {
		return 0, err
	}

	return int(g.Genome[position]), nil
}

// GetIntegerValue returns the integer value at the specified position.
func (g *Genotype) GetIntegerValue(position int) (int, error) {
	if g.GenomeType != IntegerEncoding {
		return 0, errNotIntegerGenome
	}

	if err := g.checkPosition(position); err != nil {
		return 0, err
	}

	lo := int(g.MinValues[position])
	hi := int(g.MaxValues[position])

	return lo + (int(g.Genome[position])*(hi-lo+1))/(math.MaxUint8+1), nil
}

// GetRealValue returns the real value at the specified position.
func (g *Genotype) GetRealValue(position int) (float64, error) {
	if g.GenomeType != RealEncoding {
		return 0, errNotRealGenome
	}

	if err := g.checkPosition(position); err != nil {
		return 0, err
	}

	normalizedValue := float64(g.Genome[position]) / math.MaxUint8

	return g.MinValues[position] + normalizedValue*(g.MaxValues[position]-g.MinValues[position]), nil
}

// GetPermutation returns the entire permutation as a slice of integers.
func (g *Genotype) GetPermutation() ([]int, error) {
	if g.GenomeType != PermutationEncoding {
		return nil, errNotPermutationGenome
	}

	result := make([]int, len(g.Genome))
	for i, v := range g.Genome {
		result[i] = int(v)
	}
	return result, nil
}

func (g *Genotype) checkPosition(position int) error {
	if position < 0 || position >= len(g.Genome) {
		return fmt.Errorf("%w: %d", errPositionOutOfBounds, position)
	}

	return nil
}
