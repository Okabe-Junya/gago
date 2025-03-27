package ga

import (
	"testing"
	"time"
)

func BenchmarkEvolution(b *testing.B) {
	// テストケースの設定
	testCases := []struct {
		name           string
		populationSize int
		genomeLength   int
		generations    int
		numWorkers     int
	}{
		{"Small", 100, 32, 10, 1},
		{"Small-Parallel", 100, 32, 10, 4},
		{"Medium", 1000, 64, 20, 1},
		{"Medium-Parallel", 1000, 64, 20, 8},
		{"Large", 10000, 128, 30, 1},
		{"Large-Parallel", 10000, 128, 30, 16},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			gaInstance := &GA{
				Selection:        func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
				Crossover:        SinglePointCrossover,
				Mutation:         BitFlipMutation,
				CrossoverRate:    0.7,
				MutationRate:     0.01,
				Generations:      tc.generations,
				NumParallelEvals: tc.numWorkers,
				EnableLogger:     false,
			}

			initFunc := func() *Genotype {
				return NewBinaryGenotype(tc.genomeLength)
			}

			evalFunc := func(genotype *Genotype) *Phenotype {
				// 計算コストをシミュレート
				time.Sleep(100 * time.Microsecond)
				fitness := 0.0
				for _, gene := range genotype.Genome {
					if gene == 1 {
						fitness += 1.0
					}
				}
				return &Phenotype{Fitness: fitness}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := gaInstance.Initialize(tc.populationSize, initFunc, evalFunc)
				if err != nil {
					b.Fatalf("Failed to initialize GA: %v", err)
				}
				_, err = gaInstance.Evolve(evalFunc)
				if err != nil {
					b.Fatalf("Failed to evolve population: %v", err)
				}
			}
		})
	}
}

func BenchmarkGeneticOperators(b *testing.B) {
	// テストケースの設定
	testCases := []struct {
		name           string
		populationSize int
		genomeLength   int
	}{
		{"Small", 100, 32},
		{"Medium", 1000, 64},
		{"Large", 10000, 128},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			// テスト用の個体を生成
			individuals := make([]*Individual, tc.populationSize)
			for i := range individuals {
				individuals[i] = &Individual{
					Genotype: NewBinaryGenotype(tc.genomeLength),
					Phenotype: &Phenotype{
						Fitness: float64(i),
					},
				}
			}

			// 選択演算子のベンチマーク
			b.Run("Selection", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					TournamentSelection(individuals, 3)
				}
			})

			// 交叉演算子のベンチマーク
			b.Run("Crossover", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					SinglePointCrossover(individuals, 0.7)
				}
			})

			// 突然変異演算子のベンチマーク
			b.Run("Mutation", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					BitFlipMutation(individuals, 0.01)
				}
			})
		})
	}
}

func BenchmarkMemoryUsage(b *testing.B) {
	// メモリ使用量を測定するためのテストケース
	testCases := []struct {
		name           string
		populationSize int
		genomeLength   int
		generations    int
	}{
		{"Small", 100, 32, 10},
		{"Medium", 1000, 64, 20},
		{"Large", 10000, 128, 30},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			gaInstance := &GA{
				Selection:     func(population []*Individual) []*Individual { return TournamentSelection(population, 3) },
				Crossover:     SinglePointCrossover,
				Mutation:      BitFlipMutation,
				CrossoverRate: 0.7,
				MutationRate:  0.01,
				Generations:   tc.generations,
				EnableLogger:  false,
			}

			initFunc := func() *Genotype {
				return NewBinaryGenotype(tc.genomeLength)
			}

			evalFunc := func(genotype *Genotype) *Phenotype {
				fitness := 0.0
				for _, gene := range genotype.Genome {
					if gene == 1 {
						fitness += 1.0
					}
				}
				return &Phenotype{Fitness: fitness}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := gaInstance.Initialize(tc.populationSize, initFunc, evalFunc)
				if err != nil {
					b.Fatalf("Failed to initialize GA: %v", err)
				}
				_, err = gaInstance.Evolve(evalFunc)
				if err != nil {
					b.Fatalf("Failed to evolve population: %v", err)
				}
			}
		})
	}
}
