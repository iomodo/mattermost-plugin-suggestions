package main

import (
	"encoding/json"
)

type recommendedChannel struct {
	channelID string  // identifier
	score     float64 // score
}

func (p *Plugin) saveUserRecommendations(userID string, channels []recommendedChannel) error {
	j, err := json.Marshal(channels)

	if err != nil {
		p.API.LogError("Can't marshal recommendations", "err", err.Error())
		return err
	}
	if err := p.API.KVSet(userID, j); err != nil {
		p.API.LogError("Can't set recommendations", "err", err.Error())
		return err
	}
	return nil
}

func (p *Plugin) retreiveUserRecomendations(userID string) ([]*recommendedChannel, error) {
	j, err := p.API.KVGet(userID)
	if err != nil {
		p.API.LogError("can't get recommendations", "err", err.Error())
		return nil, err
	}
	recommendations := make([]*recommendedChannel, 0)
	if err := json.Unmarshal(j, &recommendations); err != nil {
		p.API.LogError("failed to unmarshal recommendations", "err", err.Error())
		return nil, err
	}
	return recommendations, nil
}
