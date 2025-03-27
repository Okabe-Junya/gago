// Package ga provides adapters for encoding package types.
package ga

import (
	"github.com/Okabe-Junya/gago/pkg/ga/encoding"
)

// EncodingTypeAdapter provides adapters between ga package and encoding package.
// This file ensures backward compatibility with existing code.

// ConvertToEncodingGenotype converts a ga.Genotype to encoding.Genotype.
func ConvertToEncodingGenotype(g *Genotype) *encoding.Genotype {
	if g == nil {
		return nil
	}

	encodingGenotype := &encoding.Genotype{
		Genome:     make([]byte, len(g.Genome)),
		GenomeType: encoding.GenomeType(g.GenomeType),
	}

	copy(encodingGenotype.Genome, g.Genome)

	if len(g.MinValues) > 0 {
		encodingGenotype.MinValues = make([]float64, len(g.MinValues))
		copy(encodingGenotype.MinValues, g.MinValues)
	}

	if len(g.MaxValues) > 0 {
		encodingGenotype.MaxValues = make([]float64, len(g.MaxValues))
		copy(encodingGenotype.MaxValues, g.MaxValues)
	}

	return encodingGenotype
}

// ConvertFromEncodingGenotype converts an encoding.Genotype to ga.Genotype.
func ConvertFromEncodingGenotype(g *encoding.Genotype) *Genotype {
	if g == nil {
		return nil
	}

	gaGenotype := &Genotype{
		Genome:     make([]byte, len(g.Genome)),
		GenomeType: GenomeType(g.GenomeType),
	}

	copy(gaGenotype.Genome, g.Genome)

	if len(g.MinValues) > 0 {
		gaGenotype.MinValues = make([]float64, len(g.MinValues))
		copy(gaGenotype.MinValues, g.MinValues)
	}

	if len(g.MaxValues) > 0 {
		gaGenotype.MaxValues = make([]float64, len(g.MaxValues))
		copy(gaGenotype.MaxValues, g.MaxValues)
	}

	return gaGenotype
}
