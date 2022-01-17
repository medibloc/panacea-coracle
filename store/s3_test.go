package store_test

import (
	"encoding/hex"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/medibloc/panacea-data-market-validator/store"
)

// TestS3UploadAndDownload Upload file to S3 and download generated url link and verify after download
func TestS3UploadAndDownload(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	path := "temp_path"
	data := []byte("original file data")
	name := hex.EncodeToString(crypto.Hash(data))

	err := store.UploadFile(path, name, data)
	require.NoError(t, err)

	downloadURL := store.MakeDownloadURL(path, name)
	resp, err := http.Get(downloadURL)

	defer resp.Body.Close()
	require.NoError(t, err)

	resData, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, data, resData)
}