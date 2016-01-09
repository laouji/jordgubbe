package model

import (
	"github.com/laouji/jordgubbe/feed"
	"github.com/laouji/jordgubbe/middleware"
	"github.com/russross/meddler"
	"strconv"
	"time"
)

type Review struct {
	ID         int `meddler:"id"`
	Title      string
	Content    string
	Rating     int
	AuthorName string    `meddler:"author_name"`
	AuthorUri  string    `meddler:"author_uri"`
	Updated    time.Time `meddler:"updated,localtime"`
	Acquired   time.Time `meddler:"acquired,localtime"`
}

func init() {
	meddler.Default = meddler.SQLite
}

func NewReview(entry *feed.Entry) *Review {
	id, _ := strconv.Atoi(entry.Id)
	updatedTime, _ := time.Parse("2006-01-02T15:04:05-07:00", entry.Updated)

	return &Review{
		ID:         id,
		Title:      entry.Title,
		Content:    entry.Content[0],
		Rating:     entry.Rating,
		AuthorName: entry.Author.Name,
		AuthorUri:  entry.Author.Uri,
		Updated:    updatedTime.Local(),
		Acquired:   time.Now(),
	}
}

func LastSeenReviewId(platformKey string) int {
	dbh := middleware.GetDBH()
	return dbh.LastInsertId("review_" + platformKey)
}

func (rowData *Review) Save(platformKey string) error {
	dbh := middleware.GetDBH()
	err := meddler.Insert(dbh, "review_"+platformKey, rowData)

	return err
}
