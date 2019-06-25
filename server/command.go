package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	trigger = "suggest"
	desc    = "Mattermost Suggestions Plugin"
)

const commandHelp = `
* |/suggest info| - Show user info
* |/suggest channels| - Show relevant channels for the user
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
		text := "###### " + desc + " - Slash Command Help\n" + strings.Replace(commandHelp, "|", "`", -1)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}
	return nil, nil
}
