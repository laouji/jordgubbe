package platform

import (
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/feed"
	//	"google.golang.org/api/storage/v1"
)

type AndroidReviewRetriever struct {
	Conf *config.ConfData
}

func NewAndroidReviewRetriever(conf *config.ConfData) *AndroidReviewRetriever {
	return &AndroidReviewRetriever{
		Conf: conf,
	}
}

func (r *AndroidReviewRetriever) RetrieveEntries() ([]feed.Entry, error) {
	itunesFeed := feed.NewFeed(r.Conf.ItunesAppId)
	rawXml := HttpGet(itunesFeed.Uri)

	entries, err := itunesFeed.Entries(rawXml)

	return entries, err
}
