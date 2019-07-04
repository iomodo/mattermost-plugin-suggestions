package main

import "github.com/mattermost/mattermost-server/model"

// GetAllUsers returns all users
func (p *Plugin) GetAllUsers() (map[string]*model.User, *model.AppError) {
	allUsers := make(map[string]*model.User)
	teams, err := p.API.GetTeams()
	if err != nil {
		return nil, err
	}
	for _, team := range teams {
		page := 0
		perPage := 100
		for {
			users, err := p.API.GetUsersInTeam(team.Id, page, perPage)
			if err != nil {
				return nil, err
			}
			if len(users) == 0 {
				break
			}
			for _, user := range users {
				allUsers[user.Id] = user
			}
			page++
		}
	}
	return allUsers, nil
}
