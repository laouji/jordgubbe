package feed

import (
	"encoding/xml"
	"github.com/laouji/jordgubbe/model"
)

type XMLData struct {
	Entries []model.FeedEntry `xml:"entry"`
}

type Feed struct {
	Uri string
}

func BuildFeedUri(appId string) string {
	return "http://itunes.apple.com/jp/rss/customerreviews/id=" + appId + "/sortBy=mostRecent/xml"
}

func NewFeed(appId string) *Feed {
	return &Feed{
		Uri: BuildFeedUri(appId),
	}
}

func (feed *Feed) Entries(rawXml []byte) ([]model.FeedEntry, error) {
	xmlData := XMLData{}
	err := xml.Unmarshal(rawXml, &xmlData)
	if err != nil {
		return nil, err
	}

	return xmlData.Entries, nil
}
