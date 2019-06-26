package main

import (
	"encoding/json"
)

const timestamp = "timestamp"

type recommendedChannel struct {
	ChannelID string  // identifier
	Score     float64 // score
}

func (p *Plugin) saveUserRecommendations(userID string, channels []recommendedChannel) error {
	j, err := json.Marshal(channels)

	if err != nil {
		p.API.LogError("Can't marshal recommendations", "err", err.Error())
		return err
	}
	if err := p.API.KVSet(userID, j); err != nil {
		p.API.LogError("Can't set recommendations"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return nil
}

func (p *Plugin) retreiveUserRecomendations(userID string) ([]*recommendedChannel, error) {
	j, err := p.API.KVGet(userID)
	if err != nil {
		p.API.LogError("can't get recommendations"+err.Error(), "err", err.Error()) //TODO
		return nil, err
	}
	recommendations := make([]*recommendedChannel, 0)
	if err := json.Unmarshal(j, &recommendations); err != nil {
		p.API.LogError("failed to unmarshal recommendations", "err", err.Error())
		return nil, err
	}
	return recommendations, nil
}

func (p *Plugin) saveTimestamp(time int64) error {
	t, err := json.Marshal(time)
	if err != nil {
		p.API.LogError("Can't marshal time", "err", err.Error())
		return err
	}
	if err := p.API.KVSet(timestamp, t); err != nil {
		p.API.LogError("Can't set timestamp"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return nil
}

func (p *Plugin) retreiveTimestamp() (int64, error) {
	t, err := p.API.KVGet(timestamp)
	if err != nil {
		p.API.LogError("can't get timestamp"+err.Error(), "err", err.Error()) //TODO
		return 0, err
	}
	var time int64
	if err := json.Unmarshal(t, &time); err != nil {
		p.API.LogError("failed to unmarshal timestamp", "err", err.Error())
		return 0, err
	}
	return time, nil
}
