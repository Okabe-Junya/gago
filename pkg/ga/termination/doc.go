// Package termination provides various termination conditions for genetic algorithms.
//
// This package includes several types of termination conditions, such as:
//
// - Generation count-based termination: stops evolution after a specific number of generations
//
// - Time-based termination: stops evolution after a specific duration
//
//   - Fitness-based termination: stops evolution when a fitness threshold is reached
//     or when improvement stagnates
//
// - Diversity-based termination: stops evolution based on population diversity metrics
//
// - Composite termination: combines multiple termination conditions with logical operators
//
// These termination conditions can be used with the genetic algorithm to control
// when the evolution process should stop, providing flexibility in balancing
// between computation time and solution quality.
package termination
