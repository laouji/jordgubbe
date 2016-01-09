package platform

import (
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/feed"
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

func (r *IosReviewRetriever) RetrieveEntries() ([]feed.Entry, error) {
	itunesFeed := feed.NewFeed(r.Conf.ItunesAppId)
	rawXml := HttpGet(itunesFeed.Uri)

	entries, err := itunesFeed.Entries(rawXml)

	return entries, err
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
