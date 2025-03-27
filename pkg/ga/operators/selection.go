// Package operators provides genetic algorithm operations such as selection, crossover, and mutation.
package operators

import (
	"math"
	"math/rand"
	"sort"

	"github.com/Okabe-Junya/gago/pkg/ga/population"
)

// TournamentSelection performs tournament selection on the given population.
//
// In tournament selection, a subset of individuals is randomly chosen from the population,
// and the individual with the highest fitness in this subset is selected. This process is repeated
// until the desired number of individuals is selected.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
// - tournamentSize: the number of individuals to be chosen randomly for each tournament.
//
// Returns:
// - A new population of selected individuals.
func TournamentSelection(individuals []*population.Individual, tournamentSize int) []*population.Individual {
	selected := make([]*population.Individual, len(individuals))
	for i := range selected {
		best := individuals[rand.Intn(len(individuals))]
		for j := 0; j < tournamentSize-1; j++ {
			contender := individuals[rand.Intn(len(individuals))]
			if contender.Phenotype.Fitness > best.Phenotype.Fitness {
				best = contender
			}
		}
		selected[i] = best
	}
	return selected
}

// RouletteWheelSelection performs roulette wheel selection on the given population.
//
// In roulette wheel selection, individuals are selected based on their fitness proportionate to
// the total fitness of the population. This method ensures that individuals with higher fitness
// have a higher chance of being selected.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A new population of selected individuals.
func RouletteWheelSelection(individuals []*population.Individual) []*population.Individual {
	totalFitness := 0.0
	for _, ind := range individuals {
		totalFitness += ind.Phenotype.Fitness
	}
	selected := make([]*population.Individual, len(individuals))
	for i := range selected {
		pick := rand.Float64() * totalFitness
		current := 0.0
		for _, ind := range individuals {
			current += ind.Phenotype.Fitness
			if current > pick {
				selected[i] = ind
				break
			}
		}
	}
	return selected
}

// RankSelection performs rank selection on the given population.
//
// In rank selection, individuals are first ranked according to their fitness,
// and then selection is performed based on these ranks rather than on the actual fitness values.
// This reduces selection pressure when the fitness variance is high.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A new population of selected individuals.
func RankSelection(individuals []*population.Individual) []*population.Individual {
	n := len(individuals)
	selected := make([]*population.Individual, n)
	// Clone the population to avoid modifying the original
	ranks := make([]*population.Individual, n)
	copy(ranks, individuals)
	// Sort by fitness in descending order
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Phenotype.Fitness > ranks[j].Phenotype.Fitness
	})
	// Calculate total rank sum: n*(n+1)/2
	totalRank := n * (n + 1) / 2
	// Perform selection based on ranks
	for i := range selected {
		pick := rand.Float64() * float64(totalRank)
		current := 0.0
		for j, ind := range ranks {
			// Ranks are 1-based, with highest fitness individual getting highest rank
			rank := n - j
			current += float64(rank)
			if current > pick {
				selected[i] = ind
				break
			}
		}
	}
	return selected
}

// StochasticUniversalSamplingSelection performs stochastic universal sampling on the given population.
//
// Stochastic Universal Sampling uses a single random value to sample all of the required individuals
// by choosing them at evenly spaced intervals. This gives a more diverse selection than roulette wheel selection.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A new population of selected individuals.
func StochasticUniversalSamplingSelection(individuals []*population.Individual) []*population.Individual {
	n := len(individuals)
	selected := make([]*population.Individual, n)
	totalFitness := 0.0
	for _, ind := range individuals {
		totalFitness += ind.Phenotype.Fitness
	}
	// Calculate the distance between the pointers
	distance := totalFitness / float64(n)
	// Choose a random starting point
	start := rand.Float64() * distance
	// Select individuals
	for i := range selected {
		pointer := start + float64(i)*distance
		current := 0.0
		for _, ind := range individuals {
			current += ind.Phenotype.Fitness
			if current > pointer {
				selected[i] = ind
				break
			}
		}
	}
	return selected
}

// TruncationSelection performs truncation selection on the given population.
//
// In truncation selection, the top portion of individuals sorted by fitness are selected,
// and the rest are discarded. The selected individuals are then duplicated to maintain the population size.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
// - truncationThreshold: the proportion of top individuals to select (between 0 and 1).
//
// Returns:
// - A new population of selected individuals.
func TruncationSelection(individuals []*population.Individual, truncationThreshold float64) []*population.Individual {
	n := len(individuals)
	selected := make([]*population.Individual, n)
	// Clone the population to avoid modifying the original
	sorted := make([]*population.Individual, n)
	copy(sorted, individuals)
	// Sort by fitness in descending order
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Phenotype.Fitness > sorted[j].Phenotype.Fitness
	})
	// Determine how many individuals to select
	selectCount := int(math.Ceil(float64(n) * truncationThreshold))
	if selectCount < 1 {
		selectCount = 1
	} else if selectCount > n {
		selectCount = n
	}
	// Fill the selected population
	for i := range selected {
		// Duplicate the top individuals as needed
		selected[i] = sorted[i%selectCount]
	}
	return selected
}

// BoltzmannSelection performs Boltzmann selection on the given population.
//
// Boltzmann selection is based on the principles of thermodynamics and uses a temperature parameter
// to control selection pressure. High temperatures lead to more uniform selection probabilities,
// while low temperatures increase selection pressure towards higher fitness individuals.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
// - temperature: the selection temperature (higher values mean more uniform selection).
//
// Returns:
// - A new population of selected individuals.
func BoltzmannSelection(individuals []*population.Individual, temperature float64) []*population.Individual {
	n := len(individuals)
	selected := make([]*population.Individual, n)
	// Calculate Boltzmann probabilities
	boltzmannValues := make([]float64, n)
	totalBoltzmann := 0.0
	for i, ind := range individuals {
		// Compute the Boltzmann probability
		boltzmannValues[i] = math.Exp(ind.Phenotype.Fitness / temperature)
		totalBoltzmann += boltzmannValues[i]
	}
	// Perform selection based on Boltzmann probabilities
	for i := range selected {
		pick := rand.Float64() * totalBoltzmann
		current := 0.0
		for j, ind := range individuals {
			current += boltzmannValues[j]
			if current > pick {
				selected[i] = ind
				break
			}
		}
	}
	return selected
}

// MultiObjectiveSelection performs selection for multi-objective optimization.
// It uses non-dominated sorting to assign fitness based on Pareto dominance and crowding distance.
//
// Parameters:
// - individuals: a slice of pointers to Individual, representing the current population.
// - objectives: a function that evaluates an individual and returns a slice of objective values.
//
// Returns:
// - A new population of selected individuals.
func MultiObjectiveSelection(
	individuals []*population.Individual,
	objectives func(*population.Individual) []float64,
) []*population.Individual {
	n := len(individuals)
	// Calculate the objective values for each individual
	objectiveValues := make([][]float64, n)
	for i, ind := range individuals {
		objectiveValues[i] = objectives(ind)
	}
	// Identify the Pareto fronts
	fronts := nonDominatedSort(individuals, objectiveValues)
	// Calculate crowding distance within each front
	for _, front := range fronts {
		calculateCrowdingDistance(front, objectiveValues)
	}
	// Create a new population by selecting from the fronts
	selected := make([]*population.Individual, n)
	selectedCount := 0
	// Add individuals from each front, starting with the best front
	for _, front := range fronts {
		// Sort the front by crowding distance (higher is better)
		sort.Slice(front, func(i, j int) bool {
			return front[i].Phenotype.Fitness > front[j].Phenotype.Fitness
		})
		// Add individuals from this front
		for _, ind := range front {
			if selectedCount >= n {
				break
			}
			selected[selectedCount] = ind
			selectedCount++
		}
		if selectedCount >= n {
			break
		}
	}
	return selected
}

// nonDominatedSort sorts individuals into Pareto fronts based on non-dominance.
// Returns a slice of slices, where each inner slice contains individuals from one front.
func nonDominatedSort(individuals []*population.Individual, objectiveValues [][]float64) [][]*population.Individual {
	n := len(individuals)
	fronts := [][]*population.Individual{}
	// Count how many solutions dominate each solution
	dominationCount := make([]int, n)
	// For each solution, store the solutions it dominates
	dominated := make([][]int, n)
	// Calculate domination relationships
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dominates(objectiveValues[i], objectiveValues[j]) {
				dominated[i] = append(dominated[i], j)
				dominationCount[j]++
			} else if dominates(objectiveValues[j], objectiveValues[i]) {
				dominated[j] = append(dominated[j], i)
				dominationCount[i]++
			}
		}
	}
	// Add the first front (non-dominated individuals)
	front := []*population.Individual{}
	for i := 0; i < n; i++ {
		if dominationCount[i] == 0 {
			front = append(front, individuals[i])
		}
	}
	fronts = append(fronts, front)
	// Create subsequent fronts
	currentFront := 0
	for len(fronts[currentFront]) > 0 {
		nextFront := []*population.Individual{}
		for _, ind := range fronts[currentFront] {
			i := indexOfIndividual(individuals, ind)
			for _, j := range dominated[i] {
				dominationCount[j]--
				if dominationCount[j] == 0 {
					nextFront = append(nextFront, individuals[j])
				}
			}
		}
		if len(nextFront) > 0 {
			fronts = append(fronts, nextFront)
			currentFront++
		} else {
			break
		}
	}
	return fronts
}

// dominates checks if solution a dominates solution b.
// a dominates b if a is not worse than b in all objectives and strictly better in at least one.
func dominates(a, b []float64) bool {
	better := false
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return false
		}
		if a[i] > b[i] {
			better = true
		}
	}
	return better
}

// calculateCrowdingDistance calculates the crowding distance for individuals in a front.
// The crowding distance is stored in each individual's Phenotype.Fitness field.
func calculateCrowdingDistance(front []*population.Individual, objectiveValues [][]float64) {
	n := len(front)
	if n <= 2 {
		// For the boundary points, set the crowding distance to a very large value
		for _, ind := range front {
			ind.Phenotype.Fitness = math.MaxFloat64
		}
		return
	}
	// Reset crowding distances
	for _, ind := range front {
		ind.Phenotype.Fitness = 0
	}
	numObjectives := len(objectiveValues[0])
	for m := 0; m < numObjectives; m++ {
		// Sort the front by the current objective
		sortByObjective(front, objectiveValues, m)
		// The boundary points have infinite distance
		front[0].Phenotype.Fitness = math.MaxFloat64
		front[n-1].Phenotype.Fitness = math.MaxFloat64
		// Calculate crowding distance for non-boundary points
		objectiveRange := getObjectiveRange(objectiveValues, m)
		if objectiveRange > 0 {
			for i := 1; i < n-1; i++ {
				idx1 := indexOfIndividual(front, front[i-1])
				idx2 := indexOfIndividual(front, front[i+1])
				// Add normalized distance to crowding distance
				front[i].Phenotype.Fitness += (objectiveValues[idx2][m] - objectiveValues[idx1][m]) / objectiveRange
			}
		}
	}
}

// sortByObjective sorts the front based on a specific objective value.
func sortByObjective(front []*population.Individual, objectiveValues [][]float64, m int) {
	sort.Slice(front, func(i, j int) bool {
		idxI := indexOfIndividual(front, front[i])
		idxJ := indexOfIndividual(front, front[j])
		return objectiveValues[idxI][m] < objectiveValues[idxJ][m]
	})
}

// getObjectiveRange calculates the range of values for a specific objective.
func getObjectiveRange(objectiveValues [][]float64, m int) float64 {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, values := range objectiveValues {
		if values[m] < min {
			min = values[m]
		}
		if values[m] > max {
			max = values[m]
		}
	}
	return max - min
}

// indexOfIndividual returns the index of an individual in a population.
func indexOfIndividual(individuals []*population.Individual, target *population.Individual) int {
	for i, ind := range individuals {
		if ind == target {
			return i
		}
	}
	return -1
}
