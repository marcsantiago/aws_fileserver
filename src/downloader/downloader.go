package downloader

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	parentPath string
	bucket     = ""          //  ADD BUCKET NAME HERE
	region     = "us-east-1" // CHANGE REGION HERE
	svc        *s3.S3
	wg         sync.WaitGroup
	d          *s3manager.Downloader
	utc        *time.Location
)

func init() {
	pwd, _ := os.Getwd()
	parentPath = filepath.Join(pwd, "/static")
	svc = s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	d = s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(region)}))
	utc, _ = time.LoadLocation("America/New_York")
}

func currentFiles() map[string]time.Time {
	files, _ := ioutil.ReadDir(parentPath)
	fileStats := make(map[string]time.Time)
	for _, f := range files {
		fileStats[f.Name()] = f.ModTime().In(utc)
	}
	return fileStats
}

// SyncFiles ...
func syncFiles(key string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create(filepath.Join(parentPath, key))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Downloading", key)
	d.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
}

// SyncFiles ...
func SyncFiles() {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Fatalln(err)
	}
	compFiles := currentFiles()

	for _, key := range resp.Contents {
		if t, ok := compFiles[*key.Key]; ok {
			delete(compFiles, *key.Key)
			if key.LastModified.In(utc).After(t) {
				log.Printf("Fetching the new version of %s\n", *key.Key)
				wg.Add(1)
				go syncFiles(*key.Key, &wg)
			}
		} else {
			log.Printf("Fetching the new file %s\n", *key.Key)
			wg.Add(1)
			go syncFiles(*key.Key, &wg)
		}
	}
	wg.Wait()
	// remove files from the local machine since they werent found in aws
	if len(compFiles) > 0 {
		for key := range compFiles {
			log.Printf("Removing %s\n", key)
			os.Remove(filepath.Join(parentPath, key))
		}
	}
}
