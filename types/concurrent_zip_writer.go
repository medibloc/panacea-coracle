package types

import (
	"archive/zip"
	"io"
	"sync"
)

type ConcurrentZipWriter struct {
	sync.Mutex

	Zw *zip.Writer
}

func NewConcurrentZipWriter(w io.Writer) *ConcurrentZipWriter {
	return &ConcurrentZipWriter{
		Zw: zip.NewWriter(w),
	}
}

func (czw *ConcurrentZipWriter) Close() error {
	czw.Lock()
	defer czw.Unlock()

	return czw.Zw.Close()
}

func (czw *ConcurrentZipWriter) ZipWrite(filename string, data []byte) error {
	czw.Lock()
	defer czw.Unlock()

	// TODO: add file format to filename
	zipWriter, err := czw.Zw.Create(filename + ".json")
	if err != nil {
		return err
	}

	if _, err := zipWriter.Write(data); err != nil {
		return err
	}

	return nil
}
