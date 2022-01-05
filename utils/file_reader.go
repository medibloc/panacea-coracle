package utils

import (
	"bytes"
	"io"
	"net/http"
	"panacea-data-market-validator/types"
	"path/filepath"
)

// FileReader reads json file from multipart/form-data request
func FileReader(r *http.Request) (string, error) {
	// TODO limit file size
	//r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile(types.FileKey)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	fileExtension := filepath.Ext(header.Filename)
	if fileExtension != types.DataFileFormat {
		return "", types.ErrInvalidFileFormat
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return "", err
	}
	data := buf.String()
	buf.Reset()

	return data, nil
}
