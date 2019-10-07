package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

const ImageSourceTypeS3 = "s3"
const S3Key = "key"
const S3Bucket = "bucket"

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
	fmt.Print("getImage S3 - key: " + key + " bucket: "+bucket)

	// The session the S3 Downloader will use
	//sess := session.Must(session.NewSession())

	sess := createSession()
	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)
	//buf := []byte{}
	//f := aws.NewWriteAtBuffer(buf)
	filename := "/tmp/pic.png"
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

func createSession() *session.Session{
	region := "eu-central-1"
	profile := "default"

	path := "/go/config/credentials"
	ioutil.WriteFile("/go/config/here_i_am", []byte("Hello"), 0755)
	if _, err := os.Stat(path); err != nil {
		fmt.Println("Credential File not found")
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewSharedCredentials(path, profile),
	}))
	return sess
}

func uploadBufferToS3(buffer []byte, outputKey string, bucket string){
	sess := createSession()
	uploader := s3manager.NewUploader(sess)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    &outputKey,
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		_ = fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Print("upload successful to "+result.Location)
}



func parseKey(request *http.Request) (string) {
	return request.URL.Query().Get(S3Key)
}

func parseBucket(request *http.Request) (string) {
	return request.URL.Query().Get(S3Bucket)
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSource)
}
