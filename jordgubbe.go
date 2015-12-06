package main

import (
	"./config"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type ReviewData struct {
	Entries []struct {
		Id      string   `xml:"id"`
		Updated string   `xml:"updated"`
		Title   string   `xml:"title"`
		Content []string `xml:"content"`
		Rating  int      `xml:"im:rating"`
		Author  struct {
			Name string `xml:"name"`
			Uri  string `xml:"uri"`
		} `xml:"author"`
	} `xml:"entry"`
}

type SlackPayload struct {
	Text        string            `json:"text"`
	UserName    string            `json:"username"`
	IconEmoji   string            `json:"icon_emoji"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Title     string            `json:"title"`
	TitleLink string            `json:"title_link"`
	Text      string            `json:"text"`
	Fallback  string            `json:"fallback"`
	Fields    []AttachmentField `json:"fields"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func main() {
	conf := config.LoadConfig()
	rawXml := HttpGet(BuildFeedUri(conf.ItunesAppId))

	reviewData := ReviewData{}
	err := xml.Unmarshal(rawXml, &reviewData)
	if err != nil {
		panic(err)
	}

	attachments := ParseAttachments(reviewData)
	payload := PreparePayload(attachments)
	HttpPostJson(conf.WebHookUri, payload)
}

func BuildFeedUri(appId string) string {
	return "http://itunes.apple.com/jp/rss/customerreviews/id=" + appId + "/sortBy=mostRecent/xml"
}

func ParseAttachments(reviewData ReviewData) []SlackAttachment {
	attachments := []SlackAttachment{}

	for i, entry := range reviewData.Entries {
		if i > 0 && i < 5 {

			fields := []AttachmentField{}
			fields = append(fields, AttachmentField{Title: "Reviewer", Value: entry.Author.Name, Short: true})
			fields = append(fields, AttachmentField{Title: "Updated", Value: entry.Updated, Short: true})

			attachments = append(attachments, SlackAttachment{
				Title:     entry.Title,
				TitleLink: entry.Author.Uri,
				Text:      entry.Content[0],
				Fallback:  entry.Title + " " + entry.Author.Uri,
				Fields:    fields,
			})
		}
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

func HttpGet(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	defer res.Body.Close()
}
