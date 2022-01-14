package store

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

const (
	S3Bucket = "data-market"
	S3Region = endpoints.ApNortheast2RegionID
)

// createSvc create default S3 service.
func createSvc() *s3.S3 {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(S3Region),
			Credentials: credentials.NewSharedCredentials("", "default"),
		}))

	return s3.New(sess)
}

// UploadFile path is directory, name is the file name.
// It is stored in the 'data-market' bucket
func UploadFile(path string, name string, data []byte) error {
	svc := createSvc()

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(S3Bucket),
		Key:           aws.String(makeFullPath(path, name)),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		//Expires:       aws.Time(time.Now().Local().Add(time.Second * time.Duration(10))),
	})

	if err != nil {
		return err
	}

	return nil
}

// MakeDownloadURL path is directory, name is the file name.
// It is downloaded in the 'data-market' bucket
func MakeDownloadURL(path string, name string) string {
	return fmt.Sprintf("https://%v.s3.%v.amazonaws.com/%v", S3Bucket, S3Region, makeFullPath(path, name))
}

// makeFullPath simple make path
func makeFullPath(str ...string) string {
	return strings.Join(str, "/")
}
