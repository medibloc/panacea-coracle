package utils

import (
	"encoding/json"
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/types"
	"net/http"
)

// ReadData reads data from request
func ReadData(r *http.Request) (interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, types.ErrUnsupportedMediaType
	}

	var data interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("request body decode failed: %w", err)
	}

	return data, nil
}
