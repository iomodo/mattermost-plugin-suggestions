package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllUsers(t *testing.T) {
	t.Run("getTeamUsers error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		defer api.AssertExpectations(t)
		plugin.SetAPI(api)
		_, err := plugin.GetAllUsers()
		assert.NotNil(t, err)
	})

	t.Run("No error", func(t *testing.T) {
		correctUsers := map[string]*model.User{
			"userID1": &model.User{Id: "userID1"},
			"userID2": &model.User{Id: "userID2"},
		}
		plugin, api := getUsersInTeamPlugin()
		defer api.AssertExpectations(t)
		users, err := plugin.GetAllUsers()
		assert.Nil(t, err)
		assert.Equal(t, correctUsers, users)
	})
}

func TestGetAllChannels(t *testing.T) {
	t.Run("getTeamUsers error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		defer api.AssertExpectations(t)
		plugin.SetAPI(api)
		_, err := plugin.GetAllChannels()
		assert.NotNil(t, err)
	})

	t.Run("GetChannelsForTeamForUser error", func(t *testing.T) {
		plugin, api := getUsersInTeamPlugin()
		api.On("GetChannelsForTeamForUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		defer api.AssertExpectations(t)
		_, err := plugin.GetAllChannels()
		assert.NotNil(t, err)
	})

	t.Run("No error", func(t *testing.T) {
		plugin, api := getUsersInTeamPlugin()
		channels := []*model.Channel{
			&model.Channel{Id: "Id1"},
			&model.Channel{Id: "Id2"},
		}
		correctChannels := map[string]*model.Channel{
			"Id1": channels[0],
			"Id2": channels[1],
		}
		api.On("GetChannelsForTeamForUser", mock.Anything, mock.Anything, mock.Anything).Return(channels, nil)
		defer api.AssertExpectations(t)
		res, err := plugin.GetAllChannels()
		assert.Nil(t, err)
		assert.Equal(t, correctChannels, res)
	})
}

func TestGetTeamUsers(t *testing.T) {
	t.Run("GetTeams error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		defer api.AssertExpectations(t)
		plugin.SetAPI(api)
		_, err := plugin.getTeamUsers()
		assert.NotNil(t, err)
	})

	t.Run("GetUsersInTeam error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		api.On("GetUsersInTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		defer api.AssertExpectations(t)
		plugin.SetAPI(api)
		_, err := plugin.getTeamUsers()
		assert.NotNil(t, err)
	})

	t.Run("No error", func(t *testing.T) {
		correctUsers := map[string][]*model.User{
			"teamID": []*model.User{
				&model.User{Id: "userID1"},
				&model.User{Id: "userID2"},
			},
		}
		plugin, api := getUsersInTeamPlugin()
		defer api.AssertExpectations(t)
		users, err := plugin.getTeamUsers()
		assert.Nil(t, err)
		assert.Equal(t, correctUsers, users)
	})
}

func getUsersInTeamPlugin() (*Plugin, *plugintest.API) {
	api := &plugintest.API{}
	plugin := Plugin{}
	teams := make([]*model.Team, 1)
	teams[0] = &model.Team{Id: "teamID"}
	api.On("GetTeams").Return(teams, nil)
	users1 := make([]*model.User, 1)
	users1[0] = &model.User{Id: "userID1"}
	users2 := make([]*model.User, 1)
	users2[0] = &model.User{Id: "userID2"}

	api.On("GetUsersInTeam", mock.Anything, 0, mock.Anything).Return(users1, nil)
	api.On("GetUsersInTeam", mock.Anything, 1, mock.Anything).Return(users2, nil)
	api.On("GetUsersInTeam", mock.Anything, 2, mock.Anything).Return(make([]*model.User, 0), nil)
	plugin.SetAPI(api)
	return &plugin, api
}
