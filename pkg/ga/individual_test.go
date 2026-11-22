package ga

import (
	"math/rand"
	"testing"
)

func TestNewIntegerGenotypeDecodesAcrossRange(t *testing.T) {
	const minValue, maxValue = 10, 200

	genotype := NewIntegerGenotype(5000, minValue, maxValue, rand.New(rand.NewSource(1)))

	seen := make(map[int]bool)
	for i := range genotype.Genome {
		value, err := genotype.GetIntegerValue(i)
		if err != nil {
			t.Fatalf("GetIntegerValue(%d): %v", i, err)
		}
		if value < minValue || value > maxValue {
			t.Fatalf("value %d at %d outside [%d, %d]", value, i, minValue, maxValue)
		}
		seen[value] = true
	}

	if !seen[minValue] || !seen[maxValue] {
		t.Errorf("decoded values do not span [%d, %d]: min seen %v, max seen %v",
			minValue, maxValue, seen[minValue], seen[maxValue])
	}
}

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
	}
}

func TestFindBestIndividual(t *testing.T) {
	cases := []struct {
		population      []*Individual
		expectedFitness float64
	}{
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 3.0}},
				{Phenotype: &Phenotype{Fitness: 0.5}},
			},
			expectedFitness: 3.0,
		},
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 0.5}},
			},
			expectedFitness: 2.0,
		},
	}

	for _, tc := range cases {
		best := findBestIndividual(tc.population)

		if best.Phenotype.Fitness != tc.expectedFitness {
			t.Fatalf("Expected best fitness to be %f, but got %f", tc.expectedFitness, best.Phenotype.Fitness)
		}
	}
}
