package utils

import (
	"bytes"
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/types"
	"io"
	"net/http"
	"path/filepath"
)

// ReadFormFile reads json file from multipart/form-data request
func ReadFormFile(r *http.Request) (string, error) {
	// TODO: limit file size
	//r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile(types.FileKey)
	if err != nil {
		return "", fmt.Errorf("read formfile from request failed: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// TODO: handle when filename does not exist.
	fileExtension := filepath.Ext(header.Filename)
	if fileExtension != types.DataFileFormat {
		return "", types.ErrInvalidFileFormat
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return "", fmt.Errorf("io copy failed: %w", err)
	}
	data := buf.String()
	buf.Reset()

	return data, nil
}
