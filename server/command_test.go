package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
)

// func assertThat(command, text string, )

func TestExecuteCommandTrivial(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}

	args := &model.CommandArgs{
		Command: "",
	}
	resp, err := plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Nil(resp)

	args = &model.CommandArgs{
		Command: "random",
	}
	resp, err = plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Equal(&model.CommandResponse{}, resp)

	args = &model.CommandArgs{
		Command: "/suggest",
	}
	resp, err = plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Contains(resp.Text, desc)

	args = &model.CommandArgs{
		Command: "/suggest help",
	}
	resp, err = plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Contains(resp.Text, desc)
}

// func TestExecuteCommandSuggestChannels(t *testing.T) {
// 	assert := assert.New(t)
// 	plugin := Plugin{}
// 	plugint := plugintest
// 	api := &plugintest.API{}
// 	api.On("RegisterCommand", mock.Anything).Return(nil)
// 	api.On("UnregisterCommand", mock.Anything, mock.Anything).Return(nil)
// 	api.On("GetUser", mock.Anything).Return(&model.User{
// 		Id:       "userid",
// 		Nickname: "User",
// 	}, (*model.AppError)(nil))

// 	args := &model.CommandArgs{
// 		Command: "",
// 	}
// 	resp, err := plugin.ExecuteCommand(nil, args)
// 	assert.Nil(err)
// 	assert.Nil(resp)
// }
