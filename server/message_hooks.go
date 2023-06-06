package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return
	}

	user, err := p.API.GetUser(post.UserId)
	if err != nil {
		p.API.LogError(
			"Failed to query user",
			"user_id", post.UserId,
			"error", err.Error(),
		)
		return
	}

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError(
			"Failed to query channel",
			"channel_id", post.ChannelId,
			"error", err.Error(),
		)
		return
	}

	channelUsers, err := p.API.GetUsersInChannel(post.ChannelId, "username", 0, 100)
	if err != nil {
		p.API.LogError(
			"Failed to get users in channel",
			"channel_id", post.ChannelId,
			"error", err.Error(),
		)
		return
	}

	msg := fmt.Sprintf("%s -> %s", user.Username, post.Message)
	if err := p.translatePluginMessage(user.Id, channel.Id, msg, channelUsers); err != nil {
		p.API.LogError(
			"Failed to post MessageHasBeenPosted message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}
