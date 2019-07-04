package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllUsers(t *testing.T) {
	assert := assert.New(t)
	t.Run("GetTeams error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		plugin.SetAPI(api)
		_, err := plugin.GetAllUsers()
		assert.NotNil(err)
	})

	t.Run("GetUsersInTeam error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		api.On("GetUsersInTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		plugin.SetAPI(api)
		_, err := plugin.GetAllUsers()
		assert.NotNil(err)
	})

	t.Run("No error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		users1 := make([]*model.User, 1)
		users1[0] = &model.User{Id: "userID1"}
		users2 := make([]*model.User, 1)
		users2[0] = &model.User{Id: "userID2"}
		correctUsers := map[string]*model.User{
			"userID1": &model.User{Id: "userID1"},
			"userID2": &model.User{Id: "userID2"},
		}

		api.On("GetUsersInTeam", mock.Anything, 0, mock.Anything).Return(users1, nil)
		api.On("GetUsersInTeam", mock.Anything, 1, mock.Anything).Return(users2, nil)
		api.On("GetUsersInTeam", mock.Anything, 2, mock.Anything).Return(make([]*model.User, 0), nil)
		plugin.SetAPI(api)
		users, err := plugin.GetAllUsers()
		assert.Nil(err)
		assert.Equal(correctUsers, users)
	})
}
