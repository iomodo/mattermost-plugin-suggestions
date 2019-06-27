package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankingUnion(t *testing.T) {
	assert := assert.New(t)
	t.Run("Add new ranks", func(t *testing.T) {
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
		res := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		assert.Equal(res, r1)
	})

	t.Run("Mix ranks", func(t *testing.T) {
		r1 := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
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
		res := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
				"channel1": 200,
				"channel2": 400,
			},
		}
		assert.Equal(res, r1)
	})
}

func TestGetRankingsSinceForChannel(t *testing.T) {

}
