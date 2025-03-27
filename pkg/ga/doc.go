// Package ga provides a comprehensive implementation of genetic algorithms
// for solving optimization and search problems.
//
// The genetic algorithm (GA) is a metaheuristic and population-based optimization
// algorithm inspired by natural selection. It uses operators like selection,
// crossover, and mutation to evolve a population of candidate solutions toward
// better solutions for the given problem.
//
// This package offers the following main components:
//
// Core GA Structure:
//   - GA struct: the main genetic algorithm implementation with methods for
//     initialization, evolution, and evaluation
//   - Population: manages the collection of individuals undergoing evolution
//   - Individual: represents a candidate solution with genotype and phenotype
//   - Genotype: the encoded representation of a solution
//   - Phenotype: the decoded representation with its fitness value
//   - Statistics: tracks metrics about the population during evolution
//
// Genetic Operators:
//   - Selection: methods for selecting parent individuals (tournament, roulette wheel, etc.)
//   - Crossover: methods for creating offspring from parents (single-point, multi-point, etc.)
//   - Mutation: methods for introducing variation (bit-flip, swap, etc.)
//
// Termination Conditions:
//   - Generation count: stops after a specified number of generations
//   - Time-based: stops after a specified duration
//   - Fitness threshold: stops when a target fitness is reached
//   - Convergence: stops when improvement stagnates
//   - Diversity: stops based on population diversity metrics
//
// Usage:
//
//	To use this package, create a GA instance, define genotype initialization and
//	phenotype evaluation functions, configure genetic operators and parameters,
//	then call Initialize() and Evolve() methods.
//
// Example:
//
//	ga := &ga.GA{
//	    Selection:     ga.TournamentSelection,
//	    Crossover:     ga.SinglePointCrossover,
//	    Mutation:      ga.BitFlipMutation,
//	    CrossoverRate: 0.7,
//	    MutationRate:  0.01,
//	    Generations:   100,
//	}
//
//	ga.Initialize(populationSize, initGenotype, evaluatePhenotype)
//	bestSolution, err := ga.Evolve(evaluatePhenotype)
package ga
