package main

import "github.com/mattermost/mattermost-server/model"

type userChannelRank = map[string]map[string]int64 //map[userID]map[channelID][rank]

func rankingUnion(r1, r2 userChannelRank) {
	for userID2, m2 := range r2 {
		if _, ok := r1[userID2]; !ok {
			r1[userID2] = make(map[string]int64)
		}
		for channelID2, rank2 := range m2 {
			r1[userID2][channelID2] += rank2
		}
	}
}

func (p *Plugin) getRankingSince(since int64) (userChannelRank, *model.AppError) {
	var rankings userChannelRank
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError("can't get Teams", "err", err.Error())
		return nil, err
	}
	for _, team := range teams {
		rankingsForTeam, err := p.getRankingSinceForTeam(team.Id, since)
		if err != nil {
			return nil, err
		}
		rankingUnion(rankings, rankingsForTeam)
	}
	return rankings, nil
}

func (p *Plugin) getRankingSinceForTeam(teamID string, since int64) (userChannelRank, *model.AppError) {
	var rankings userChannelRank
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
			rankingsForChannel, err := p.getRankingsSinceForChannel(channel.Id, since)
			if err != nil {
				return nil, err
			}
			rankingUnion(rankings, rankingsForChannel)
		}
	}
	return rankings, nil
}

func (p *Plugin) getRankingsSinceForChannel(channelID string, since int64) (userChannelRank, *model.AppError) {
	var rankings userChannelRank
	postList, err := p.API.GetPostsSince(channelID, since)
	if err != nil {
		p.API.LogError("can't get posts since", "err", err.Error())
		return nil, err
	}
	posts := postList.ToSlice()
	for i := 0; i < len(posts); i++ {
		if _, ok := rankings[posts[i].UserId]; !ok {
			rankings[posts[i].UserId] = make(map[string]int64)
		}
		rankings[posts[i].UserId][channelID]++
	}
	return rankings, nil
}
