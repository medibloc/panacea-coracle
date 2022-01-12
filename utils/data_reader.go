package utils

import (
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/types"
	"io/ioutil"
	"net/http"
)

// ReadData reads data from request body
func ReadData(r *http.Request) ([]byte, error) {

	// content type check from header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, types.ErrUnsupportedMediaType
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP request body: %w", err)
	}

	return body, nil
}
