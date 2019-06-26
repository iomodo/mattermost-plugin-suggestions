package main

import "github.com/mattermost/mattermost-server/model"

func (p *Plugin) preCalculateRecommendations() {
	// since, err := p.retreiveTimestamp()
	// if err != nil {
	// 	return
	// }
	// postCounts, err := p.getUserPostCountsSince(since)
	// if err != nil {
	// 	return
	// }

}

func mapUnion(m1, m2 map[string]int64) {
	for k, v := range m2 {
		m1[k] += v
	}
}

func (p *Plugin) getUserPostCountsSince(since int64) (map[string]int64, *model.AppError) {
	postCounts := make(map[string]int64)
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError("can't get Teams", "err", err.Error())
		return nil, err
	}
	for _, team := range teams {
		postCountsForTeam, err := p.getUserPostCountsSinceForTeam(team.Id, since)
		if err != nil {
			return nil, err
		}
		mapUnion(postCounts, postCountsForTeam)
	}
	return postCounts, nil
}

func (p *Plugin) getUserPostCountsSinceForTeam(teamID string, since int64) (map[string]int64, *model.AppError) {
	postCounts := make(map[string]int64)
	page := 0
	perPage := 100
	for {
		channels, err := p.API.GetPublicChannelsForTeam(teamID, page, perPage)
		if err != nil {
			p.API.LogError("can't get public channels for a team", "err", err.Error())
			return nil, err
		}
		if len(channels) == 0 {
			break
		}
		for _, channel := range channels {
			postCountsForChannel, err := p.getUserPostCountSinceForChannel(channel.Id, since)
			if err != nil {
				return nil, err
			}
			mapUnion(postCounts, postCountsForChannel)
		}
	}
	return postCounts, nil
}

func (p *Plugin) getUserPostCountSinceForChannel(channelID string, since int64) (map[string]int64, *model.AppError) {
	postCounts := make(map[string]int64)
	postList, err := p.API.GetPostsSince(channelID, since)
	if err != nil {
		p.API.LogError("can't get posts since", "err", err.Error())
		return nil, err
	}
	posts := postList.ToSlice()
	for i := 0; i < len(posts); i++ {
		postCounts[posts[i].UserId]++
	}
	return postCounts, nil
}

func (p *Plugin) getUserPostCountsSince2(since int64) (map[string]int64, *model.AppError) {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError("can't get Teams", "err", err.Error())
		return nil, err
	}
	channels := make([]*model.Channel, 0)
	for i := 0; i < len(teams); i++ {
		page := 0
		perPage := 100
		for {
			c, err := p.API.GetPublicChannelsForTeam(teams[i].Id, page, perPage)
			if err != nil {
				p.API.LogError("can't get public channels for a team", "err", err.Error())
				return nil, err
			}
			if len(c) == 0 {
				break
			}
			channels = append(channels, c...)
			page++
		}
	}
	posts := make([]*model.Post, 0)
	for i := 0; i < len(channels); i++ {
		postList, err := p.API.GetPostsSince(channels[i].Id, since)
		if err != nil {
			p.API.LogError("can't get posts since", "err", err.Error())
			return nil, err
		}
		posts = append(posts, postList.ToSlice()...)
	}

	postCounts := make(map[string]int64)

	for i := 0; i < len(posts); i++ {
		postCounts[posts[i].UserId]++
	}
	return postCounts, nil
}
