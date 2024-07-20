package main

// func main() {
// 	logger := logger.NewLogger(true)

// 	gaInstance := &ga.GA{
// 		Selection:     func(population []*ga.Individual) []*ga.Individual { return ga.TournamentSelection(population, 3) },
// 		Crossover:     ga.SinglePointCrossover,
// 		Mutation:      ga.BitFlipMutation,
// 		CrossoverRate: 0.7,
// 		MutationRate:  0.01,
// 		Generations:   100,
// 		Logger:        logger,
// 	}

// 	gaInstance.Initialize(50, func() *ga.Genotype {
// 		genotype := ga.NewGenotype(10)
// 		for i := range genotype.Genome {
// 			genotype.Genome[i] = byte(rand.Intn(2))
// 		}
// 		return genotype
// 	}, func(genotype *ga.Genotype) *ga.Phenotype {
// 		fitness := 0.0
// 		for _, gene := range genotype.Genome {
// 			if gene == 1 {
// 				fitness++
// 			}
// 		}
// 		return &ga.Phenotype{Fitness: fitness}
// 	})

// 	gaInstance.Evolve(func(genotype *ga.Genotype) *ga.Phenotype {
// 		fitness := 0.0
// 		for _, gene := range genotype.Genome {
// 			if gene == 1 {
// 				fitness++
// 			}
// 		}
// 		return &ga.Phenotype{Fitness: fitness}
// 	})

// 	for _, ind := range gaInstance.Population {
// 		fmt.Printf("Genotype: %v, Fitness: %f\n", ind.Genotype.Genome, ind.Phenotype.Fitness)
// 	}
// }
