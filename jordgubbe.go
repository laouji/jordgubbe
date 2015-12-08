package main

import (
	"./config"
	"./sqlite3"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ReviewData struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Id      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Title   string   `xml:"title"`
	Content []string `xml:"content"`
	Rating  int      `xml:"im:rating"`
	Author  struct {
		Name string `xml:"name"`
		Uri  string `xml:"uri"`
	} `xml:"author"`
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

	unseenEntries := FilterUnseenEntries(reviewData)
	if len(unseenEntries) <= 0 {
		//no new content
		return
	}

	attachments := ParseAttachments(unseenEntries)
	payload := PreparePayload(attachments)
	HttpPostJson(conf.WebHookUri, payload)
}

func BuildFeedUri(appId string) string {
	return "http://itunes.apple.com/jp/rss/customerreviews/id=" + appId + "/sortBy=mostRecent/xml"
}

func FilterUnseenEntries(reviewData ReviewData) []Entry {
	entries := []Entry{}

	sqlite3.Init()
	dbh := sqlite3.GetDBH()
	lastSeenReviewId := dbh.LatestId("review")

	for i, entry := range reviewData.Entries {
		// first entry is the summary of the app so skip it
		if i == 0 {
			continue
		}

		entryId, _ := strconv.Atoi(entry.Id)
		if entryId <= lastSeenReviewId {
			continue
		}

		sth, err := dbh.Prepare(`
INSERT INTO review(id, title, content, rating, author_name, author_uri, updated, acquired) 
values(?,?,?,?,?,?,?,?)
`,
		)
		if err != nil {
			panic(err)
		}

		updatedTime, _ := time.Parse("2006-01-02T15:04:05-07:00", entry.Updated)
		updated := updatedTime.Local().Format("2006-01-02 15:04:05")
		now := time.Now().Format("2006-01-02 15:04:05")

		_, err = sth.Exec(entry.Id, entry.Title, entry.Content[0], entry.Rating, entry.Author.Name, entry.Author.Uri, updated, now)
		if err != nil {
			panic(err)
		}

		entry.Updated = updated
		entries = append(entries, entry)
	}

	return entries
}

func ParseAttachments(entries []Entry) []SlackAttachment {
	attachments := []SlackAttachment{}

	for _, entry := range entries {
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
