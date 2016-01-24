package main

import (
	"bytes"
	"flag"
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/model"
	"github.com/laouji/jordgubbe/platform"
	"github.com/laouji/jordgubbe/slack"
	"log"
	"net/http"
)

var (
	conf *config.ConfData
)

func main() {
	flag.Parse()
	conf = config.LoadConfig()

	var retriever interface {
		Retrieve() []*model.Review
	}

	switch conf.PlatformName {
	case "android":
		retriever = platform.NewAndroidReviewRetriever(conf)
	case "ios":
		retriever = platform.NewIosReviewRetriever(conf)
	default:
		log.Fatal("unsupported platform: " + conf.PlatformName)
	}

	reviews := retriever.Retrieve()

	newReviews := FilterAndSaveReviews(reviews)
	if len(newReviews) == 0 {
		//no new content
		return
	}

	attachments := slack.GenerateAttachments(newReviews)
	payload := slack.PreparePayload(attachments)
	HttpPostJson(conf.WebHookUri, payload)
}

func FilterAndSaveReviews(candidates []*model.Review) []*model.Review {
	newReviews := []*model.Review{}

	if len(candidates) == 0 {
		return newReviews
	}

	lastSeenId := model.LastSeenReviewId(candidates[0].DeviceType)

	for _, candidate := range candidates {
		if candidate.ID <= lastSeenId {
			break
		}

		err := candidate.Save()
		if err != nil {
			log.Fatal(err)
		}
		newReviews = append(newReviews, candidate)
	}

	return newReviews
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
