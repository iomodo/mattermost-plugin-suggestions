package main

import (
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type userChannelRank = map[string]map[string]int64 //map[userID]map[channelID][rank]

func (p *Plugin) getRanking() (userChannelRank, error) {
	previousTimestamp, err := p.retreiveTimestamp()
	if err != nil {
		return nil, err
	}
	timestampNow := time.Now().Unix()
	ranksSince, appErr := p.getRankingSince(previousTimestamp) // TODO what about the posts that where added between those lines?
	if appErr != nil {
		return nil, appErr
	}
	ranksUntil, err := p.retreiveUserChannelRanks()
	if err != nil {
		return nil, err
	}

	rankingUnion(ranksSince, ranksUntil)
	if err = p.saveTimestamp(timestampNow); err != nil {
		return nil, err
	}
	return ranksSince, nil
}

func (p *Plugin) getRankingSince(since int64) (userChannelRank, *model.AppError) {
	rankings := make(userChannelRank)
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
	rankings := make(userChannelRank)
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
		page++
	}
	return rankings, nil
}

func (p *Plugin) getRankingsSinceForChannel(channelID string, since int64) (userChannelRank, *model.AppError) {
	rankings := make(userChannelRank)
	postList, err := p.API.GetPostsSince(channelID, since)
	if err != nil {
		p.API.LogError("can't get posts since", "err", err.Error())
		return nil, err
	}
	posts := postList.ToSlice()
	for _, post := range posts {
		if _, ok := rankings[post.UserId]; !ok {
			rankings[post.UserId] = make(map[string]int64)
		}
		rankings[post.UserId][channelID]++
	}
	return rankings, nil
}

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
