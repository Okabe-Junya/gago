// Package ga provides functionalities for implementing genetic algorithms,
// including the main GA struct and its methods for initialization and evolution.
package ga

import (
	"fmt"

	"github.com/Okabe-Junya/gago/internal/logger"
)

// GA represents the genetic algorithm, including its population, genetic operators,
// and parameters for crossover and mutation rates, and the number of generations to evolve.
type GA struct {
	Population    []*Individual
	Selection     func([]*Individual) []*Individual
	Crossover     func([]*Individual, float64) []*Individual
	Mutation      func([]*Individual, float64)
	CrossoverRate float64
	MutationRate  float64
	Generations   int
	Logger        *logger.Logger
}

// Initialize initializes the population with the specified size, using the provided
// functions to create and evaluate genotypes.
//
// Parameters:
// - populationSize: the size of the population to be initialized.
// - initializeGenotype: a function to create a new Genotype.
// - evaluatePhenotype: a function to evaluate a Genotype and return its Phenotype.
func (ga *GA) Initialize(populationSize int, initializeGenotype func() *Genotype, evaluatePhenotype func(*Genotype) *Phenotype) {
	ga.Population = make([]*Individual, populationSize)
	for i := 0; i < populationSize; i++ {
		genotype := initializeGenotype()
		phenotype := evaluatePhenotype(genotype)
		ga.Population[i] = &Individual{Genotype: genotype, Phenotype: phenotype}
	}
}

// Evolve evolves the population over the specified number of generations, using the provided
// function to evaluate the fitness of each individual after applying selection, crossover,
// and mutation operations.
//
// Parameters:
// - evaluatePhenotype: a function to evaluate a Genotype and return its Phenotype.
func (ga *GA) Evolve(evaluatePhenotype func(*Genotype) *Phenotype) {
	for gen := 0; gen < ga.Generations; gen++ {
		ga.log(fmt.Sprintf("Generation %d", gen), "BestFitness", findBestIndividual(ga.Population).Phenotype.Fitness)
		ga.Population = ga.Selection(ga.Population)
		ga.Population = ga.Crossover(ga.Population, ga.CrossoverRate)
		ga.Mutation(ga.Population, ga.MutationRate)
		for _, ind := range ga.Population {
			ind.Phenotype = evaluatePhenotype(ind.Genotype)
		}
	}
}

// log logs a message with a key-value pair if the logger is set.
//
// Parameters:
// - msg: the message to log.
// - key: the key for the value being logged.
// - value: the value to log.
func (ga *GA) log(msg string, key string, value interface{}) {
	if ga.Logger != nil {
		ga.Logger.Log(msg, key, value)
	}
}
