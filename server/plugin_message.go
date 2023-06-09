package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/mattermost/mattermost-server/v6/model"
)

func (p *Plugin) translatePluginMessage(userID string, channelID string, post *model.Post) *model.AppError {
	configuration := p.getConfiguration()
	userInfo, _ := p.getUserInfo(userID)

	if configuration.disabled || !userInfo.Activated {
		return nil
	}

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
		Text:               &post.Message,
	}

	output, awsErr := svc.Text(&input)
	if awsErr != nil {
		return model.NewAppError(awsErr.Error(), channelID, nil, "error in translation", 400)
	}

	p.API.UpdatePost(&model.Post{
		Id:      post.Id,
		Message: *output.TranslatedText,
	})

	// for _, user := range channel {
	// 	userInfo, err := p.getUserInfo(user.Id)

	// 	if err != nil {
	// 		return model.NewAppError(err.ID, channelID, nil, err.Message, 400)
	// 	}

	// 	if user.Id != userID && userInfo.Activated {
	// 		configuration := p.getConfiguration()
	// 		sess := session.Must(session.NewSession())
	// 		creds := credentials.NewStaticCredentials(configuration.AWSAccessKeyID, configuration.AWSSecretAccessKey, "")
	// 		_, awsErr := creds.Get()
	// 		if awsErr != nil {
	// 			return model.NewAppError(awsErr.Error(), channelID, nil, "error in credentials", 400)
	// 		}

	// 		svc := translate.New(sess, aws.NewConfig().WithCredentials(creds).WithRegion(configuration.AWSRegion))

	// 		input := translate.TextInput{
	// 			SourceLanguageCode: &userInfo.SourceLanguage,
	// 			TargetLanguageCode: &userInfo.TargetLanguage,
	// 			Text:               &msg,
	// 		}

	// 		output, awsErr := svc.Text(&input)
	// 		if awsErr != nil {
	// 			return model.NewAppError(awsErr.Error(), channelID, nil, "error in translation", 400)
	// 		}

	// 		p.API.SendEphemeralPost(user.Id, &model.Post{
	// 			UserId:    p.botID,
	// 			ChannelId: channelID,
	// 			Message:   *output.TranslatedText,
	// 		})
	// 	}
	// }

	return nil
}
