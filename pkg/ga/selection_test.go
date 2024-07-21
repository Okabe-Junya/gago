package ga

import (
	"reflect"
	"testing"
)

func TestTournamentSelection(t *testing.T) {
	cases := []struct {
		population     []*Individual
		tournamentSize int
	}{
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 3.0}},
				{Phenotype: &Phenotype{Fitness: 4.0}},
			},
			tournamentSize: 2,
		},
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 1.5}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 2.5}},
			},
			tournamentSize: 3,
		},
	}

	for _, tc := range cases {
		selected := TournamentSelection(tc.population, tc.tournamentSize)

		if len(selected) != len(tc.population) {
			t.Fatalf("Expected selected length %d, but got %d", len(tc.population), len(selected))
		}

		for _, ind := range selected {
			found := false
			for _, original := range tc.population {
				if reflect.DeepEqual(ind, original) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Selected individual %+v not found in the original population", ind)
			}
		}
	}
}

func TestRouletteWheelSelection(t *testing.T) {
	cases := []struct {
		population []*Individual
	}{
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 3.0}},
				{Phenotype: &Phenotype{Fitness: 4.0}},
			},
		},
		{
			population: []*Individual{
				{Phenotype: &Phenotype{Fitness: 1.0}},
				{Phenotype: &Phenotype{Fitness: 1.5}},
				{Phenotype: &Phenotype{Fitness: 2.0}},
				{Phenotype: &Phenotype{Fitness: 2.5}},
			},
		},
	}

	for _, tc := range cases {
		selected := RouletteWheelSelection(tc.population)

		if len(selected) != len(tc.population) {
			t.Fatalf("Expected selected length %d, but got %d", len(tc.population), len(selected))
		}

		for _, ind := range selected {
			found := false
			for _, original := range tc.population {
				if reflect.DeepEqual(ind, original) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Selected individual %+v not found in the original population", ind)
			}
		}
	}
}
