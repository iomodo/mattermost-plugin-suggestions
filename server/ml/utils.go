package ml

import "math"

type funcSimilarity func(a, b []float64) float64

// cosineSimilarity function
func cosineSimilarity(a, b []float64) float64 {
	ab := .0
	aa := .0
	bb := .0
	for i := 0; i < len(a); i++ {
		ab += a[i] * b[i]
		aa += a[i] * a[i]
		bb += b[i] * b[i]
	}
	return ab / math.Sqrt(aa*bb)
}
