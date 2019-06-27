package main

import (
	"encoding/json"
)

const (
	timestampKey        = "timestamp"
	userChannelRanksKey = "userChannelRanks"
)

type recommendedChannel struct {
	ChannelID string  // identifier
	Score     float64 // score
}

func (p *Plugin) saveUserRecommendations(userID string, channels []*recommendedChannel) error {
	return p.save(userID, channels)
}

func (p *Plugin) retreiveUserRecomendations(userID string) ([]*recommendedChannel, error) {
	recommendations := make([]*recommendedChannel, 0)
	err := p.retreive(userID, &recommendations)
	return recommendations, err
}

func (p *Plugin) saveTimestamp(time int64) error {
	return p.save(timestampKey, time)
}

func (p *Plugin) retreiveTimestamp() (int64, error) {
	var time int64
	err := p.retreive(timestampKey, &time)
	return time, err
}

func (p *Plugin) saveUserChannelRanks(ranks userChannelRank) error {
	return p.save(userChannelRanksKey, ranks)
}

func (p *Plugin) retreiveUserChannelRanks() (userChannelRank, error) {
	var ranks userChannelRank
	err := p.retreive(userChannelRanksKey, &ranks)
	return ranks, err
}

func (p *Plugin) save(key string, value interface{}) error {
	j, err := json.Marshal(value)
	if err != nil {
		p.API.LogError("Can't marshal time", "err", err.Error())
		return err
	}
	return p.API.KVSet(key, j)
}

func (p *Plugin) retreive(key string, value interface{}) error {
	v, err := p.API.KVGet(key)
	if err != nil {
		p.API.LogError("can't get timestamp"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return json.Unmarshal(v, value)
}
