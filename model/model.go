package model

import (
	"github.com/laouji/jordgubbe/middleware"
	"github.com/russross/meddler"
	"log"
	"strconv"
	"time"
)

func init() {
	meddler.Default = meddler.SQLite
}

type Review struct {
	ID         int64 `meddler:"id"`
	Title      string
	Content    string
	Rating     int
	DeviceType int       `meddler:"device_type"`
	DeviceName string    `meddler:"device_name"`
	AuthorName string    `meddler:"author_name"`
	AuthorUri  string    `meddler:"author_uri"`
	Created    time.Time `meddler:"created,localtime"`
	Updated    time.Time `meddler:"updated,localtime"`
	Acquired   time.Time `meddler:"acquired,localtime"`
}

func (rowData *Review) Save() error {
	dbh := middleware.GetDBH()
	err := meddler.Insert(dbh, "review", rowData)

	return err
}

type ReviewSlice []*Review

func (s ReviewSlice) Len() int           { return len(s) }
func (s ReviewSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ReviewSlice) Less(i, j int) bool { return s[i].ID < s[j].ID }

type FeedEntry struct {
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

type CsvRow struct {
	ID         int64  `csv:Review Submit Millis Since Epoch`
	Title      string `csv:Review Title`
	Content    string `csv:Review Text`
	Rating     int    `csv:Star Rating`
	DeviceName string `csv:Device`
	AuthorUri  string `csv:Review Link`
	Created    string `csv:Review Submit Date and Time`
	Updated    string `csv:Review Update Date and Time`
}

func LastSeenReviewId(deviceType int) int64 {
	dbh := middleware.GetDBH()
	row := dbh.QueryRow(`SELECT id FROM review WHERE device_type = ` + strconv.Itoa(deviceType) + ` ORDER BY id DESC LIMIT 1`)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0
		}
		log.Fatal(err)
	}

	return id
}
