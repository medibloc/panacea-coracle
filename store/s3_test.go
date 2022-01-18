package store_test

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/medibloc/panacea-data-market-validator/store"
)

// TestS3UploadAndDownload Upload file to s3Store and download generated url link and verify after download
func TestS3UploadAndDownload(t *testing.T) {
	s3Store, err := store.NewDefaultS3Store()
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
