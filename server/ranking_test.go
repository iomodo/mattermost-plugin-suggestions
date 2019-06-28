package main

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRankingUnion(t *testing.T) {
	assert := assert.New(t)
	t.Run("Add new ranks", func(t *testing.T) {
		r1 := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		r2 := userChannelRank{
			"user2": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		rankingUnion(r1, r2)
		res := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		assert.Equal(res, r1)
	})

	t.Run("Mix ranks", func(t *testing.T) {
		r1 := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		r2 := userChannelRank{
			"user2": {
				"channel1": 100,
				"channel2": 200,
			},
		}
		rankingUnion(r1, r2)
		res := userChannelRank{
			"user1": {
				"channel1": 100,
				"channel2": 200,
			},
			"user2": {
				"channel1": 200,
				"channel2": 400,
			},
		}
		assert.Equal(res, r1)
	})
}

func createMockPostList() (*model.PostList, userChannelRank) {
	postList := model.NewPostList()
	channelID := "channel1"
	post1 := model.Post{Id: "post1", UserId: "user1", ChannelId: channelID}
	post2 := model.Post{Id: "post2", UserId: "user1", ChannelId: channelID}
	post3 := model.Post{Id: "post3", UserId: "user2", ChannelId: channelID}

	postList.AddPost(&post1)
	postList.AddOrder("post1")
	postList.AddPost(&post2)
	postList.AddOrder("post2")
	postList.AddPost(&post3)
	postList.AddOrder("post3")

	rankings := make(userChannelRank)
	rankings["user1"] = map[string]int64{"channel1": 2}
	rankings["user2"] = map[string]int64{"channel1": 1}
	return postList, rankings
}
func TestGetRankingsSinceForChannel(t *testing.T) {
	assert := assert.New(t)
	t.Run("GetPostsSince error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetPostsSince", mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRankingsSinceForChannel("", 0)
		assert.NotNil(err)
	})

	t.Run("No error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		channelID := "channel1"
		postList, correctRanks := createMockPostList()
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		plugin.SetAPI(api)
		ranks, err := plugin.getRankingsSinceForChannel(channelID, 0)
		assert.Nil(err)
		assert.Equal(correctRanks, ranks)
	})
}

func TestGetRankingSinceForTeam(t *testing.T) {
	assert := assert.New(t)
	t.Run("GetPublicChannelsForTeam error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetPublicChannelsForTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRankingSinceForTeam("", 0)
		assert.NotNil(err)
	})

	t.Run("GetPublicChannelsForTeam 0 channels", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		channels := make([]*model.Channel, 0)
		api.On("GetPublicChannelsForTeam", mock.Anything, mock.Anything, mock.Anything).Return(channels, nil)
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		ranks, err := plugin.getRankingSinceForTeam("", 0)
		assert.Nil(err)
		assert.Equal(0, len(ranks))
	})

	t.Run("GetPublicChannelsForTeam many channels, GetPostsSince error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: "channelId"}
		api.On("GetPublicChannelsForTeam", mock.Anything, mock.Anything, mock.Anything).Return(channels, nil)
		api.On("GetPostsSince", mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRankingSinceForTeam("", 0)
		assert.NotNil(err)
	})

	t.Run("GetPublicChannelsForTeam many channels, no error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		channelID := "channel1"
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: channelID}
		postList, correctRanks := createMockPostList()
		api.On("GetPublicChannelsForTeam", mock.Anything, 0, mock.Anything).Return(channels, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, 1, mock.Anything).Return(make([]*model.Channel, 0), nil)
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		ranks, err := plugin.getRankingSinceForTeam("", 0)
		assert.Nil(err)
		assert.Equal(correctRanks, ranks)
	})
}

func TestGetRankingSince(t *testing.T) {
	assert := assert.New(t)
	t.Run("GetTeams error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRankingSince(0)
		assert.NotNil(err)
	})

	t.Run("getRankingSinceForTeam error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRankingSince(0)
		assert.NotNil(err)
	})

	t.Run("No error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		channelID := "channel1"
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: channelID}
		postList, correctRanks := createMockPostList()
		api.On("GetPublicChannelsForTeam", mock.Anything, 0, mock.Anything).Return(channels, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, 1, mock.Anything).Return(make([]*model.Channel, 0), nil)
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		ranks, err := plugin.getRankingSince(0)
		assert.Nil(err)
		assert.Equal(correctRanks, ranks)
	})
}

func TestGetRanking(t *testing.T) {
	assert := assert.New(t)
	t.Run("retreiveTimestamp error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("KVGet", timestampKey).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRanking()
		assert.NotNil(err)
	})

	t.Run("getRankingSince error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("KVGet", timestampKey).Return([]byte(`0`), nil)
		api.On("GetTeams").Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRanking()
		assert.NotNil(err)
	})

	t.Run("retreiveUserChannelRanks error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("KVGet", timestampKey).Return([]byte(`0`), nil)
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		channelID := "channel1"
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: channelID}
		postList, _ := createMockPostList()
		api.On("GetPublicChannelsForTeam", mock.Anything, 0, mock.Anything).Return(channels, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, 1, mock.Anything).Return(make([]*model.Channel, 0), nil)
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		api.On("KVGet", userChannelRanksKey).Return(nil, model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRanking()
		assert.NotNil(err)
	})

	t.Run("saveTimestamp error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("KVGet", timestampKey).Return([]byte(`0`), nil)
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		channelID := "channel1"
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: channelID}
		postList, _ := createMockPostList()
		api.On("GetPublicChannelsForTeam", mock.Anything, 0, mock.Anything).Return(channels, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, 1, mock.Anything).Return(make([]*model.Channel, 0), nil)
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		um := make(userChannelRank)
		um["user10"] = map[string]int64{"chan": 100}
		j, _ := json.Marshal(um)
		api.On("KVGet", userChannelRanksKey).Return(j, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(model.NewAppError("", "", nil, "", 404))
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		_, err := plugin.getRanking()
		assert.NotNil(err)
	})

	t.Run("no error", func(t *testing.T) {
		api := &plugintest.API{}
		plugin := Plugin{}
		api.On("KVGet", timestampKey).Return([]byte(`0`), nil)
		teams := make([]*model.Team, 1)
		teams[0] = &model.Team{Id: "teamID"}
		api.On("GetTeams").Return(teams, nil)
		channelID := "channel1"
		channels := make([]*model.Channel, 1)
		channels[0] = &model.Channel{Id: channelID}
		postList, correct := createMockPostList()
		api.On("GetPublicChannelsForTeam", mock.Anything, 0, mock.Anything).Return(channels, nil)
		api.On("GetPublicChannelsForTeam", mock.Anything, 1, mock.Anything).Return(make([]*model.Channel, 0), nil)
		api.On("GetPostsSince", channelID, mock.Anything).Return(postList, nil)
		um := make(userChannelRank)
		um["user10"] = map[string]int64{"chan": 100}
		j, _ := json.Marshal(um)
		api.On("KVGet", userChannelRanksKey).Return(j, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything)
		plugin.SetAPI(api)
		ranks, err := plugin.getRanking()
		assert.Nil(err)
		rankingUnion(correct, um)
		assert.Equal(correct, ranks)
	})
}
