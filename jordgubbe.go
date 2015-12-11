package main

import (
	"bytes"
	"encoding/json"
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/feed"
	"github.com/laouji/jordgubbe/model"
	"log"
	"io/ioutil"
	"net/http"
	"strconv"
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

var (
	conf *config.ConfData
)

func main() {
	conf = config.LoadConfig()

	itunesFeed := feed.NewFeed(conf.ItunesAppId)
	rawXml := HttpGet(itunesFeed.Uri)

	entries, err := itunesFeed.Entries(rawXml)
	if err != nil {
		log.Fatal(err)
	}

	unseenReviews := SaveUnseen(entries)
	if len(unseenReviews) <= 0 {
		//no new content
		return
	}

	attachments := GenerateAttachments(unseenReviews)
	payload := PreparePayload(attachments)
	HttpPostJson(conf.WebHookUri, payload)
}

func SaveUnseen(entries []feed.Entry) []*model.Review {
	reviews := []*model.Review{}

	lastSeenReviewId := model.LastSeenReviewId()

	for i, entry := range entries {
		// first entry is the summary of the app so skip it
		if i == 0 {
			continue
		}

		entryId, _ := strconv.Atoi(entry.Id)
		if entryId <= lastSeenReviewId {
			break
		}

		review := model.NewReview(&entry)
		err := review.Save()
		if err != nil {
			log.Fatal(err)
		}
		reviews = append(reviews, review)
	}

	return reviews
}

func GenerateAttachments(reviews []*model.Review) []SlackAttachment {
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
	slackPayload := SlackPayload{
		UserName:    conf.BotName,
		IconEmoji:   conf.IconEmoji,
		Text:        conf.MessageText,
		Attachments: attachments,
	}
	payload, _ := json.Marshal(slackPayload)

	return payload
}

func HttpGet(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return body
}

func HttpPostJson(url string, jsonPayload []byte) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
}
