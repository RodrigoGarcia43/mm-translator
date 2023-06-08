package main

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {
	if err := p.registerCommands(); err != nil {
		return errors.Wrap(err, "failed to register commands")
	}

	p.client = pluginapi.NewClient(p.API, p.Driver)

	botID, err := p.client.Bot.EnsureBot(&model.Bot{
		Username:    "translator-bot",
		DisplayName: "Translator",
		Description: "A bot account created by the translator plugin",
	},
		pluginapi.ProfileImagePath("assets/R.jpeg"))
	if err != nil {
		return errors.Wrap(err, "can't ensure bot")
	}
	p.botID = botID

	return nil
}
