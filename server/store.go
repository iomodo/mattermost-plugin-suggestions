package main

import (
	"encoding/json"
	"fmt"
)

const (
	timestampKey        = "timestamp"
	userChannelRanksKey = "userChannelRanks"
)

type recommendedChannel struct {
	ChannelID string  // identifier
	Score     float64 // score
}

// initStore method is for initializing the KVStore.
func (p *Plugin) initStore() error {
	err := p.saveTimestamp(0)
	if err != nil {
		return err
	}
	return p.saveUserChannelRanks(make(userChannelRank))
}

// saveUserRecommendations saves user recommendations in the KVStore.
func (p *Plugin) saveUserRecommendations(userID string, channels []*recommendedChannel) error {
	return p.save(userID, channels)
}

// retreiveUserRecomendations gets user recommendations from the KVStore.
func (p *Plugin) retreiveUserRecomendations(userID string) ([]*recommendedChannel, error) {
	recommendations := make([]*recommendedChannel, 0)
	err := p.retreive(userID, &recommendations)
	return recommendations, err
}

// saveTimestamp saves timestamp in the KVStore.
// All posts until this timestamp should already be analyzed.
func (p *Plugin) saveTimestamp(time int64) error {
	return p.save(timestampKey, time)
}

// retreiveTimestamp gets timestamp from KVStore.
func (p *Plugin) retreiveTimestamp() (int64, error) {
	var time int64
	err := p.retreive(timestampKey, &time)
	return time, err
}

// saveUserChannelRanks saves user-channel ranks in the KVStore.
func (p *Plugin) saveUserChannelRanks(ranks userChannelRank) error {
	return p.save(userChannelRanksKey, ranks)
}

// retreiveUserChannelRanks gets user-channel ranks from the KVStore.
func (p *Plugin) retreiveUserChannelRanks() (userChannelRank, error) {
	var ranks userChannelRank
	err := p.retreive(userChannelRanksKey, &ranks)
	return ranks, err
}

// save method saves generic value in the KVStore
func (p *Plugin) save(key string, value interface{}) (err error) {
	j, err := json.Marshal(value)
	println(fmt.Sprintf("%v", err))
	if err != nil {
		p.API.LogError("Can't marshal time", "err", err.Error())
		return err
	}
	appErr := p.API.KVSet(key, j)
	if appErr != nil {
		p.API.LogError("Can't set key", "err", appErr.Error())
		return appErr
	}
	return nil
}

// retreive method gets saved generic value from the KVStore
func (p *Plugin) retreive(key string, value interface{}) error {
	v, err := p.API.KVGet(key)
	if err != nil {
		p.API.LogError("can't get timestamp"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return json.Unmarshal(v, value)
}
