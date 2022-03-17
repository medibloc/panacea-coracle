package store_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/medibloc/panacea-data-market-validator/store"
)

// TestS3UploadAndDownload Upload file to s3Store and download generated url link and verify after download
func TestS3UploadAndDownload(t *testing.T) {
	bucket := os.Getenv("EDG_DATAVAL_AWS_S3_BUCKET")
	region := os.Getenv("EDG_DATAVAL_AWS_S3_REGION")
	accessKeyID := os.Getenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID")
	secretAccessKeyID := os.Getenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY")

	s3Store, err := store.NewS3Store(bucket, region, accessKeyID, secretAccessKeyID)
	require.NoError(t, err)

	path := "temp_path"
	name := s3Store.MakeRandomFilename()
	data := []byte(name)

	err = s3Store.UploadFile(path, name, data)
	require.NoError(t, err)

	downloadURL := s3Store.MakeDownloadURL(path, name)
	resp, err := http.Get(downloadURL)

	defer resp.Body.Close()
	require.NoError(t, err)

	resData, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, data, resData)
}
