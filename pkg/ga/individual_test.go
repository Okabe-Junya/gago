package ga

import "testing"

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
