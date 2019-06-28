package ml

const defaultK = 10

// SimpleKNN struct
type SimpleKNN struct {
	params                  map[string]interface{}
	channelSimilarityMatrix [][]float64
	activityMatrix          [][]float64
	userIndexes             map[string]int64
	channelIndexes          map[string]int64
	similarity              funcSimilarity
	k                       int
}

// BaseEstimator determines interface for all estimators for user-channel suggestions
type BaseEstimator interface {
	SetParams(params map[string]interface{})
	Predict(userID, channelID string) float64
	Fit(rankings map[string]map[string]int64)
}

// NewSimpleKNN returns Simple KNN Estimator
func NewSimpleKNN(params map[string]interface{}) BaseEstimator {
	simpleKNN := new(SimpleKNN)
	simpleKNN.SetParams(params)
	return simpleKNN
}

// SetParams sets parameters for KNN estimator
func (knn *SimpleKNN) SetParams(params map[string]interface{}) {
	if val, exist := params["similarity"]; exist {
		switch val.(type) {
		case funcSimilarity:
			knn.similarity = val.(funcSimilarity)
		default:
			knn.similarity = cosineSimilarity
		}
	} else {
		knn.similarity = cosineSimilarity
	}
	if val, exist := params["k"]; exist {
		switch val.(type) {
		case int:
			knn.k = val.(int)
		default:
			knn.k = defaultK
		}
	} else {
		knn.k = defaultK
	}
}

func (knn *SimpleKNN) computeActivityMatrix(rankings map[string]map[string]int64) {
	knn.userIndexes = indexUsers(rankings)
	knn.channelIndexes = indexChannels(rankings)
	knn.activityMatrix = make([][]float64, len(knn.channelIndexes))
	for i := 0; i < len(knn.channelIndexes); i++ {
		knn.activityMatrix[i] = make([]float64, len(knn.userIndexes))
	}

	for user, m := range rankings {
		for channel, rank := range m {
			uIndex := knn.userIndexes[user]
			chIndex := knn.channelIndexes[channel]
			knn.activityMatrix[chIndex][uIndex] = float64(rank)
		}
	}

}

func (knn *SimpleKNN) computeSimilarityMatrix() {
	channelCount := len(knn.activityMatrix)
	knn.channelSimilarityMatrix = make([][]float64, channelCount)
	for i := 0; i < channelCount; i++ {
		knn.channelSimilarityMatrix[i] = make([]float64, channelCount)
	}
	for i := 0; i < channelCount; i++ {
		for j := 0; j < channelCount; j++ {
			knn.channelSimilarityMatrix[i][j] = knn.similarity(knn.activityMatrix[i], knn.activityMatrix[j])
		}
	}
}

// Fit the KNN estimator
func (knn *SimpleKNN) Fit(rankings map[string]map[string]int64) {
	knn.computeActivityMatrix(rankings)
	knn.computeSimilarityMatrix()
}

func (knn *SimpleKNN) getNeighbors(channel int64) []int64 {
	// chanVector := knn.channelSimilarityMatrix[channel]
	return nil
}

// Predict the rank of channel channelID for userID
func (knn *SimpleKNN) Predict(userID, channelID string) float64 {
	// channel, exists := knn.channelIndexes[channelID]
	// if !exists {
	// 	panic("Unknown channelID")
	// }
	// user, exists := knn.userIndexes[userID]
	// if !exists {
	// 	panic("Unknown userID")
	// }
	// neighbors := knn.getNeighbors(channel)

	return 0
}
