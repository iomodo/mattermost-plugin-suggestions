package main

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestExecuteCommandSuggestChannelsZero(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}

	channelsZero := make([]*recommendedChannel, 0)
	bytes, _ := json.Marshal(channelsZero)
	api.On("KVGet", mock.Anything).Return(bytes, (*model.AppError)(nil))
	plugin.SetAPI(api)

	args := &model.CommandArgs{
		Command: "/suggest channels",
	}
	resp, err := plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Equal(noNewChannelsText, resp.Text)
}

func TestExecuteCommandSuggestChannels(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}

	channels := make([]*recommendedChannel, 1)
	channels[0] = &recommendedChannel{ChannelID: "chan", Score: 0.1}
	bytes, _ := json.Marshal(channels)

	api.On("KVGet", mock.Anything).Return(bytes, (*model.AppError)(nil))
	api.On("GetChannel", mock.Anything).Return(&model.Channel{DisplayName: "CoolChannel"}, model.NewAppError("", "", nil, "", 404))
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
	plugin.SetAPI(api)

	args := &model.CommandArgs{
		Command: "/suggest channels",
	}
	_, err := plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
}

func TestExecuteCommandSuggestChannelError(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}

	channels := make([]*recommendedChannel, 1)
	channels[0] = &recommendedChannel{ChannelID: "chan", Score: 0.1}
	bytes, _ := json.Marshal(channels)

	api.On("KVGet", mock.Anything).Return(bytes, (*model.AppError)(nil))
	api.On("GetChannel", mock.Anything).Return(&model.Channel{Name: "CoolChannel"}, (*model.AppError)(nil))
	plugin.SetAPI(api)

	args := &model.CommandArgs{
		Command: "/suggest channels",
	}
	resp, err := plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Equal("Channels we recommend\n * ~CoolChannel - \n", resp.Text)
}

func TestExecuteCommandReset(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}
	api.On("KVSet", mock.Anything, mock.Anything).Return((*model.AppError)(nil))
	plugin.SetAPI(api)
	args := &model.CommandArgs{
		Command: "/suggest reset",
	}
	resp, err := plugin.ExecuteCommand(nil, args)
	assert.Nil(err)
	assert.Equal(resetText, resp.Text)

}
