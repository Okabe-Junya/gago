package ga

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestSinglePointCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
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
		offspring := SinglePointCrossover(tc.population, tc.crossoverRate, rng)

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
	rng := rand.New(rand.NewSource(42))
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

			offspring := UniformCrossover(tc.population, tc.crossoverRate, rng)

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
			offspring := UniformCrossover(tc.population, tc.crossoverRate, rng)
			for i := 0; i < len(tc.population); i++ {
				if !reflect.DeepEqual(offspring[i], tc.population[i]) {
					t.Errorf("Expected no crossover to occur, but crossover happened for individual %d", i)
				}
			}
		}
	}
}

// permutationPair builds a population of two parents that are random
// permutations of [0, n).
func permutationPair(t *testing.T, n int) []*Individual {
	t.Helper()
	p1 := make([]byte, n)
	p2 := make([]byte, n)
	for i := 0; i < n; i++ {
		p1[i] = byte(i)
		p2[i] = byte(n - 1 - i)
	}
	return []*Individual{
		{Genotype: &Genotype{Genome: p1, GenomeType: PermutationEncoding}},
		{Genotype: &Genotype{Genome: p2, GenomeType: PermutationEncoding}},
	}
}

// assertIsPermutation verifies that the genome is a valid permutation of [0, n).
func assertIsPermutation(t *testing.T, label string, genome []byte) {
	t.Helper()
	n := len(genome)
	seen := make([]bool, n)
	for _, v := range genome {
		if int(v) >= n {
			t.Fatalf("%s: gene %d out of range for permutation of length %d", label, v, n)
		}
		if seen[v] {
			t.Fatalf("%s: duplicate gene %d in genome %v", label, v, genome)
		}
		seen[v] = true
	}
}

func TestTwoPointCrossover(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	cases := []struct {
		name           string
		population     []*Individual
		crossoverRate  float64
		expectedLength int
	}{
		{
			name: "crossover occurs",
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1, 1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0, 0, 0, 0, 0}}},
			},
			crossoverRate:  1.0,
			expectedLength: 2,
		},
		{
			name: "no crossover",
			population: []*Individual{
				{Genotype: &Genotype{Genome: []byte{1, 1, 1, 1}}},
				{Genotype: &Genotype{Genome: []byte{0, 0, 0, 0}}},
			},
			crossoverRate:  0.0,
			expectedLength: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			offspring := TwoPointCrossover(tc.population, tc.crossoverRate, rng)
			if len(offspring) != tc.expectedLength {
				t.Fatalf("Expected offspring length %d, got %d", tc.expectedLength, len(offspring))
			}
			for _, ind := range offspring {
				if len(ind.Genotype.Genome) != len(tc.population[0].Genotype.Genome) {
					t.Fatalf("Offspring genome length mismatch")
				}
			}
		})
	}
}

// TestOrderBasedCrossoverPreservesPermutation is the regression test for the
// previously broken OX1 implementation: every child must be a valid
// permutation of [0, n) when both parents are.
func TestOrderBasedCrossoverPreservesPermutation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		population := permutationPair(t, 20)
		offspring := OrderBasedCrossover(population, 1.0, rng)
		if len(offspring) != 2 {
			t.Fatalf("trial %d: expected 2 offspring, got %d", trial, len(offspring))
		}
		assertIsPermutation(t, "OX1 child1", offspring[0].Genotype.Genome)
		assertIsPermutation(t, "OX1 child2", offspring[1].Genotype.Genome)
	}
}

func TestPMXCrossoverPreservesPermutation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		population := permutationPair(t, 20)
		offspring := PMXCrossover(population, 1.0, rng)
		if len(offspring) != 2 {
			t.Fatalf("trial %d: expected 2 offspring, got %d", trial, len(offspring))
		}
		assertIsPermutation(t, "PMX child1", offspring[0].Genotype.Genome)
		assertIsPermutation(t, "PMX child2", offspring[1].Genotype.Genome)
	}
}

func TestCycleCrossoverPreservesPermutation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 50; trial++ {
		population := permutationPair(t, 20)
		offspring := CycleCrossover(population, 1.0, rng)
		if len(offspring) != 2 {
			t.Fatalf("trial %d: expected 2 offspring, got %d", trial, len(offspring))
		}
		assertIsPermutation(t, "CX child1", offspring[0].Genotype.Genome)
		assertIsPermutation(t, "CX child2", offspring[1].Genotype.Genome)
	}
}

// oddPopulation builds a population of the given odd size where every
// individual is a valid permutation of [0, genomeLen), so it is usable by both
// the binary/point crossovers and the permutation crossovers (OX1, PMX, CX).
func oddPopulation(size, genomeLen int, rng *rand.Rand) []*Individual {
	population := make([]*Individual, size)
	for i := range population {
		population[i] = &Individual{Genotype: NewPermutationGenotype(genomeLen, rng)}
	}
	return population
}

// TestCrossoverOddPopulationNoNilTail is the regression test for the odd-sized
// population bug: the pairwise crossover loop only filled indices
// [0, 2*(len/2)), leaving offspring[len-1] nil for odd populations, which
// panicked downstream (mutation) on a nil *Individual. Every returned slot must
// be a fully-formed individual.
func TestCrossoverOddPopulationNoNilTail(t *testing.T) {
	t.Parallel()

	const genomeLen = 6

	crossovers := []struct {
		name string
		fn   func([]*Individual, float64, *rand.Rand) []*Individual
	}{
		{"SinglePointCrossover", SinglePointCrossover},
		{"UniformCrossover", UniformCrossover},
		{"MultiPointCrossover", func(p []*Individual, r float64, rng *rand.Rand) []*Individual {
			return MultiPointCrossover(p, r, 2, rng)
		}},
		{"TwoPointCrossover", TwoPointCrossover},
		{"OrderBasedCrossover", OrderBasedCrossover},
		{"PMXCrossover", PMXCrossover},
		{"CycleCrossover", CycleCrossover},
	}

	sizes := []int{3, 5}

	for _, cx := range crossovers {
		cx := cx
		t.Run(cx.name, func(t *testing.T) {
			t.Parallel()
			for _, size := range sizes {
				size := size
				t.Run("size="+strconv.Itoa(size), func(t *testing.T) {
					t.Parallel()
					rng := rand.New(rand.NewSource(int64(size)))
					population := oddPopulation(size, genomeLen, rng)

					offspring := cx.fn(population, 1.0, rng)

					if len(offspring) != size {
						t.Fatalf("expected %d offspring, got %d", size, len(offspring))
					}
					for i, ind := range offspring {
						if ind == nil {
							t.Fatalf("offspring[%d] is nil", i)
						}
						if ind.Genotype == nil {
							t.Fatalf("offspring[%d].Genotype is nil", i)
						}
						if got := len(ind.Genotype.Genome); got != genomeLen {
							t.Fatalf("offspring[%d] genome length = %d, want %d", i, got, genomeLen)
						}
					}
				})
			}
		})
	}
}

// CX with identical parents must reproduce them exactly.
func TestCycleCrossoverIdentityWithIdenticalParents(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	uniq := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	population := []*Individual{
		{Genotype: &Genotype{Genome: append([]byte(nil), uniq...), GenomeType: PermutationEncoding}},
		{Genotype: &Genotype{Genome: append([]byte(nil), uniq...), GenomeType: PermutationEncoding}},
	}
	offspring := CycleCrossover(population, 1.0, rng)
	if !reflect.DeepEqual(offspring[0].Genotype.Genome, uniq) {
		t.Errorf("CX with identical parents: child1 = %v, want %v", offspring[0].Genotype.Genome, uniq)
	}
	if !reflect.DeepEqual(offspring[1].Genotype.Genome, uniq) {
		t.Errorf("CX with identical parents: child2 = %v, want %v", offspring[1].Genotype.Genome, uniq)
	}
}
