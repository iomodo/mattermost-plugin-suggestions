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
	if err := p.API.KVSet(timestampKey, t); err != nil {
		p.API.LogError("Can't set timestamp"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return nil
}

func (p *Plugin) retreiveTimestamp() (int64, error) {
	t, err := p.API.KVGet(timestampKey)
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

func (p *Plugin) saveUserChannelRanks(ranks userChannelRank) error {
	r, err := json.Marshal(ranks)
	if err != nil {
		p.API.LogError("Can't marshal time", "err", err.Error())
		return err
	}
	if err := p.API.KVSet(userChannelRanksKey, r); err != nil {
		p.API.LogError("Can't set timestamp"+err.Error(), "err", err.Error()) //TODO
		return err
	}
	return nil
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

/*
func (p *Plugin) retreiveTimestamp() (int64, error) {
	t, err := p.API.KVGet(timestampKey)
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
*/
