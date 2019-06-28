package main

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	trigger          = "suggest"
	channelAction    = "channels"
	addRandomChannel = "add"

	desc                 = "Mattermost Suggestions Plugin"
	noNewChannels        = "No new channels for you."
	addRandomChannelText = "Channel was successfully added."
)

const commandHelp = `
* |/suggest info| - Shows user info
* |/suggest channels| - Suggests relevant channels for the user
* |/suggest add| - Adds random channel to a current user. For testing only.
`

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          trigger,
		DisplayName:      "Suggestions",
		Description:      desc,
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: info, channels, help",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Type:         model.POST_DEFAULT,
	}
}

func helpResponse() (*model.CommandResponse, *model.AppError) {
	text := "###### " + desc + " - Slash Command Help\n" + strings.Replace(commandHelp, "|", "`", -1)
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
}

func appError(message string, err error) *model.AppError {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return model.NewAppError("Suggestions Plugin", message, nil, errorMessage, http.StatusBadRequest)
}

func (p *Plugin) getChannelListFromRecommendations(recommendations []*recommendedChannel) []*model.Channel {
	channels := make([]*model.Channel, 0)
	for _, rec := range recommendations {
		channel, err := p.API.GetChannel(rec.ChannelID)
		if err != nil {
			p.API.LogError("Can't get channel - "+rec.ChannelID, "err", err.Error())
			continue
		}
		channels = append(channels, channel)
	}
	return channels
}

func (p *Plugin) suggestChannelResponse(userID string) (*model.CommandResponse, *model.AppError) {
	recommendations, err := p.retreiveUserRecomendations(userID)
	if err != nil {
		return nil, appError("Can't retreive user recommendations.", err)
	}
	channels := p.getChannelListFromRecommendations(recommendations)
	if len(channels) == 0 {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, noNewChannels), nil
	}
	text := ""
	for _, channel := range channels {
		text += " * " + channel.DisplayName + " \n"
	}
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
}

func (p *Plugin) addRandomChannel(teamID, userID string) (*model.CommandResponse, *model.AppError) {
	channels, appErr := p.API.GetPublicChannelsForTeam(teamID, 0, 100)
	if appErr != nil {
		return nil, appError("Can't get channels for team for user", appErr)
	}
	randChan := channels[rand.Intn(len(channels))]
	recommendations, err := p.retreiveUserRecomendations(userID)
	if err != nil {
		return nil, appError("Can't retreive user recommendations", err)
	}
	recommend := &recommendedChannel{ChannelID: randChan.Id, Score: rand.Float64()}
	recommendations = append(recommendations, recommend)
	p.saveUserRecommendations(userID, recommendations)
	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, addRandomChannelText+" Channel name "+randChan.DisplayName), nil
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	if len(split) == 0 {
		return nil, nil
	}
	command := split[0]
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	if command != "/"+trigger {
		return &model.CommandResponse{}, nil
	}

	if action == "" || action == "help" {
		return helpResponse()
	}

	if action == channelAction {
		return p.suggestChannelResponse(args.UserId)
	}

	if action == addRandomChannel {
		return p.addRandomChannel(args.TeamId, args.UserId)
	}
	return nil, nil
}
