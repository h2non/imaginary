package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const ImageSourceTypeS3 = "s3"
const S3Key = "s3key"
const S3OutputKey = "outputKey"
const S3Bucket = "bucket"
const S3Region = "region"

type S3ImageSource struct {
	Config *SourceConfig
}

func NewS3ImageSource(config *SourceConfig) ImageSource {
	return &S3ImageSource{config}
}

func (s *S3ImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && r.URL.Query().Get(S3Key) != ""
}

func (s *S3ImageSource) GetImage(req *http.Request) ([]byte, error) {

	key := parseKey(req)
	bucket := parseBucket(req)
	region := parseS3Region(req)

	fmt.Print("getImage S3 - key: " + key + " bucket: "+bucket + " region: "+region)

	sess := createSession(region)

	downloader := s3manager.NewDownloader(sess)
	filename := "/tmp/" + randString() + ".png"

	f, errFile := os.Create(filename)
	if errFile != nil {
		_ = fmt.Errorf("failed to create file, %v", errFile)
	}

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    &key,
	})
	if err != nil {
		_ = fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf(" file downloaded, %d bytes\n", n)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	return buffer, err
}

func createSession(region string) *session.Session{
	profile := "default"
	path := "/go/config/credentials"

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewSharedCredentials(path, profile),
	}))
	return sess
}

func uploadBufferToS3(buffer []byte, outputKey string, bucket string, region string){
	sess := createSession(region)
	uploader := s3manager.NewUploader(sess)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    &outputKey,
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		_ = fmt.Errorf("failed to upload file, %v ", err)
	}
	fmt.Print("upload successful to "+result.Location + " ")
}



func parseKey(request *http.Request) (string) {
	return request.URL.Query().Get(S3Key)
}

func parseBucket(request *http.Request) (string) {
	return request.URL.Query().Get(S3Bucket)
}

func parseS3Region(request *http.Request) (string) {
	return request.URL.Query().Get(S3Region)
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSource)
}

func randString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}
