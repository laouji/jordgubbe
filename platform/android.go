package platform

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/laouji/jordgubbe/config"
	"github.com/laouji/jordgubbe/factory"
	"github.com/laouji/jordgubbe/model"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

type AndroidReviewRetriever struct {
	Conf *config.ConfData
}

func NewAndroidReviewRetriever(conf *config.ConfData) *AndroidReviewRetriever {
	return &AndroidReviewRetriever{
		Conf: conf,
	}
}

func (r *AndroidReviewRetriever) Retrieve() []*model.Review {
	r.CheckGsutil()
	r.Download()

	var reviews []*model.Review
	csvParser := r.BuildCSVParser()
	for {
		rawData, err := csvParser.Read()
		if err == io.EOF {
			break
		}
		if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrBareQuote {
			log.Print(err)
			continue
		}
		if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
			//EOF
			break
		}

		id, err := strconv.ParseInt(rawData[6], 10, 64)
		if err != nil {
			//column name is not a valid integer so will be skipped
			continue
		}
		rating, _ := strconv.Atoi(rawData[9])
		csvRow := model.CsvRow{
			ID:         id,
			Title:      rawData[10],
			Content:    rawData[11],
			Rating:     rating,
			DeviceName: rawData[4],
			AuthorUri:  rawData[15],
			Created:    rawData[5],
			Updated:    rawData[7],
		}
		review := factory.NewAndroidReview(&csvRow)
		reviews = append(reviews, review)
	}

	var sortedReviews model.ReviewSlice = reviews
	sort.Sort(sort.Reverse(sortedReviews[:]))

	r.Cleanup()
	return sortedReviews
}

func (r *AndroidReviewRetriever) Download() {
	gsString := fmt.Sprintf("gs://%s/reviews/%s", r.Conf.GCSBucketId, r.CsvFileName())
	gsutil := exec.Command("gsutil", "cp", gsString, r.Conf.TmpDir+"/")

	err := gsutil.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *AndroidReviewRetriever) CsvFileName() string {
	dateStr := time.Now().Format("200601")
	return fmt.Sprintf("reviews_%s_%s.csv", r.Conf.AndroidPackageName, dateStr)
}

func (r *AndroidReviewRetriever) CsvFilePath() string {
	return r.Conf.TmpDir + "/" + r.CsvFileName()
}

func (r *AndroidReviewRetriever) Cleanup() {
	file, err := os.Open(r.CsvFilePath())
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	os.Remove(file.Name())
}

func (r *AndroidReviewRetriever) BuildCSVParser() *csv.Reader {
	rawBytes, err := ioutil.ReadFile(r.CsvFilePath())
	if err != nil {
		log.Fatal(err)
	}
	utf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	transformer := unicode.BOMOverride(utf16le.NewDecoder())

	return csv.NewReader(transform.NewReader(bytes.NewReader(rawBytes), transformer))
}

func (r *AndroidReviewRetriever) CheckGsutil() {
	var out bytes.Buffer
	cmd := exec.Command("which", "gsutil")
	cmd.Stdout = &out
	cmd.Run()
	if len(out.String()) < 1 {
		log.Fatal("gsutil not found, https://cloud.google.com/storage/docs/gsutil")
	}
}
