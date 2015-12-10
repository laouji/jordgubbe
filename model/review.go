package model

import (
	"../feed"
	"../middleware"
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

func LastSeenReviewId() int {
	dbh := middleware.GetDBH()
	return dbh.LastInsertId("review")
}

func (rowData *Review) Save() error {
	dbh := middleware.GetDBH()
	meddler.Default = meddler.SQLite
	err := meddler.Insert(dbh, "review", rowData)

	return err
}
