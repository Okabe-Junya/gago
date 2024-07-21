# GAGO: Simple Genetic Algorithm Library written in Go

[![Test](https://github.com/Okabe-Junya/gago/actions/workflows/test.yml/badge.svg)](https://github.com/Okabe-Junya/gago/actions/workflows/test.yml) [![CodeQL](https://github.com/Okabe-Junya/gago/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/Okabe-Junya/gago/actions/workflows/github-code-scanning/codeql) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/Okabe-Junya/gago)](https://goreportcard.com/report/github.com/Okabe-Junya/gago) [![Go Reference](https://pkg.go.dev/badge/github.com/Okabe-Junya/gago.svg)](https://pkg.go.dev/github.com/Okabe-Junya/gago)

## Overview

**GAGO** is a go library for implementing genetic algorithms. The library is designed to be flexible and extensible, allowing users to define their own selection, crossover, and mutation functions.

## Installation

```bash
go get -u github.com/Okabe-Junya/gago
```

## Usage

See the [examples](./examples/) directory for more details. The following is a simple example of how to use the library.

> [!NOTE]
> You need to define some functions to use this example.

```go
func main() {
    gaInstance := &ga.GA{
        Selection:     func(population []*ga.Individual) []*ga.Individual { return ga.TournamentSelection(population, 3) },
        Crossover:     ga.SinglePointCrossover,
        Mutation:      ga.BitFlipMutation,
        CrossoverRate: crossoverRate,
        MutationRate:  mutationRate,
        Generations:   generations,
        EnableLogger:  true,
    }

    gaInstance.Initialize(populationSize, initializeGenotype, evaluatePhenotype)
    gaInstance.Evolve(evaluatePhenotype)

    bestIndividual := findBestIndividual(gaInstance.Population)
    bestX := decodeGenotype(bestIndividual.Genotype)

    fmt.Printf("Best x: %f, Fitness: %f\n", bestX, bestIndividual.Phenotype.Fitness)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
