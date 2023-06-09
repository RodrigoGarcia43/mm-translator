package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

func (p *Plugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {

	// Make sure that Feedbackbot doesn't respond to itself
	if post.UserId == p.botID {
		return
	}

	// Or to system messages
	if post.IsSystemMessage() {
		return
	}

	configuration := p.getConfiguration()
	if configuration.disabled {
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

	// Make sure this is not a post sent by another bot
	user, appErr := p.API.GetUser(post.UserId)
	if appErr != nil {
		p.API.LogError("Unable to get sender translator bot", "err", appErr)
		return
	}

	if user.IsBot {
		return
	}

	userInfo, _ := p.getUserInfo(post.UserId)

	if configuration.disabled || !userInfo.Activated {
		return
	}
	// channelUsers, err := p.API.GetUsersInChannel(post.ChannelId, "username", 0, 100)
	// if err != nil {
	// 	p.API.LogError(
	// 		"Failed to get users in channel",
	// 		"channel_id", post.ChannelId,
	// 		"error", err.Error(),
	// 	)
	// 	return
	// }

	// msg := fmt.Sprintf("%s -> %s", user.Username, post.Message)
	if err := p.translatePluginMessage(user.Id, channel.Id, post); err != nil {
		p.API.LogError(
			"Failed to post MessageHasBeenPosted message",
			"channel_id", channel.Id,
			"user_id", user.Id,
			"error", err.Error(),
		)
	}
}
