package slack

import (
	"encoding/json"
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/model"
	"strings"
)

type SlackPayload struct {
	Text        string            `json:"text"`
	UserName    string            `json:"username"`
	IconEmoji   string            `json:"icon_emoji"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Title     string                 `json:"title"`
	TitleLink string                 `json:"title_link"`
	Text      string                 `json:"text"`
	Fallback  string                 `json:"fallback"`
	Fields    []SlackAttachmentField `json:"fields"`
}

type SlackAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func GenerateAttachments(reviews []*model.Review) []SlackAttachment {
	conf := config.LoadConfig()
	attachments := []SlackAttachment{}

	for i, review := range reviews {
		if i > conf.MaxAttachmentCount {
			break
		}

		fields := []SlackAttachmentField{}
		fields = append(fields, SlackAttachmentField{Title: "Rating", Value: strings.Repeat(":star:", review.Rating), Short: true})
		fields = append(fields, SlackAttachmentField{Title: "Updated", Value: review.Updated.Format("2006-01-02 15:04:05"), Short: true})

		attachments = append(attachments, SlackAttachment{
			Title:     review.Title,
			TitleLink: review.AuthorUri,
			Text:      review.Content,
			Fallback:  review.Title + " " + review.AuthorUri,
			Fields:    fields,
		})
	}

	return attachments
}

func PreparePayload(attachments []SlackAttachment) []byte {
	conf := config.LoadConfig()
	slackPayload := SlackPayload{
		UserName:    conf.BotName,
		IconEmoji:   conf.IconEmoji,
		Text:        conf.MessageText,
		Attachments: attachments,
	}
	payload, _ := json.Marshal(slackPayload)

	return payload
}
