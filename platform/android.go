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
	"os/exec"
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

	dateStr := time.Now().Format("200601")
	fileName := fmt.Sprintf("reviews_%s_%s.csv", r.Conf.AndroidPackageName, dateStr)

	var reviews []*model.Review
	csvParser := r.BuildCSVParser(fileName)
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

		id, err := strconv.ParseInt(rawData[5], 10, 64)
		if err != nil {
			//column name is not a valid integer so will be skipped
			continue
		}
		rating, _ := strconv.Atoi(rawData[8])
		csvRow := model.CsvRow{
			ID:         id,
			Title:      rawData[9],
			Content:    rawData[10],
			Rating:     rating,
			DeviceName: rawData[3],
			AuthorUri:  rawData[14],
			Created:    rawData[4],
			Updated:    rawData[6],
		}
		review := factory.NewAndroidReview(&csvRow)
		reviews = append(reviews, review)
	}

	return reviews
}

func (r *AndroidReviewRetriever) Download() {
	dateStr := time.Now().Format("200601")
	fileName := fmt.Sprintf("reviews_%s_%s.csv", r.Conf.AndroidPackageName, dateStr)
	gsString := fmt.Sprintf("gs://%s/reviews/%s", r.Conf.GCSBucketId, fileName)
	gsutil := exec.Command("gsutil", "cp", gsString, r.Conf.TmpDir+"/")
	err := gsutil.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *AndroidReviewRetriever) BuildCSVParser(fileName string) *csv.Reader {
	rawBytes, err := ioutil.ReadFile(r.Conf.TmpDir + "/" + fileName)
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
