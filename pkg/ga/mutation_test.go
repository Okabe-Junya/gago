package ga

import (
	"reflect"
	"testing"
)

func TestBitFlipMutation(t *testing.T) {
	cases := []struct {
		population   []*Individual
		mutationRate float64
	}{
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1, 1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0, 0, 0, 0, 0}}},
			},
			mutationRate: 1.0, // All bits should be flipped
		},
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0}}},
			},
			mutationRate: 0.0, // No bits should be flipped
		},
	}

	for _, tc := range cases {
		original := make([]*Individual, len(tc.population))
		for i, ind := range tc.population {
			original[i] = &Individual{
				Genotype: &Genotype{
					Genome: append([]byte(nil), ind.Genotype.Genome...),
				},
				Phenotype: ind.Phenotype,
			}
		}

		BitFlipMutation(tc.population, tc.mutationRate)

		if tc.mutationRate == 1.0 {
			for i, ind := range tc.population {
				anyUnflippedBit := false
				for j, gene := range ind.Genotype.Genome {
					if gene == original[i].Genotype.Genome[j] {
						anyUnflippedBit = true
					}
				}
				if anyUnflippedBit {
					t.Errorf("Expected all bits to be flipped with mutation rate 1.0, but some bits remained unchanged in individual %d", i)
				}
			}
		} else if tc.mutationRate == 0.0 {
			for i, ind := range tc.population {
				if !reflect.DeepEqual(ind.Genotype.Genome, original[i].Genotype.Genome) {
					t.Errorf("Expected no mutation, but mutation occurred in individual %d", i)
				}
			}
		}
	}
}

func TestSwapMutation(t *testing.T) {
	cases := []struct {
		population   []*Individual
		mutationRate float64
	}{
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 2, 3, 4, 5, 6, 7, 8}}},
				{Genotype: &Genotype{Genome: []byte{8, 7, 6, 5, 4, 3, 2, 1}}},
			},
			mutationRate: 1.0, // Ensures swaps will happen
		},
		{
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 2, 3, 4}}},
				{Genotype: &Genotype{Genome: []byte{5, 6, 7, 8}}},
			},
			mutationRate: 0.0, // Ensures no swaps will happen
		},
	}

	for _, tc := range cases {
		original := make([]*Individual, len(tc.population))
		for i, ind := range tc.population {
			original[i] = &Individual{
				Genotype: &Genotype{
					Genome: append([]byte(nil), ind.Genotype.Genome...),
				},
				Phenotype: ind.Phenotype,
			}
		}

		// Since the same random seed would produce the same results,
		// we try multiple attempts to ensure the test is deterministic
		anyMutationOccurred := false
		for attempt := 0; attempt < 10; attempt++ {
			// Reset genomes to original state
			for i, ind := range tc.population {
				ind.Genotype.Genome = append([]byte(nil), original[i].Genotype.Genome...)
			}

			SwapMutation(tc.population, tc.mutationRate)

			if tc.mutationRate > 0.0 {
				// Check if mutation occurred in at least one individual
				for i, ind := range tc.population {
					if !reflect.DeepEqual(ind.Genotype.Genome, original[i].Genotype.Genome) {
						anyMutationOccurred = true
						break
					}
				}

				if anyMutationOccurred {
					break
				}
			}
		}

		if tc.mutationRate > 0.0 && !anyMutationOccurred {
			t.Errorf("Expected swap mutation to occur in at least one individual, but no mutation occurred")
		} else if tc.mutationRate == 0.0 {
			for i, ind := range tc.population {
				if !reflect.DeepEqual(ind.Genotype.Genome, original[i].Genotype.Genome) {
					t.Errorf("Expected no mutation, but mutation occurred in individual %d", i)
				}
			}
		}
	}
}
