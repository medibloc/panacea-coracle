package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Store struct {
	bucket string
	region string
}

// NewS3Store Create S3Store with bucket and region.
func NewS3Store(bucket, region string) (S3Store, error) {
	if bucket == "" {
		return S3Store{}, fmt.Errorf("'bucket' should not be empty")
	}
	if region == "" {
		return S3Store{}, fmt.Errorf("'region' should not be empty")
	}

	return S3Store{bucket: bucket, region: region}, nil
}

// UploadFile path is directory, name is the file name.
// It is stored in the 'data-market' bucket
func (s S3Store) UploadFile(path string, name string, data []byte) error {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(s.region),
				// There are several ways to set credit.
				// By default, the SDK detects AWS credentials set in your environment and uses them to sign requests to AWS
				// AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN(optionals)
				// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
				//Credentials: credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN"),
			}))
	svc := s3.New(sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(makeFullPath(path, name)),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	})

	if err != nil {
		return err
	}

	return nil
}

// MakeDownloadURL path is directory, name is the file name.
// It is downloaded in the 'data-market' bucket
func (s S3Store) MakeDownloadURL(path string, name string) string {
	return fmt.Sprintf("https://%v.s3.%v.amazonaws.com/%v", s.bucket, s.region, makeFullPath(path, name))
}

// MakeRandomFilename Create filename with UUID
func (s S3Store) MakeRandomFilename() string {
	return uuid.New().String()
}

// makeFullPath simple make path
func makeFullPath(str ...string) string {
	return strings.Join(str, "/")
}
