package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/mattermost/mattermost-server/v6/model"
)

func (p *Plugin) translatePluginMessage(userID string, channelID string, msg string, channel []*model.User) *model.AppError {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return nil
	}

	for _, user := range channel {
		userInfo, _ := p.getUserInfo(user.Id)

		if user.Id != userID && userInfo.Activated {
			configuration := p.getConfiguration()
			sess := session.Must(session.NewSession())
			creds := credentials.NewStaticCredentials(configuration.AWSAccessKeyID, configuration.AWSSecretAccessKey, "")
			_, awsErr := creds.Get()
			if awsErr != nil {
				return model.NewAppError(awsErr.Error(), channelID, nil, "error in credentials", 400)
			}

			svc := translate.New(sess, aws.NewConfig().WithCredentials(creds).WithRegion(configuration.AWSRegion))

			input := translate.TextInput{
				SourceLanguageCode: &userInfo.SourceLanguage,
				TargetLanguageCode: &userInfo.TargetLanguage,
				Text:               &msg,
			}

			output, awsErr := svc.Text(&input)
			if awsErr != nil {
				return model.NewAppError(awsErr.Error(), channelID, nil, "error in translation", 400)
			}

			// bot, e := p.API.GetBot(p.getConfiguration().BotID, true)
			// p.API.LogInfo(bot.UserId)
			// if e != nil {
			// 	return e
			// }
			// p.API.LogInfo("BOTS:")
			// bots, _ := p.API.GetBots(&model.BotGetOptions{})
			// for _, b := range bots {
			// 	p.API.LogInfo(b.UserId)
			// 	p.API.LogInfo(b.Description)
			// }

			_, err := p.API.CreatePost(p.API.SendEphemeralPost(user.Id, &model.Post{
				UserId:    p.botID,
				ChannelId: channelID,
				Message:   *output.TranslatedText,
			}))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
