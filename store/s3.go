package store

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/medibloc/panacea-data-market-validator/config"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

var _ Storage = (*AWSS3Storage)(nil)

type AWSS3Storage struct {
	bucket          string
	region          string
	accessKeyID     string
	secretAccessKey string
}

// NewS3Store Create AWSS3Storage with bucket and region.
func NewS3Store(conf *config.Config) (Storage, error) {
	if conf.AWSS3.Bucket == "" {
		return nil, fmt.Errorf("'bucket' should not be empty")
	}
	if conf.AWSS3.Region == "" {
		return nil, fmt.Errorf("'region' should not be empty")
	}
	if conf.AWSS3.AccessKeyID == "" {
		return nil, fmt.Errorf("'accessKeyID' should not be empty")
	}
	if conf.AWSS3.SecretAccessKey == "" {
		return nil, fmt.Errorf("'secretAccessKey' should not be empty")
	}

	return AWSS3Storage{
		bucket:          conf.AWSS3.Bucket,
		region:          conf.AWSS3.Region,
		accessKeyID:     conf.AWSS3.AccessKeyID,
		secretAccessKey: conf.AWSS3.SecretAccessKey,
	}, nil
}

// UploadFile path is directory, name is the file name.
// It is stored in the 'data-market' bucket
func (s AWSS3Storage) UploadFile(path string, name string, data []byte) error {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(s.region),
				// There are several ways to set credit.
				// By default, the SDK detects AWS credentials set in your environment and uses them to sign requests to AWS
				// AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN(optionals)
				// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
				Credentials: credentials.NewStaticCredentials(s.accessKeyID, s.secretAccessKey, ""),
			},
		),
	)
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
func (s AWSS3Storage) MakeDownloadURL(path string, name string) string {
	return fmt.Sprintf("https://%v.s3.%v.amazonaws.com/%v", s.bucket, s.region, makeFullPath(path, name))
}

// MakeRandomFilename Create filename with UUID
func (s AWSS3Storage) MakeRandomFilename() string {
	return uuid.New().String()
}

// makeFullPath simple make path
func makeFullPath(str ...string) string {
	return strings.Join(str, "/")
}
