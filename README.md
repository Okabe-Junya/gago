# GAGO: Simple Genetic Algorithm Library written in Go

[![Test](https://github.com/Okabe-Junya/gago/actions/workflows/test.yml/badge.svg)](https://github.com/Okabe-Junya/gago/actions/workflows/test.yml) [![CodeQL](https://github.com/Okabe-Junya/gago/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/Okabe-Junya/gago/actions/workflows/github-code-scanning/codeql) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/Okabe-Junya/gago)](https://goreportcard.com/report/github.com/Okabe-Junya/gago) [![Go Reference](https://pkg.go.dev/badge/github.com/Okabe-Junya/gago.svg)](https://pkg.go.dev/github.com/Okabe-Junya/gago)

## Overview

GAGO is a comprehensive Go library for implementing genetic algorithms. The library provides a flexible and extensible framework that supports various genetic algorithm operations and encodings.

### Key Features

- Multiple genome encodings:
  - Binary encoding
  - Integer encoding
  - Real-value encoding
  - Permutation encoding

- Selection methods:
  - Tournament selection
  - Roulette wheel selection
  - Rank selection
  - Stochastic universal sampling
  - Truncation selection
  - Boltzmann selection
  - Multi-objective selection (NSGA-II style)

- Crossover operations:
  - Single-point crossover
  - Multi-point crossover
  - Uniform crossover
  - Order-based crossover (for permutations)

- Mutation operations:
  - Bit-flip mutation
  - Swap mutation
  - Gaussian mutation
  - Inversion mutation
  - Scramble mutation
  - Uniform mutation
  - Adaptive mutation

- Termination conditions:
  - Generation count
  - Convergence threshold
  - Time-based
  - Fitness threshold

- Additional features:
  - Elitism
  - Parallel fitness evaluation
  - Adaptive parameter control
  - Comprehensive logging
  - Error handling
  - Type-safe genome operations

## Installation

```bash
go get -u github.com/Okabe-Junya/gago
```

## Quick Start

Here's a simple example of using GAGO to solve a maximization problem:

```go
func main() {
    // Configure the genetic algorithm
    gaInstance := &ga.GA{
        Selection:     func(population []*ga.Individual) []*ga.Individual {
            return ga.TournamentSelection(population, 3)
        },
        Crossover:     ga.SinglePointCrossover,
        Mutation:      ga.BitFlipMutation,
        CrossoverRate: 0.7,
        MutationRate:  0.01,
        Generations:   100,
        EnableLogger:  true,
        LogLevel:      logger.LevelInfo,
    }

    // Initialize the population
    gaInstance.Initialize(50, initializeGenotype, evaluatePhenotype)

    // Set termination condition
    gaInstance.TermCondition = ga.FitnessThresholdTermination(targetFitness)

    // Run the evolution
    bestIndividual := gaInstance.Evolve(evaluatePhenotype)

    // Get results
    fmt.Printf("Best fitness: %f\n", bestIndividual.Phenotype.Fitness)
    fmt.Printf("Total generations: %d\n", len(gaInstance.History)-1)
    fmt.Printf("Total runtime: %v\n", gaInstance.GetRuntime())
}
```

## Package Structure

- `pkg/ga`: Main genetic algorithm implementation
  - `encoding/`: Genome encoding types and operations
  - `operators/`: Selection, crossover, and mutation operators
  - `population/`: Population and individual management
- `internal/logger`: Logging functionality
- `examples/`: Example implementations
  - `find_max/`: Simple function maximization
  - `tsp/`: Traveling Salesman Problem solver
  - More examples to come

## Error Handling

The library provides comprehensive error handling:
- Input validation for all parameters
- Type-safe genome operations with bounds checking
- Panic recovery in parallel evaluations
- Logging of errors and warnings

## Logging

GAGO includes a flexible logging system that supports:
- Multiple log levels (Debug, Info, Warn, Error)
- JSON or text output format
- Custom output destinations
- Generation statistics
- Performance metrics
- Error tracking

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
