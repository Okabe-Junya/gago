// Package encoding provides encoding types and operations for genetic algorithms.
// This package implements various genome encodings (binary, integer, real, and permutation)
// and provides methods for manipulating these encodings in a type-safe manner.
package encoding

import (
	"errors"
	"fmt"
	"math/rand"
)

// Common errors for the encoding package
var (
	ErrInvalidGenomeType     = errors.New("invalid genome type")
	ErrInvalidGenomeLength   = errors.New("invalid genome length")
	ErrInvalidGenomePosition = errors.New("invalid genome position")
	ErrValueRangeMismatch    = errors.New("min/max values length must match genome length")
)

// GenomeType represents the type of genome encoding.
type GenomeType int

const (
	BinaryEncoding GenomeType = iota
	IntegerEncoding
	RealEncoding
	PermutationEncoding
)

// String returns a string representation of the GenomeType.
func (gt GenomeType) String() string {
	switch gt {
	case BinaryEncoding:
		return "Binary"
	case IntegerEncoding:
		return "Integer"
	case RealEncoding:
		return "Real"
	case PermutationEncoding:
		return "Permutation"
	default:
		return "Unknown"
	}
}

// Genotype represents the genetic makeup of an individual, encoded as a sequence of bytes.
type Genotype struct {
	Genome     []byte
	MinValues  []float64
	MaxValues  []float64
	GenomeType GenomeType
}

// NewBinaryGenotype creates a new binary-encoded Genotype with the specified genome length.
// Binary encoding represents genes as a sequence of 0s and 1s.
//
// Parameters:
// - genomeLength: the length of the genome to be created.
//
// Returns:
// - A pointer to the newly created Genotype.
// - Panics if genomeLength is less than or equal to 0.
func NewBinaryGenotype(genomeLength int) *Genotype {
	if genomeLength <= 0 {
		panic(fmt.Errorf("%w: %d", ErrInvalidGenomeLength, genomeLength))
	}
	return &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: BinaryEncoding,
	}
}

// NewIntegerGenotype creates a new integer-encoded Genotype with the specified range.
// Integer encoding represents genes as integer values within a specified range.
//
// Parameters:
// - genomeLength: the length of the genome to be created.
// - minValue: the minimum value for each gene.
// - maxValue: the maximum value for each gene.
//
// Returns:
// - A pointer to the newly created Genotype with random integer values between minValue and maxValue.
// - Panics if genomeLength is less than or equal to 0.
func NewIntegerGenotype(genomeLength int, minValue, maxValue int) *Genotype {
	if genomeLength <= 0 {
		panic(fmt.Errorf("%w: %d", ErrInvalidGenomeLength, genomeLength))
	}
	if minValue > maxValue {
		minValue, maxValue = maxValue, minValue
	}

	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: IntegerEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	for i := 0; i < genomeLength; i++ {
		genotype.MinValues[i] = float64(minValue)
		genotype.MaxValues[i] = float64(maxValue)
		genotype.Genome[i] = byte(minValue + rand.Intn(maxValue-minValue+1))
	}

	return genotype
}

// NewRealGenotype creates a new real-encoded Genotype with the specified range.
// Real encoding represents genes as real numbers within specified ranges.
//
// Parameters:
// - genomeLength: the length of the genome to be created.
// - minValues: slice of minimum values for each gene.
// - maxValues: slice of maximum values for each gene.
//
// Returns:
// - A pointer to the newly created Genotype with random real values between minValues and maxValues.
// - Panics if genomeLength is less than or equal to 0 or if minValues/maxValues lengths don't match genomeLength.
func NewRealGenotype(genomeLength int, minValues, maxValues []float64) *Genotype {
	if genomeLength <= 0 {
		panic(fmt.Errorf("%w: %d", ErrInvalidGenomeLength, genomeLength))
	}
	if len(minValues) != genomeLength || len(maxValues) != genomeLength {
		panic(ErrValueRangeMismatch)
	}

	genotype := &Genotype{
		Genome:     make([]byte, genomeLength),
		GenomeType: RealEncoding,
		MinValues:  make([]float64, genomeLength),
		MaxValues:  make([]float64, genomeLength),
	}

	for i := 0; i < genomeLength; i++ {
		// Ensure min is actually less than max
		min, max := minValues[i], maxValues[i]
		if min > max {
			min, max = max, min
		}

		genotype.MinValues[i] = min
		genotype.MaxValues[i] = max

		// Initialize with random values scaled to the appropriate range
		normalizedValue := rand.Float64()
		// スケーリングされた値を保存（エンコードされた値として）
		genotype.Genome[i] = byte(255 * normalizedValue)
	}

	return genotype
}

// NewPermutationGenotype creates a new permutation-encoded Genotype.
// Permutation encoding represents genes as a sequence of unique integers,
// useful for problems like the traveling salesman problem.
//
// Parameters:
// - size: the size of the permutation (number of elements to be permuted).
//
// Returns:
// - A pointer to the newly created Genotype with a random permutation of integers from 0 to size-1.
// - Panics if size is less than or equal to 0.
func NewPermutationGenotype(size int) *Genotype {
	if size <= 0 {
		panic(fmt.Errorf("%w: %d", ErrInvalidGenomeLength, size))
	}

	genotype := &Genotype{
		Genome:     make([]byte, size),
		GenomeType: PermutationEncoding,
	}

	// Initialize with values 0 to size-1
	for i := 0; i < size; i++ {
		genotype.Genome[i] = byte(i)
	}

	// Shuffle to create a random permutation
	rand.Shuffle(size, func(i, j int) {
		genotype.Genome[i], genotype.Genome[j] = genotype.Genome[j], genotype.Genome[i]
	})

	return genotype
}

// checkBounds verifies that a position is within the valid range of the genome.
//
// Parameters:
// - position: the index to check.
//
// Returns:
// - An error if the position is out of bounds, nil otherwise.
func (g *Genotype) checkBounds(position int) error {
	if position < 0 || position >= len(g.Genome) {
		return fmt.Errorf("%w: %d (length: %d)", ErrInvalidGenomePosition, position, len(g.Genome))
	}
	return nil
}

// GetBinaryValue returns the binary value (0 or 1) at the specified position.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The binary value at the specified position.
// - An error if the position is invalid or if the genome is not binary-encoded.
func (g *Genotype) GetBinaryValue(position int) (int, error) {
	if g.GenomeType != BinaryEncoding {
		return 0, fmt.Errorf("%w: expected %s, got %s", ErrInvalidGenomeType, BinaryEncoding, g.GenomeType)
	}

	if err := g.checkBounds(position); err != nil {
		return 0, err
	}

	return int(g.Genome[position]), nil
}

// GetBinaryValueUnsafe returns the binary value without bounds or type checking.
// This method should only be used when performance is critical and you are certain
// that the position is valid and the encoding type is correct.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The binary value at the specified position.
func (g *Genotype) GetBinaryValueUnsafe(position int) int {
	return int(g.Genome[position])
}

// GetIntegerValue returns the integer value at the specified position.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The integer value at the specified position.
// - An error if the position is invalid or if the genome is not integer-encoded.
func (g *Genotype) GetIntegerValue(position int) (int, error) {
	if g.GenomeType != IntegerEncoding {
		return 0, fmt.Errorf("%w: expected %s, got %s", ErrInvalidGenomeType, IntegerEncoding, g.GenomeType)
	}

	if err := g.checkBounds(position); err != nil {
		return 0, err
	}

	min := int(g.MinValues[position])
	max := int(g.MaxValues[position])
	range_ := max - min + 1

	// Map the byte value (0-255) to the specified integer range
	return min + (int(g.Genome[position])*range_)/256, nil
}

// GetIntegerValueUnsafe returns the integer value without bounds or type checking.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The integer value at the specified position.
func (g *Genotype) GetIntegerValueUnsafe(position int) int {
	min := int(g.MinValues[position])
	max := int(g.MaxValues[position])
	range_ := max - min + 1

	// Map the byte value (0-255) to the specified integer range
	return min + (int(g.Genome[position])*range_)/256
}

// GetRealValue returns the real value at the specified position.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The real value at the specified position.
// - An error if the position is invalid or if the genome is not real-encoded.
func (g *Genotype) GetRealValue(position int) (float64, error) {
	if g.GenomeType != RealEncoding {
		return 0, fmt.Errorf("%w: expected %s, got %s", ErrInvalidGenomeType, RealEncoding, g.GenomeType)
	}

	if err := g.checkBounds(position); err != nil {
		return 0, err
	}

	// Convert the byte (0-255) back to the original range
	normalizedValue := float64(g.Genome[position]) / 255.0
	return g.MinValues[position] + normalizedValue*(g.MaxValues[position]-g.MinValues[position]), nil
}

// GetRealValueUnsafe returns the real value without bounds or type checking.
//
// Parameters:
// - position: the index in the genome to read.
//
// Returns:
// - The real value at the specified position.
func (g *Genotype) GetRealValueUnsafe(position int) float64 {
	// Convert the byte (0-255) back to the original range
	normalizedValue := float64(g.Genome[position]) / 255.0
	return g.MinValues[position] + normalizedValue*(g.MaxValues[position]-g.MinValues[position])
}

// SetRealValue sets a real value at the specified position.
//
// Parameters:
// - position: the index in the genome to write.
// - value: the value to set.
//
// Returns:
// - An error if the position is invalid or if the genome is not real-encoded.
func (g *Genotype) SetRealValue(position int, value float64) error {
	if g.GenomeType != RealEncoding {
		return fmt.Errorf("%w: expected %s, got %s", ErrInvalidGenomeType, RealEncoding, g.GenomeType)
	}

	if err := g.checkBounds(position); err != nil {
		return err
	}

	// Clamp the value to the allowed range
	if value < g.MinValues[position] {
		value = g.MinValues[position]
	}
	if value > g.MaxValues[position] {
		value = g.MaxValues[position]
	}

	// Convert the value back to a byte (0-255)
	normalizedValue := (value - g.MinValues[position]) / (g.MaxValues[position] - g.MinValues[position])
	g.Genome[position] = byte(normalizedValue * 255)
	return nil
}

// SetRealValueUnsafe sets a real value without bounds or type checking.
//
// Parameters:
// - position: the index in the genome to write.
// - value: the value to set.
func (g *Genotype) SetRealValueUnsafe(position int, value float64) {
	// Clamp the value to the allowed range
	if value < g.MinValues[position] {
		value = g.MinValues[position]
	}
	if value > g.MaxValues[position] {
		value = g.MaxValues[position]
	}

	// Convert the value back to a byte (0-255)
	normalizedValue := (value - g.MinValues[position]) / (g.MaxValues[position] - g.MinValues[position])
	g.Genome[position] = byte(normalizedValue * 255)
}

// GetPermutation returns the entire permutation as a slice of integers.
//
// Returns:
// - A slice containing the permutation.
// - An error if the genome is not permutation-encoded.
func (g *Genotype) GetPermutation() ([]int, error) {
	if g.GenomeType != PermutationEncoding {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrInvalidGenomeType, PermutationEncoding, g.GenomeType)
	}

	result := make([]int, len(g.Genome))
	for i, v := range g.Genome {
		result[i] = int(v)
	}
	return result, nil
}

// GetPermutationUnsafe returns the permutation without type checking.
//
// Returns:
// - A slice containing the permutation.
func (g *Genotype) GetPermutationUnsafe() []int {
	result := make([]int, len(g.Genome))
	for i, v := range g.Genome {
		result[i] = int(v)
	}
	return result
}

// Clone creates a deep copy of the Genotype.
//
// Returns:
// - A pointer to a new Genotype with identical contents.
func (g *Genotype) Clone() *Genotype {
	genomeClone := make([]byte, len(g.Genome))
	copy(genomeClone, g.Genome)

	minValuesClone := make([]float64, len(g.MinValues))
	maxValuesClone := make([]float64, len(g.MaxValues))
	copy(minValuesClone, g.MinValues)
	copy(maxValuesClone, g.MaxValues)

	return &Genotype{
		Genome:     genomeClone,
		GenomeType: g.GenomeType,
		MinValues:  minValuesClone,
		MaxValues:  maxValuesClone,
	}
}

// String returns a string representation of the Genotype.
//
// Returns:
// - A string describing the genotype's type and length.
func (g *Genotype) String() string {
	return fmt.Sprintf("Genotype{Type: %s, Length: %d}", g.GenomeType, len(g.Genome))
}
