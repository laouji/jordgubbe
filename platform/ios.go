package platform

import (
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/factory"
	"github.com/laouji/jordgubbe/feed"
	"github.com/laouji/jordgubbe/model"
	"io/ioutil"
	"log"
	"net/http"
)

type IosReviewRetriever struct {
	Conf *config.ConfData
}

func NewIosReviewRetriever(conf *config.ConfData) *IosReviewRetriever {
	return &IosReviewRetriever{
		Conf: conf,
	}
}

func (r *IosReviewRetriever) Retrieve() []*model.Review {
	itunesFeed := feed.NewFeed(r.Conf.ItunesAppId)
	rawXml := HttpGet(itunesFeed.Uri)

	entries, err := itunesFeed.Entries(rawXml)
	if err != nil {
		log.Fatal(err)
	}

	var reviews []*model.Review
	for i, entry := range entries {
		// first entry is the summary of the app so skip it
		if i == 0 {
			continue
		}

		review := factory.NewIosReview(&entry)
		reviews = append(reviews, review)
	}

	return reviews
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
