package main

import (
	"github.com/iomodo/mattermost-plugin-suggestions/server/ml"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

func mapToSlice(m map[string]*model.Channel) []*model.Channel {
	channels := make([]*model.Channel, len(m))
	index := 0
	for _, channel := range m {
		channels[index] = channel
		index++
	}
	return channels
}

func (p *Plugin) isChannelOk(channelID string) bool {
	posts, err := p.API.GetPostsForChannel(channelID, 0, 1)

	if err != nil || len(posts.Order) == 0 {
		return false
	}
	return true
}

func (p *Plugin) preCalculateRecommendations() {
	userActivity, err := p.getActivity()
	// mlog.Info(fmt.Sprintf("userActivity : %v", userActivity))
	// allChannels, _ := p.GetAllChannels()
	// mlog.Info("==============")
	// for chann := range allChannels {
	// 	mlog.Info(chann)
	// }
	// mlog.Info("==============")

	if err != nil {
		mlog.Error("Can't get user activity. " + err.Error())
		return
	}
	params := map[string]interface{}{"k": 3}
	knn := ml.NewSimpleKNN(params)
	knn.Fit(userActivity)
	for userID := range userActivity {
		recommendedChannels := make([]*recommendedChannel, 0)
		channels, appErr := p.GetAllPublicChannelsForUser(userID)
		// mlog.Info("=======" + userID + "=======")
		// for chann := range channels {
		// 	mlog.Info(chann)
		// }
		// mlog.Info("=======xxx=======")
		if appErr != nil {
			mlog.Error("Can't get public channels for user. " + appErr.Error())
			return
		}
		for _, channel := range channels {
			if _, ok := userActivity[userID][channel.Id]; !ok {
				if !p.isChannelOk(channel.Id) {
					// mlog.Info(fmt.Sprintf("channel not Ok - %v", channel.Id))
					continue
				}
				pred, err := knn.Predict(userID, channel.Id)
				// mlog.Info(fmt.Sprintf("User - %v, Channel - %v, pred - %v", userID, channel.Id, pred))

				if err != nil {
					// unknown user or unknown channel
					continue
				}
				recommendedChannels = append(recommendedChannels, &recommendedChannel{
					ChannelID: channel.Id,
					Score:     pred,
				})
			}
		}
		p.saveUserRecommendations(userID, recommendedChannels)
		// mlog.Info(fmt.Sprintf("---------userID: %v. Recommendations: %v-----", userID, len(recommendedChannels)))
		// for i := 0; i < len(recommendedChannels); i++ {
		// 	mlog.Info(fmt.Sprintf("channelID : %v, Score : %v", recommendedChannels[i].ChannelID, recommendedChannels[i].Score))
		// }
	}

}
