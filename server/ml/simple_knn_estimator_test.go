package ml

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetParams(t *testing.T) {
	assert := assert.New(t)
	t.Run("wrong types in params", func(t *testing.T) {
		params := make(map[string]interface{})
		params["similarity"] = 1
		params["k"] = "bla"
		knn := new(SimpleKNN)
		knn.SetParams(params)
		assert.Equal(defaultK, knn.k)
		f1 := reflect.ValueOf(cosineSimilarity)
		f2 := reflect.ValueOf(knn.similarity)
		assert.Equal(f1.Pointer(), f2.Pointer())
	})
	t.Run("nil params", func(t *testing.T) {
		knn := new(SimpleKNN)
		knn.SetParams(nil)
		assert.Equal(defaultK, knn.k)
		f1 := reflect.ValueOf(cosineSimilarity)
		f2 := reflect.ValueOf(knn.similarity)
		assert.Equal(f1.Pointer(), f2.Pointer())
	})

	t.Run("good params", func(t *testing.T) {
		params := make(map[string]interface{})
		params["similarity"] = (funcSimilarity)(cosineSimilarity)
		params["k"] = 17
		knn := new(SimpleKNN)
		knn.SetParams(params)
		assert.Equal(17, knn.k)
		f1 := reflect.ValueOf(cosineSimilarity)
		f2 := reflect.ValueOf(knn.similarity)
		assert.Equal(f1.Pointer(), f2.Pointer())
	})

}

func getUserChannelRanks() map[string]map[string]int64 {
	m := make(map[string]map[string]int64)
	m["user1"] = make(map[string]int64)
	m["user1"]["chan1"] = 1
	m["user1"]["chan2"] = 1
	m["user1"]["chan3"] = 1
	m["user2"] = make(map[string]int64)
	m["user2"]["chan2"] = 1
	m["user2"]["chan3"] = 1
	m["user2"]["chan4"] = 1
	return m
}

func TestComputeActivityMatrix(t *testing.T) {
	assert := assert.New(t)
	m := getUserChannelRanks()
	knn := new(SimpleKNN)
	knn.SetParams(make(map[string]interface{}))
	knn.computeActivityMatrix(m)
	assert.Equal(4, len(knn.activityMatrix))
	for i := 0; i < 4; i++ {
		assert.Equal(2, len(knn.activityMatrix[i]))
	}
	for i := 0; i < 2; i++ {
		zeros := 0
		ones := 0
		for j := 0; j < 4; j++ {
			if knn.activityMatrix[j][i] == 0 {
				zeros++
			} else if knn.activityMatrix[j][i] == 1 {
				ones++
			} else {
				assert.Fail("in activity matrix should be only 0s and 1s")
			}
		}
		assert.Equal(1, zeros)
		assert.Equal(3, ones)
	}
}

func TestComputeSimilarityMatrix(t *testing.T) {
	assert := assert.New(t)
	m := getUserChannelRanks()
	knn := new(SimpleKNN)
	knn.SetParams(make(map[string]interface{}))
	knn.computeActivityMatrix(m)
	knn.computeSimilarityMatrix()
	assert.Equal(4, len(knn.channelSimilarityMatrix))
	for i := 0; i < 4; i++ {
		assert.Equal(4, len(knn.channelSimilarityMatrix[i]))
		assert.Equal(1.0, knn.channelSimilarityMatrix[i][i])
	}
	epsilon := 0.0001

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sim := knn.channelSimilarityMatrix[i][j]
			assert.True(sim == 0.0 || sim == 1.0 || sim-0.707106 < epsilon)
		}
	}
}

func TestFit(t *testing.T) {
	assert := assert.New(t)
	m := getUserChannelRanks()
	kn := NewSimpleKNN(nil)
	knn := kn.(*SimpleKNN)
	knn.Fit(m)

	assert.Equal(4, len(knn.activityMatrix))
	for i := 0; i < 4; i++ {
		assert.Equal(2, len(knn.activityMatrix[i]))
	}
	for i := 0; i < 2; i++ {
		zeros := 0
		ones := 0
		for j := 0; j < 4; j++ {
			if knn.activityMatrix[j][i] == 0 {
				zeros++
			} else if knn.activityMatrix[j][i] == 1 {
				ones++
			} else {
				assert.Fail("in activity matrix should be only 0s and 1s")
			}
		}
		assert.Equal(1, zeros)
		assert.Equal(3, ones)
	}

	assert.Equal(4, len(knn.channelSimilarityMatrix))
	for i := 0; i < 4; i++ {
		assert.Equal(4, len(knn.channelSimilarityMatrix[i]))
		assert.Equal(1.0, knn.channelSimilarityMatrix[i][i])
	}
	epsilon := 0.0001

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sim := knn.channelSimilarityMatrix[i][j]
			assert.True(sim == 0.0 || sim == 1.0 || sim-0.707106 < epsilon)
		}
	}
}
