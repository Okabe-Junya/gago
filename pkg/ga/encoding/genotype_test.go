package encoding

import (
	"reflect"
	"testing"
)

func TestNewBinaryGenotype(t *testing.T) {
	cases := []struct {
		genomeLength   int
		expectedLength int
	}{
		{genomeLength: 5, expectedLength: 5},
		{genomeLength: 10, expectedLength: 10},
	}

	for _, tc := range cases {
		genotype := NewBinaryGenotype(tc.genomeLength)

		if len(genotype.Genome) != tc.expectedLength {
			t.Fatalf("Expected genome length %d, but got %d", tc.expectedLength, len(genotype.Genome))
		}

		if genotype.GenomeType != BinaryEncoding {
			t.Fatalf("Expected genome type BinaryEncoding, but got %v", genotype.GenomeType)
		}
	}
}

func TestNewIntegerGenotype(t *testing.T) {
	genomeLength := 5
	minValue := 10
	maxValue := 20

	genotype := NewIntegerGenotype(genomeLength, minValue, maxValue)

	if len(genotype.Genome) != genomeLength {
		t.Fatalf("Expected genome length %d, but got %d", genomeLength, len(genotype.Genome))
	}

	if genotype.GenomeType != IntegerEncoding {
		t.Fatalf("Expected genome type IntegerEncoding, but got %v", genotype.GenomeType)
	}

	for i := 0; i < genomeLength; i++ {
		value, err := genotype.GetIntegerValue(i)
		if err != nil {
			t.Errorf("Failed to get integer value at position %d: %v", i, err)
		}
		if value < minValue || value > maxValue {
			t.Errorf("Integer value %d at position %d is outside the expected range [%d, %d]",
				value, i, minValue, maxValue)
		}
	}
}

func TestNewRealGenotype(t *testing.T) {
	genomeLength := 5
	minValues := []float64{0.0, 1.0, 2.0, 3.0, 4.0}
	maxValues := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	genotype := NewRealGenotype(genomeLength, minValues, maxValues)

	if len(genotype.Genome) != genomeLength {
		t.Fatalf("Expected genome length %d, but got %d", genomeLength, len(genotype.Genome))
	}

	if genotype.GenomeType != RealEncoding {
		t.Fatalf("Expected genome type RealEncoding, but got %v", genotype.GenomeType)
	}

	for i := 0; i < genomeLength; i++ {
		value, err := genotype.GetRealValue(i)
		if err != nil {
			t.Errorf("Failed to get real value at position %d: %v", i, err)
		}
		if value < minValues[i] || value > maxValues[i] {
			t.Errorf("Real value %f at position %d is outside the expected range [%f, %f]",
				value, i, minValues[i], maxValues[i])
		}
	}
}

func TestNewPermutationGenotype(t *testing.T) {
	size := 5

	genotype := NewPermutationGenotype(size)

	if len(genotype.Genome) != size {
		t.Fatalf("Expected genome length %d, but got %d", size, len(genotype.Genome))
	}

	if genotype.GenomeType != PermutationEncoding {
		t.Fatalf("Expected genome type PermutationEncoding, but got %v", genotype.GenomeType)
	}

	// Check that all values from 0 to size-1 are present (permutation)
	permutation, err := genotype.GetPermutation()
	if err != nil {
		t.Fatalf("Failed to get permutation: %v", err)
	}

	valueCounts := make(map[int]int)

	for _, value := range permutation {
		valueCounts[value]++
	}

	for i := 0; i < size; i++ {
		if count, exists := valueCounts[i]; !exists || count != 1 {
			t.Errorf("Value %d appears %d times in permutation, expected exactly once", i, count)
		}
	}
}

func TestGetSetRealValue(t *testing.T) {
	genomeLength := 3
	minValues := []float64{0.0, 1.0, 2.0}
	maxValues := []float64{10.0, 11.0, 12.0}

	genotype := NewRealGenotype(genomeLength, minValues, maxValues)

	// Test setting and getting values
	testValues := []float64{5.0, 6.0, 7.0}
	for i, value := range testValues {
		err := genotype.SetRealValue(i, value)
		if err != nil {
			t.Errorf("Failed to set real value %f at position %d: %v", value, i, err)
		}

		retrieved, err := genotype.GetRealValue(i)
		if err != nil {
			t.Errorf("Failed to get real value at position %d: %v", i, err)
		}

		// Check if the retrieved value is close enough to what we set
		// Due to encoding as bytes, there may be small rounding errors
		if retrieved < value-0.1 || retrieved > value+0.1 {
			t.Errorf("Expected real value close to %f, got %f", value, retrieved)
		}
	}

	// Test clamping of values outside range
	if err := genotype.SetRealValue(0, -1.0); err != nil { // Below min
		t.Errorf("Failed to set real value -1.0 at position 0: %v", err)
	}

	v0, err := genotype.GetRealValue(0)
	if err != nil {
		t.Errorf("Failed to get real value at position 0: %v", err)
	}
	if v0 != minValues[0] {
		t.Errorf("Value should be clamped to min, got %f", v0)
	}

	if err := genotype.SetRealValue(1, 20.0); err != nil { // Above max
		t.Errorf("Failed to set real value 20.0 at position 1: %v", err)
	}

	v1, err := genotype.GetRealValue(1)
	if err != nil {
		t.Errorf("Failed to get real value at position 1: %v", err)
	}
	if v1 != maxValues[1] {
		t.Errorf("Value should be clamped to max, got %f", v1)
	}
}

func TestGenotypeClone(t *testing.T) {
	original := NewBinaryGenotype(5)
	original.Genome = []byte{1, 0, 1, 0, 1}

	clone := original.Clone()

	// Check that the clone has the same values
	if !reflect.DeepEqual(original.Genome, clone.Genome) {
		t.Errorf("Clone genome values don't match original: %v vs %v",
			original.Genome, clone.Genome)
	}

	// Check that modifying the clone doesn't affect original
	clone.Genome[0] = 0
	if original.Genome[0] != 1 {
		t.Error("Modifying clone affected original genome")
	}
}

func TestGenotypeString(t *testing.T) {
	genotype := NewBinaryGenotype(5)
	str := genotype.String()

	// Just check that it returns something non-empty
	if str == "" {
		t.Error("String() method returned empty string")
	}
}
