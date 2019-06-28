package ml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	epsilon := 0.0001
	assert := assert.New(t)
	t.Run("same vectors", func(t *testing.T) {
		a := []float64{1, 2, 3}
		b := []float64{1, 2, 3}
		assert.Equal(1.0, cosineSimilarity(a, b))
	})

	t.Run("scaled vectors", func(t *testing.T) {
		a := []float64{1, 2, 3}
		b := []float64{2, 4, 6}
		assert.Equal(1.0, cosineSimilarity(a, b))
	})

	t.Run("random vectors", func(t *testing.T) {
		a := []float64{1, 2, 3}
		b := []float64{2, 3, 4}
		assert.InDelta(0.99258, cosineSimilarity(a, b), epsilon)
	})
}
