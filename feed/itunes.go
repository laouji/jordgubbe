package feed

import (
	"encoding/xml"
)

type XMLData struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Id      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Title   string   `xml:"title"`
	Content []string `xml:"content"`
	Rating  int      `xml:"rating"`
	Author  struct {
		Name string `xml:"name"`
		Uri  string `xml:"uri"`
	} `xml:"author"`
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

func (feed *Feed) Entries(rawXml []byte) ([]Entry, error) {
	xmlData := XMLData{}
	err := xml.Unmarshal(rawXml, &xmlData)
	if err != nil {
		return nil, err
	}

	return xmlData.Entries, nil
}
