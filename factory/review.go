package factory

import (
	"fmt"
	"github.com/laouji/jordgubbe/model"
	"strconv"
	"time"
)

const DEVICE_TYPE_IOS = 1
const DEVICE_TYPE_ANDROID = 2

func NewIosReview(entry *model.FeedEntry) *model.Review {
	id, _ := strconv.ParseInt(entry.Id, 10, 64)
	updatedTime, _ := time.Parse("2006-01-02T15:04:05-07:00", entry.Updated)

	return &model.Review{
		ID:         id,
		Title:      entry.Title,
		Content:    entry.Content[0],
		Rating:     entry.Rating,
		DeviceType: DEVICE_TYPE_IOS,
		DeviceName: "iOS",
		AuthorName: entry.Author.Name,
		AuthorUri:  entry.Author.Uri,
		Updated:    updatedTime.Local(),
		Acquired:   time.Now(),
	}
}

func NewAndroidReview(row *model.CsvRow) *model.Review {
	updatedTime, _ := time.Parse("2006-01-02T15:04:05Z", row.Updated)
	createdTime, _ := time.Parse("2006-01-02T15:04:05Z", row.Created)

	return &model.Review{
		ID:         row.ID,
		Title:      row.Title,
		Content:    row.Content,
		Rating:     row.Rating,
		DeviceType: DEVICE_TYPE_ANDROID,
		DeviceName: row.DeviceName,
		AuthorUri:  row.AuthorUri,
		Created:    createdTime.Local(),
		Updated:    updatedTime.Local(),
		Acquired:   time.Now(),
	}
}
