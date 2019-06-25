package main

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveUserRecommendationsNoError(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}
	api.On("KVSet", mock.Anything, mock.Anything).Return((*model.AppError)(nil))
	plugin.SetAPI(api)
	var channels []recommendedChannel
	err := plugin.saveUserRecommendations("randomUser", channels)
	assert.Nil(err)
}

func TestSaveUserRecommendationsWithError(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}
	api.On("KVSet", mock.Anything, mock.Anything).Return(model.NewAppError("", "", nil, "", 404))
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
	plugin.SetAPI(api)
	var channels []recommendedChannel
	err := plugin.saveUserRecommendations("randomUser", channels)
	assert.NotNil(err)
}

func TestRetreiveUserRecomendationsNoError(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}
	channels := make([]*recommendedChannel, 1)

	channels[0] = &recommendedChannel{ChannelID: "chan", Score: 0.1}
	bytes, _ := json.Marshal(channels)

	api.On("KVGet", mock.Anything).Return(bytes, (*model.AppError)(nil))
	plugin.SetAPI(api)
	c, err := plugin.retreiveUserRecomendations("randomUser")
	assert.Nil(err)
	assert.Equal(1, len(c))
	assert.Equal(channels[0], c[0])
}

func TestRetreiveUserRecomendationsWithError(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	api := &plugintest.API{}

	api.On("KVGet", mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
	plugin.SetAPI(api)
	_, err := plugin.retreiveUserRecomendations("randomUser")
	assert.NotNil(err)

}
