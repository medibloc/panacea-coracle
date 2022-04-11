package store_test

import (
	"crypto/rand"
	"github.com/medibloc/panacea-data-market-validator/config"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/medibloc/panacea-data-market-validator/store"
)

func TestStorageUploadWithSgx(t *testing.T) {
	conf := &config.Config{
		AWSS3: config.AWSS3Config{
			Region:          "ap-northeast-2",
			Bucket:          "data-market-test",
			AccessKeyID:     os.Getenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY")},
	}

	s3Store, err := store.NewS3Store(conf)
	require.NoError(t, err)

	path := "temp_path"
	name := "name"

	data, err := randomBytes(100000)
	require.NoError(t, err)

	err = s3Store.UploadFileWithSgx(path, name, data)
	require.NoError(t, err)
}

// TestS3UploadAndDownload Upload file to s3Store and download generated url link and verify after download
func TestS3UploadAndDownload(t *testing.T) {
	conf := &config.Config{
		AWSS3: config.AWSS3Config{
			Region:          "ap-northeast-2",
			Bucket:          "data-market-test",
			AccessKeyID:     os.Getenv("EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY")},
	}

	s3Store, err := store.NewS3Store(conf)
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

func randomBytes(size int) ([]byte, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
