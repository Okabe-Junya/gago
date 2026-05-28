// Package ga provides functionalities for implementing genetic algorithms,
// including crossover operations for generating offspring from parent individuals.
package ga

import (
	"math/rand"
	"sort"
)

// SinglePointCrossover performs a single-point crossover on the given population.
//
// In single-point crossover, a random crossover point is selected, and the
// offspring are created by exchanging the segments of the parent individuals' genomes
// after this point.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func SinglePointCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			point := rng.Intn(len(parent1.Genome))
			child1 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			child2 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			copy(child1.Genome[:point], parent1.Genome[:point])
			copy(child1.Genome[point:], parent2.Genome[point:])
			copy(child2.Genome[:point], parent2.Genome[:point])
			copy(child2.Genome[point:], parent1.Genome[point:])
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// UniformCrossover performs a uniform crossover on the given population.
//
// In uniform crossover, each gene from the parent individuals is independently
// chosen with a 50% probability to be included in the offspring. This allows
// for more genetic diversity in the offspring compared to single-point crossover.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func UniformCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			child1 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			child2 := &Genotype{Genome: make([]byte, len(parent1.Genome))}
			for j := range parent1.Genome {
				if rng.Float64() < 0.5 {
					child1.Genome[j] = parent1.Genome[j]
					child2.Genome[j] = parent2.Genome[j]
				} else {
					child1.Genome[j] = parent2.Genome[j]
					child2.Genome[j] = parent1.Genome[j]
				}
			}
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// MultiPointCrossover performs a multi-point crossover on the given population.
//
// In multi-point crossover, multiple crossover points are selected, and the
// offspring are created by exchanging segments of the parent individuals' genomes
// between these points. This allows for more genetic material exchange compared to
// single-point crossover.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
// - numPoints: the number of crossover points to use.
//
// Returns:
// - A new population of offspring generated from the input population.
func MultiPointCrossover(population []*Individual, crossoverRate float64, numPoints int, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype

			// Generate crossover points
			genomeLength := len(parent1.Genome)
			if numPoints > genomeLength-1 {
				numPoints = genomeLength - 1
			}

			points := make([]int, numPoints)
			for j := 0; j < numPoints; j++ {
				points[j] = rng.Intn(genomeLength)
			}
			sort.Ints(points)

			// Create children
			child1 := &Genotype{Genome: make([]byte, genomeLength)}
			child2 := &Genotype{Genome: make([]byte, genomeLength)}

			// Start with parent1's genes for child1 and parent2's genes for child2
			swap := false
			start := 0

			for j := 0; j < numPoints; j++ {
				end := points[j]

				if !swap {
					copy(child1.Genome[start:end], parent1.Genome[start:end])
					copy(child2.Genome[start:end], parent2.Genome[start:end])
				} else {
					copy(child1.Genome[start:end], parent2.Genome[start:end])
					copy(child2.Genome[start:end], parent1.Genome[start:end])
				}

				swap = !swap
				start = end
			}

			// Handle the last segment
			if !swap {
				copy(child1.Genome[start:], parent1.Genome[start:])
				copy(child2.Genome[start:], parent2.Genome[start:])
			} else {
				copy(child1.Genome[start:], parent2.Genome[start:])
				copy(child2.Genome[start:], parent1.Genome[start:])
			}

			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// TwoPointCrossover performs a two-point crossover on the given population.
//
// Two cut points are selected uniformly at random and the genes between them
// are swapped between the parents to form the two children.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func TwoPointCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))
	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			genomeLength := len(parent1.Genome)
			if genomeLength < 3 {
				// Fall back to single-point semantics for very short genomes.
				offspring[2*i] = population[2*i]
				offspring[2*i+1] = population[2*i+1]
				continue
			}
			a := rng.Intn(genomeLength - 1)
			b := a + 1 + rng.Intn(genomeLength-1-a)

			child1 := &Genotype{Genome: make([]byte, genomeLength), GenomeType: parent1.GenomeType}
			child2 := &Genotype{Genome: make([]byte, genomeLength), GenomeType: parent2.GenomeType}
			copy(child1.Genome[:a], parent1.Genome[:a])
			copy(child1.Genome[a:b], parent2.Genome[a:b])
			copy(child1.Genome[b:], parent1.Genome[b:])
			copy(child2.Genome[:a], parent2.Genome[:a])
			copy(child2.Genome[a:b], parent1.Genome[a:b])
			copy(child2.Genome[b:], parent2.Genome[b:])
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// OrderBasedCrossover performs Davis Order Crossover (OX1) on the given population.
// Suitable for permutation encodings (e.g., TSP): both children remain valid
// permutations of the parents' shared gene multiset.
//
// A random segment [a, b) of parent1 is copied to child1 at the same positions.
// The remaining positions are filled with parent2's genes in their original
// order, starting at index b (with wrap-around) and skipping genes already
// present in the copied segment. child2 is produced symmetrically.
//
// Parameters:
// - population: a slice of pointers to Individual, representing the current population.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func OrderBasedCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))

	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			genomeLength := len(parent1.Genome)
			if genomeLength < 2 {
				offspring[2*i] = population[2*i]
				offspring[2*i+1] = population[2*i+1]
				continue
			}

			a := rng.Intn(genomeLength)
			b := rng.Intn(genomeLength)
			if a > b {
				a, b = b, a
			}
			// Use the half-open interval [a, b+1) so the segment is non-empty.
			end := b + 1

			child1 := &Genotype{Genome: orderCrossoverChild(parent1.Genome, parent2.Genome, a, end), GenomeType: parent1.GenomeType}
			child2 := &Genotype{Genome: orderCrossoverChild(parent2.Genome, parent1.Genome, a, end), GenomeType: parent2.GenomeType}
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}

	return offspring
}

// orderCrossoverChild builds one OX1 child: copy parent1[start:end] into the
// child at the same positions, then fill the remaining slots with parent2's
// genes in their original order starting from index end (with wrap-around),
// skipping genes already in the copied segment.
func orderCrossoverChild(p1, p2 []byte, start, end int) []byte {
	n := len(p1)
	child := make([]byte, n)
	used := make(map[byte]bool, end-start)
	for i := start; i < end; i++ {
		child[i] = p1[i]
		used[p1[i]] = true
	}
	write := end % n
	for offset := 0; offset < n; offset++ {
		gene := p2[(end+offset)%n]
		if used[gene] {
			continue
		}
		// Skip the segment positions.
		if write == start {
			write = end % n
		}
		child[write] = gene
		used[gene] = true
		write = (write + 1) % n
		if write == start {
			write = end % n
		}
	}
	return child
}

// PMXCrossover performs Partially-Mapped Crossover (PMX) on the given population.
// Suitable for permutation encodings.
//
// A random segment [a, b] of parent1 is copied to child1 at the same positions.
// The remaining positions are inherited from parent2; any gene already used by
// the copied segment is resolved by walking a mapping table derived from the
// two segments until an unused gene is found. child2 is produced symmetrically.
//
// Parameters:
// - population: a slice of pointers to Individual.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func PMXCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))

	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			n := len(parent1.Genome)
			if n < 2 {
				offspring[2*i] = population[2*i]
				offspring[2*i+1] = population[2*i+1]
				continue
			}

			a := rng.Intn(n)
			b := rng.Intn(n)
			if a > b {
				a, b = b, a
			}

			child1 := &Genotype{Genome: pmxChild(parent1.Genome, parent2.Genome, a, b), GenomeType: parent1.GenomeType}
			child2 := &Genotype{Genome: pmxChild(parent2.Genome, parent1.Genome, a, b), GenomeType: parent2.GenomeType}
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

// pmxChild builds one PMX child by copying p1[a:b+1] into the child and
// resolving conflicts in positions outside the segment via the mapping
// p1[i] -> p2[i] for i in [a, b].
func pmxChild(p1, p2 []byte, a, b int) []byte {
	n := len(p1)
	child := make([]byte, n)
	copy(child, p2)
	inSegment := make(map[byte]bool, b-a+1)
	mapping := make(map[byte]byte, b-a+1)
	for i := a; i <= b; i++ {
		child[i] = p1[i]
		inSegment[p1[i]] = true
		mapping[p1[i]] = p2[i]
	}
	for i := 0; i < n; i++ {
		if i >= a && i <= b {
			continue
		}
		val := p2[i]
		seen := make(map[byte]bool)
		for inSegment[val] {
			if seen[val] {
				// Defensive: malformed input (parents not the same multiset).
				break
			}
			seen[val] = true
			val = mapping[val]
		}
		child[i] = val
	}
	return child
}

// CycleCrossover performs Cycle Crossover (CX) on the given population.
// Suitable for permutation encodings.
//
// CX identifies the positional cycles between the two parents. Even-indexed
// cycles take their values from parent1 for child1 (and from parent2 for
// child2); odd-indexed cycles swap. Every position in each child comes from
// one of the parents at the same index, so absolute positions are preserved.
//
// Parameters:
// - population: a slice of pointers to Individual.
// - crossoverRate: the probability with which crossover will occur.
//
// Returns:
// - A new population of offspring generated from the input population.
func CycleCrossover(population []*Individual, crossoverRate float64, rng *rand.Rand) []*Individual {
	offspring := make([]*Individual, len(population))

	for i := 0; i < len(population)/2; i++ {
		if rng.Float64() < crossoverRate {
			parent1 := population[2*i].Genotype
			parent2 := population[2*i+1].Genotype
			n := len(parent1.Genome)
			if n < 2 {
				offspring[2*i] = population[2*i]
				offspring[2*i+1] = population[2*i+1]
				continue
			}

			c1Genome, c2Genome := cycleCrossoverChildren(parent1.Genome, parent2.Genome)
			child1 := &Genotype{Genome: c1Genome, GenomeType: parent1.GenomeType}
			child2 := &Genotype{Genome: c2Genome, GenomeType: parent2.GenomeType}
			offspring[2*i] = &Individual{Genotype: child1}
			offspring[2*i+1] = &Individual{Genotype: child2}
		} else {
			offspring[2*i] = population[2*i]
			offspring[2*i+1] = population[2*i+1]
		}
	}
	return offspring
}

func cycleCrossoverChildren(p1, p2 []byte) (child1, child2 []byte) {
	n := len(p1)
	child1 = make([]byte, n)
	child2 = make([]byte, n)
	visited := make([]bool, n)
	indexInP2 := make(map[byte]int, n)
	for i, v := range p2 {
		indexInP2[v] = i
	}
	cycle := 0
	for start := 0; start < n; start++ {
		if visited[start] {
			continue
		}
		i := start
		for !visited[i] {
			visited[i] = true
			if cycle%2 == 0 {
				child1[i] = p1[i]
				child2[i] = p2[i]
			} else {
				child1[i] = p2[i]
				child2[i] = p1[i]
			}
			next, ok := indexInP2[p1[i]]
			if !ok {
				// Defensive: malformed input. Break the cycle.
				break
			}
			i = next
		}
		cycle++
	}
	return child1, child2
}
