package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankingUnion(t *testing.T) {
	assert := assert.New(t)

	r1 := userChannelRank{
		"user1": {
			"channel1": 100,
			"channel2": 200,
		},
	}
	r2 := userChannelRank{
		"user2": {
			"channel1": 100,
			"channel2": 200,
		},
	}
	rankingUnion(r1, r2)
	res1 := userChannelRank{
		"user1": {
			"channel1": 100,
			"channel2": 200,
		},
		"user2": {
			"channel1": 100,
			"channel2": 200,
		},
	}
	assert.Equal(res1, r1)

	rankingUnion(r1, r2)
	res2 := userChannelRank{
		"user1": {
			"channel1": 100,
			"channel2": 200,
		},
		"user2": {
			"channel1": 200,
			"channel2": 400,
		},
	}
	assert.Equal(res2, r1)
}
