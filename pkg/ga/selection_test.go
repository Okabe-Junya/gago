package ga

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func TestTournamentSelection(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
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
		selected := TournamentSelection(tc.population, tc.tournamentSize, rng)

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
	rng := rand.New(rand.NewSource(42))
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
		selected := RouletteWheelSelection(tc.population, rng)

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

// TestProportionalSelectionNoNilEntries is a regression test for a defect where
// proportional selectors assigned selected[i] only inside the inner selection
// loop with no fallback. Non-positive total fitness (common for minimization,
// where fitness is negative), NaN/Inf weights, or floating-point rounding could
// leave nil entries, which panic when later dereferenced during crossover.
func TestProportionalSelectionNoNilEntries(t *testing.T) {
	fitnessSets := map[string][]float64{
		"all-negative": {-1.0, -2.0, -3.0, -4.0},
		"all-zero":     {0.0, 0.0, 0.0, 0.0},
		"mixed-sign":   {-2.0, 1.0, -3.0, 4.0},
		"nan":          {math.NaN(), -1.0, 2.0, -3.0},
	}

	selectors := map[string]func([]*Individual, *rand.Rand) []*Individual{
		"RouletteWheelSelection":               RouletteWheelSelection,
		"RankSelection":                        RankSelection,
		"StochasticUniversalSamplingSelection": StochasticUniversalSamplingSelection,
		"BoltzmannSelection": func(pop []*Individual, rng *rand.Rand) []*Individual {
			return BoltzmannSelection(pop, 1.0, rng)
		},
		// A tiny temperature drives math.Exp to overflow to +Inf, exercising the
		// non-finite total-weight fallback path.
		"BoltzmannSelectionOverflow": func(pop []*Individual, rng *rand.Rand) []*Individual {
			return BoltzmannSelection(pop, 1e-6, rng)
		},
	}

	for setName, fitnesses := range fitnessSets {
		for selName, sel := range selectors {
			t.Run(setName+"/"+selName, func(t *testing.T) {
				rng := rand.New(rand.NewSource(1))
				population := make([]*Individual, len(fitnesses))
				for i, f := range fitnesses {
					population[i] = &Individual{
						Genotype:  &Genotype{Genome: []byte{byte(i)}},
						Phenotype: &Phenotype{Fitness: f},
					}
				}

				selected := sel(population, rng)
				if len(selected) != len(population) {
					t.Fatalf("expected %d selected, got %d", len(population), len(selected))
				}
				for i, ind := range selected {
					if ind == nil {
						t.Fatalf("selected[%d] is nil (would panic on later dereference)", i)
					}
				}
			})
		}
	}
}
