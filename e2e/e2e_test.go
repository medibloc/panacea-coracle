package e2e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateData(t *testing.T) {
	addr := os.Getenv("E2E_DATAVAL_ADDR")

	data := `{
		"name": "This is a name",
		"description": "This is a description"
	}`

	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/validate-data/1?requester_address=panacea1c7yh0ql0rhvyqm4vuwgaqu0jypafnwqdc6x60e", addr),
		strings.NewReader(data),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	//TODO: testify more
	t.Log(string(body))
}
