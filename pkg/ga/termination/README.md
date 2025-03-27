# Termination Conditions for Genetic Algorithms

This package provides various termination conditions for controlling when a genetic algorithm should stop its evolution process. These conditions can balance computation time against solution quality, allowing you to adapt the algorithm to your specific needs.

## Available Termination Conditions

### Basic Termination Conditions

- **Generation Count**: Terminates after a specified number of generations
  ```go
  // Stops after 100 generations
  ga.TermCondition = termination.GenerationCountTermination(100)
  ```

- **Time-Based**: Terminates after a specified duration
  ```go
  // Stops after 5 minutes
  ga.TermCondition = termination.TimeBasedTermination(5 * time.Minute)
  ```

### Fitness-Based Termination

- **Fitness Threshold**: Terminates when a fitness threshold is reached
  ```go
  // Stops when best fitness is at least 0.95
  ga.TermCondition = termination.FitnessThresholdTermination(0.95)
  ```

- **Fitness Stagnation**: Terminates when fitness hasn't improved for a number of generations
  ```go
  // Stops when there's no improvement for 20 generations
  ga.TermCondition = termination.FitnessStagnationTermination(20)
  ```

- **Fitness Improvement**: Terminates when the rate of fitness improvement falls below a threshold
  ```go
  // Stops when improvement rate is less than 0.001
  ga.TermCondition = termination.FitnessImprovementTermination(0.001)
  ```

### Diversity-Based Termination

- **Diversity Threshold**: Terminates when population diversity falls below a threshold
  ```go
  // Stops when diversity falls below 0.1
  ga.TermCondition = termination.DiversityThresholdTermination(0.1)
  ```

- **Diversity Stagnation**: Terminates when diversity hasn't improved for a number of generations
  ```go
  // Stops when diversity hasn't changed for 15 generations
  ga.TermCondition = termination.DiversityStagnationTermination(15)
  ```

- **Diversity Improvement**: Terminates when the rate of diversity improvement falls below a threshold
  ```go
  // Stops when diversity improvement rate is less than 0.005
  ga.TermCondition = termination.DiversityImprovementTermination(0.005)
  ```

### Composite Termination

You can combine multiple termination conditions using logical operators:

- **Any Termination (OR)**: Terminates when any of the conditions is met
  ```go
  // Stops after 100 generations OR when fitness reaches 0.95
  ga.TermCondition = termination.NewCompositeTermination(
      termination.AnyTermination,
      termination.GenerationCountTermination(100),
      termination.FitnessThresholdTermination(0.95),
  )
  ```

- **All Termination (AND)**: Terminates only when all conditions are met
  ```go
  // Stops when fitness is at least 0.9 AND diversity is below 0.2
  ga.TermCondition = termination.NewCompositeTermination(
      termination.AllTermination,
      termination.FitnessThresholdTermination(0.9),
      termination.DiversityThresholdTermination(0.2),
  )
  ```

## Example Usage

```go
package main

import (
    "time"

    "github.com/Okabe-Junya/gago/pkg/ga"
    "github.com/Okabe-Junya/gago/pkg/ga/termination"
)

func main() {
    // Create GA instance
    gaInstance := &ga.GA{
        Selection:     ga.TournamentSelection,
        Crossover:     ga.SinglePointCrossover,
        Mutation:      ga.BitFlipMutation,
        CrossoverRate: 0.7,
        MutationRate:  0.01,
        Generations:   500, // Maximum generations as a safety
    }

    // Set a composite termination condition
    gaInstance.TermCondition = termination.NewCompositeTermination(
        termination.AnyTermination,
        termination.GenerationCountTermination(500),           // Stop after 500 generations
        termination.FitnessThresholdTermination(0.95),         // or when fitness reaches 0.95
        termination.FitnessStagnationTermination(50),          // or when fitness stagnates for 50 generations
        termination.TimeBasedTermination(10 * time.Minute),    // or after 10 minutes
    )

    // Initialize and evolve...
    gaInstance.Initialize(populationSize, initGenotype, evaluatePhenotype)
    bestSolution, err := gaInstance.Evolve(evaluatePhenotype)

    // After evolution, we can find out which condition triggered termination
    // by examining the GA state
}
