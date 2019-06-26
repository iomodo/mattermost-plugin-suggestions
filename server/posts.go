package main

import "github.com/mattermost/mattermost-server/model"

// GetPostsForUserInChannelAfter gets a page of posts that were posted by
// the user in the channel after the post provided.
func (p *Plugin) GetPostsForUserInChannelAfter(channelID, userID, postID string, page, perPage int) (*model.PostList, *model.AppError) {
	return nil, nil
}

// GetPostCountForUserInChannelAfter gets a post count for the user in the channel.
func (p *Plugin) GetPostCountForUserInChannelAfter(channelID, userID string) (int64, *model.AppError) {
	return 0, nil
}
