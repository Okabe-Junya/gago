// Package ga provides functionalities for implementing genetic algorithms,
// including the definitions and operations related to individuals in the population.
package ga

import (
	"fmt"
	"math/rand"
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
func NewIntegerGenotype(genomeLength int, minValue, maxValue int) *Genotype {
	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: IntegerEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	// Initialize with random values
	for i := range genotype.Genome {
		rangeValue := maxValue - minValue + 1
		genotype.Genome[i] = byte(rand.Intn(rangeValue) + minValue)
		genotype.MinValues[i] = float64(minValue)
		genotype.MaxValues[i] = float64(maxValue)
	}

	return genotype
}

// NewRealGenotype creates a new real-valued genotype with the specified length,
// and values between minValues and maxValues.
func NewRealGenotype(genomeLength int, minValues, maxValues []float64) *Genotype {
	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: RealEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	// Initialize with random values
	for i := range genotype.Genome {
		minIdx := i % len(minValues)
		maxIdx := i % len(maxValues)
		min := minValues[minIdx]
		max := maxValues[maxIdx]
		genotype.MinValues[i] = min
		genotype.MaxValues[i] = max

		// 正規化された値をバイトとして保存
		normalizedValue := rand.Float64()
		genotype.Genome[i] = byte(255 * normalizedValue)
	}

	return genotype
}

// NewPermutationGenotype creates a new permutation genotype with values [0, 1, ..., size-1].
func NewPermutationGenotype(size int) *Genotype {
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
		j := rand.Intn(i + 1)
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
func MutateReal(genotype *Genotype, minValues, maxValues []float64, mutationRate float64, sigma float64) {
	if genotype == nil || len(genotype.Genome) == 0 {
		return
	}

	for i := range genotype.Genome {
		if rand.Float64() < mutationRate {
			// Calculate the valid range for this gene
			rangeValue := maxValues[i%len(maxValues)] - minValues[i%len(minValues)]

			// Add Gaussian noise scaled by sigma and the range
			delta := rand.NormFloat64() * sigma * rangeValue

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

	// クローンGenome
	genomeClone := make([]byte, len(ind.Genotype.Genome))
	copy(genomeClone, ind.Genotype.Genome)

	// クローンMinValues（もし存在すれば）
	var minValuesClone []float64
	if len(ind.Genotype.MinValues) > 0 {
		minValuesClone = make([]float64, len(ind.Genotype.MinValues))
		copy(minValuesClone, ind.Genotype.MinValues)
	}

	// クローンMaxValues（もし存在すれば）
	var maxValuesClone []float64
	if len(ind.Genotype.MaxValues) > 0 {
		maxValuesClone = make([]float64, len(ind.Genotype.MaxValues))
		copy(maxValuesClone, ind.Genotype.MaxValues)
	}

	// クローンFeatures
	var featuresClone []float64
	if len(ind.Phenotype.Features) > 0 {
		featuresClone = make([]float64, len(ind.Phenotype.Features))
		copy(featuresClone, ind.Phenotype.Features)
	}

	return &Individual{
		Genotype: &Genotype{
			Genome:     genomeClone,
			MinValues:  minValuesClone,
			MaxValues:  maxValuesClone,
			GenomeType: ind.Genotype.GenomeType,
		},
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
		return 0, fmt.Errorf("GetBinaryValue called on non-binary encoded genome")
	}
	if position < 0 || position >= len(g.Genome) {
		return 0, fmt.Errorf("position out of bounds: %d", position)
	}
	return int(g.Genome[position]), nil
}

// GetIntegerValue returns the integer value at the specified position.
func (g *Genotype) GetIntegerValue(position int) (int, error) {
	if g.GenomeType != IntegerEncoding {
		return 0, fmt.Errorf("GetIntegerValue called on non-integer encoded genome")
	}
	if position < 0 || position >= len(g.Genome) {
		return 0, fmt.Errorf("position out of bounds: %d", position)
	}

	min := int(g.MinValues[position])
	max := int(g.MaxValues[position])
	rangeValue := max - min + 1

	// Map the byte value (0-255) to the specified integer range
	return min + (int(g.Genome[position])*rangeValue)/256, nil
}

// GetRealValue returns the real value at the specified position.
func (g *Genotype) GetRealValue(position int) (float64, error) {
	if g.GenomeType != RealEncoding {
		return 0, fmt.Errorf("GetRealValue called on non-real encoded genome")
	}
	if position < 0 || position >= len(g.Genome) {
		return 0, fmt.Errorf("position out of bounds: %d", position)
	}

	// Convert the byte (0-255) back to the original range
	normalizedValue := float64(g.Genome[position]) / 255.0
	return g.MinValues[position] + normalizedValue*(g.MaxValues[position]-g.MinValues[position]), nil
}

// GetPermutation returns the entire permutation as a slice of integers.
func (g *Genotype) GetPermutation() ([]int, error) {
	if g.GenomeType != PermutationEncoding {
		return nil, fmt.Errorf("GetPermutation called on non-permutation encoded genome")
	}

	result := make([]int, len(g.Genome))
	for i, v := range g.Genome {
		result[i] = int(v)
	}
	return result, nil
}
