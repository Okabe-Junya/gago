package ga

import (
	"reflect"
	"testing"
)

func TestSinglePointCrossover(t *testing.T) {
	cases := []struct {
		population     []*Individual
		crossoverRate  float64
		expectedLength int
	}{
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0}}},
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0}}},
			},
			crossoverRate:  1.0,
			expectedLength: 4,
		},
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0}}},
			},
			crossoverRate:  0.0,
			expectedLength: 2,
		},
	}

	for _, tc := range cases {
		offspring := SinglePointCrossover(tc.population, tc.crossoverRate)

		if len(offspring) != tc.expectedLength {
			t.Fatalf("Expected offspring length %d, but got %d", tc.expectedLength, len(offspring))
		}

		for i := 0; i < len(tc.population)/2; i++ {
			if tc.crossoverRate == 1.0 && reflect.DeepEqual(offspring[2*i], tc.population[2*i]) && reflect.DeepEqual(offspring[2*i+1], tc.population[2*i+1]) {
				t.Errorf("Expected crossover to occur, but no crossover happened for pair %d", i)
			} else if tc.crossoverRate == 0.0 && (!reflect.DeepEqual(offspring[2*i], tc.population[2*i]) || !reflect.DeepEqual(offspring[2*i+1], tc.population[2*i+1])) {
				t.Errorf("Expected no crossover to occur, but crossover happened for pair %d", i)
			}
		}
	}
}

func TestUniformCrossover(t *testing.T) {
	cases := []struct {
		population     []*Individual
		crossoverRate  float64
		expectedLength int
	}{
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1, 1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0, 0, 0, 0, 0}}},
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1, 1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0, 0, 0, 0, 0}}},
			},
			crossoverRate:  1.0,
			expectedLength: 4,
		},
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0}}},
			},
			crossoverRate:  0.0,
			expectedLength: 2,
		},
	}

	for _, tc := range cases {
		// Store original individuals
		original := make([]*Individual, len(tc.population))
		for i, ind := range tc.population {
			original[i] = &Individual{
				Genotype: &Genotype{
					Genome: append([]byte(nil), ind.Genotype.Genome...),
				},
				Phenotype: ind.Phenotype,
			}
		}

		// Try multiple attempts to make the test deterministic
		anyCrossoverOccurred := false
		for attempt := 0; attempt < 10; attempt++ {
			// Reset test case individuals
			for i, ind := range tc.population {
				ind.Genotype.Genome = append([]byte(nil), original[i].Genotype.Genome...)
			}

			offspring := UniformCrossover(tc.population, tc.crossoverRate)

			if len(offspring) != tc.expectedLength {
				t.Fatalf("Expected offspring length %d, but got %d", tc.expectedLength, len(offspring))
			}

			if tc.crossoverRate > 0.0 {
				// Check if crossover occurred in at least one pair
				for i := 0; i < len(tc.population)/2; i++ {
					if !reflect.DeepEqual(offspring[2*i], tc.population[2*i]) ||
						!reflect.DeepEqual(offspring[2*i+1], tc.population[2*i+1]) {
						anyCrossoverOccurred = true
						break
					}
				}

				if anyCrossoverOccurred {
					break
				}
			}
		}

		if tc.crossoverRate > 0.0 && !anyCrossoverOccurred {
			t.Errorf("Expected crossover to occur in at least one pair, but no crossover happened")
		} else if tc.crossoverRate == 0.0 {
			// Crossover rate of 0 should yield identical offspring
			offspring := UniformCrossover(tc.population, tc.crossoverRate)
			for i := 0; i < len(tc.population); i++ {
				if !reflect.DeepEqual(offspring[i], tc.population[i]) {
					t.Errorf("Expected no crossover to occur, but crossover happened for individual %d", i)
				}
			}
		}
	}
}
