// Package ga provides functionalities for implementing genetic algorithms.
package ga

import (
	"math"
	"math/rand"
	"sort"
)

// TournamentSelection implements tournament selection for selecting individuals.
// It randomly selects tournamentSize individuals and returns the best one.
func TournamentSelection(population []*Individual, tournamentSize int, rng *rand.Rand) []*Individual {
	if len(population) == 0 {
		return nil
	}

	selected := make([]*Individual, len(population))
	for i := range selected {
		// Select tournamentSize random individuals
		tournament := make([]*Individual, tournamentSize)
		for j := range tournament {
			tournament[j] = population[rng.Intn(len(population))]
		}

		// Find the best individual in the tournament
		best := tournament[0]
		for _, ind := range tournament[1:] {
			if ind.Phenotype.Fitness > best.Phenotype.Fitness {
				best = ind
			}
		}

		selected[i] = best
	}

	return selected
}

// isUsableTotalWeight reports whether a total weight is a finite, strictly
// positive number, i.e. usable as a denominator for proportional selection.
// Non-positive totals (e.g. all-negative fitness used for minimization) or
// non-finite totals (NaN/Inf) make proportional selection undefined.
func isUsableTotalWeight(total float64) bool {
	return total > 0 && !math.IsInf(total, 1)
}

// uniformSelection returns len(population) individuals chosen uniformly at
// random with replacement. Used as a fallback when proportional selection is
// undefined so the returned slice never contains nil entries.
func uniformSelection(population []*Individual, rng *rand.Rand) []*Individual {
	selected := make([]*Individual, len(population))
	for i := range selected {
		selected[i] = population[rng.Intn(len(population))]
	}
	return selected
}

// RouletteWheelSelection implements roulette wheel selection for selecting individuals.
// The probability of selection is proportional to the individual's fitness.
func RouletteWheelSelection(population []*Individual, rng *rand.Rand) []*Individual {
	if len(population) == 0 {
		return nil
	}

	// Calculate total fitness
	totalFitness := 0.0
	for _, ind := range population {
		totalFitness += ind.Phenotype.Fitness
	}

	// Proportional selection is only defined for a positive, finite total.
	// Fall back to uniform selection to avoid nil entries in the result.
	if !isUsableTotalWeight(totalFitness) {
		return uniformSelection(population, rng)
	}

	// Create cumulative fitness array
	cumulativeFitness := make([]float64, len(population))
	cumulativeFitness[0] = population[0].Phenotype.Fitness / totalFitness
	for i := 1; i < len(population); i++ {
		cumulativeFitness[i] = cumulativeFitness[i-1] + population[i].Phenotype.Fitness/totalFitness
	}

	// Select individuals using roulette wheel
	selected := make([]*Individual, len(population))
	for i := range selected {
		r := rng.Float64()
		// Default to the last bucket: floating-point rounding can leave the
		// final cumulative value just below r, which would otherwise leave
		// this slot nil.
		selected[i] = population[len(population)-1]
		for j, cumFitness := range cumulativeFitness {
			if r <= cumFitness {
				selected[i] = population[j]
				break
			}
		}
	}

	return selected
}

// RankSelection performs rank selection on the given population.
// The probability of selection is proportional to the individual's rank rather than its fitness.
func RankSelection(population []*Individual, rng *rand.Rand) []*Individual {
	if len(population) == 0 {
		return nil
	}

	// Create a copy of the population and sort by fitness
	sorted := make([]*Individual, len(population))
	copy(sorted, population)

	// Sort by fitness in descending order
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Phenotype.Fitness > sorted[j].Phenotype.Fitness
	})

	// Assign ranks
	ranks := make([]float64, len(sorted))
	totalRanks := 0.0

	for i := range ranks {
		// Rank starts from 1 (worst) to N (best)
		ranks[i] = float64(i + 1)
		totalRanks += ranks[i]
	}

	// Create selection probabilities based on rank
	probabilities := make([]float64, len(sorted))
	cumulativeProbability := 0.0

	for i := range probabilities {
		probabilities[i] = ranks[i] / totalRanks
		cumulativeProbability += probabilities[i]
		probabilities[i] = cumulativeProbability
	}

	// Select individuals based on rank probabilities
	selected := make([]*Individual, len(population))
	for i := range selected {
		r := rng.Float64()
		// Default to the last bucket to guard against floating-point rounding
		// leaving the final cumulative probability just below r.
		selected[i] = sorted[len(sorted)-1]
		for j, prob := range probabilities {
			if r <= prob {
				selected[i] = sorted[j]
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
// - population: a slice of pointers to Individual, representing the current population.
//
// Returns:
// - A new population of selected individuals.
func StochasticUniversalSamplingSelection(population []*Individual, rng *rand.Rand) []*Individual {
	n := len(population)
	if n == 0 {
		return nil
	}
	selected := make([]*Individual, n)
	totalFitness := 0.0

	for _, ind := range population {
		totalFitness += ind.Phenotype.Fitness
	}

	// SUS requires a positive, finite total; otherwise fall back to uniform
	// selection to avoid nil entries in the result.
	if !isUsableTotalWeight(totalFitness) {
		return uniformSelection(population, rng)
	}

	// Calculate the distance between the pointers
	distance := totalFitness / float64(n)

	// Choose a random starting point
	start := rng.Float64() * distance

	// Select individuals
	for i := range selected {
		pointer := start + float64(i)*distance
		current := 0.0
		// Default to the last individual so floating-point rounding on the
		// final pointer cannot leave this slot nil.
		selected[i] = population[n-1]
		for _, ind := range population {
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
// - population: a slice of pointers to Individual, representing the current population.
// - truncationThreshold: the proportion of top individuals to select (between 0 and 1).
//
// Returns:
// - A new population of selected individuals.
func TruncationSelection(population []*Individual, truncationThreshold float64) []*Individual {
	n := len(population)
	selected := make([]*Individual, n)

	// Clone the population to avoid modifying the original
	sorted := make([]*Individual, n)
	copy(sorted, population)

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
// - population: a slice of pointers to Individual, representing the current population.
// - temperature: the selection temperature (higher values mean more uniform selection).
//
// Returns:
// - A new population of selected individuals.
func BoltzmannSelection(population []*Individual, temperature float64, rng *rand.Rand) []*Individual {
	n := len(population)
	if n == 0 {
		return nil
	}
	selected := make([]*Individual, n)

	// Calculate Boltzmann probabilities
	boltzmannValues := make([]float64, n)
	totalBoltzmann := 0.0

	for i, ind := range population {
		// Compute the Boltzmann probability
		boltzmannValues[i] = math.Exp(ind.Phenotype.Fitness / temperature)
		totalBoltzmann += boltzmannValues[i]
	}

	// math.Exp can overflow to +Inf (or a non-positive temperature can yield
	// NaN); fall back to uniform selection when the total weight is unusable
	// to avoid nil entries in the result.
	if !isUsableTotalWeight(totalBoltzmann) {
		return uniformSelection(population, rng)
	}

	// Perform selection based on Boltzmann probabilities
	for i := range selected {
		pick := rng.Float64() * totalBoltzmann
		current := 0.0

		// Default to the last individual so floating-point rounding cannot
		// leave this slot nil.
		selected[i] = population[n-1]
		for j, ind := range population {
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
// - population: a slice of pointers to Individual, representing the current population.
// - objectives: a function that evaluates an individual and returns a slice of objective values.
//
// Returns:
// - A new population of selected individuals.
func MultiObjectiveSelection(
	population []*Individual,
	objectives func(*Individual) []float64,
) []*Individual {
	n := len(population)

	// Calculate the objective values for each individual
	objectiveValues := make([][]float64, n)
	for i, ind := range population {
		objectiveValues[i] = objectives(ind)
	}

	// Identify the Pareto fronts
	fronts := nonDominatedSort(population, objectiveValues)

	// Calculate crowding distance within each front
	for _, front := range fronts {
		calculateCrowdingDistance(front, objectiveValues)
	}

	// Create a new population by selecting from the fronts
	selected := make([]*Individual, n)
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
func nonDominatedSort(population []*Individual, objectiveValues [][]float64) [][]*Individual {
	n := len(population)
	fronts := [][]*Individual{}

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
	front := []*Individual{}
	for i := 0; i < n; i++ {
		if dominationCount[i] == 0 {
			front = append(front, population[i])
		}
	}
	fronts = append(fronts, front)

	// Create subsequent fronts
	currentFront := 0
	for len(fronts[currentFront]) > 0 {
		nextFront := []*Individual{}

		for _, ind := range fronts[currentFront] {
			i := indexOfIndividual(population, ind)

			for _, j := range dominated[i] {
				dominationCount[j]--
				if dominationCount[j] == 0 {
					nextFront = append(nextFront, population[j])
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
func calculateCrowdingDistance(front []*Individual, objectiveValues [][]float64) {
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
func sortByObjective(front []*Individual, objectiveValues [][]float64, m int) {
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
func indexOfIndividual(population []*Individual, target *Individual) int {
	for i, ind := range population {
		if ind == target {
			return i
		}
	}
	return -1
}
