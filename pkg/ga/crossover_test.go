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
		offspring := UniformCrossover(tc.population, tc.crossoverRate)

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
