package types

import (
	"archive/zip"
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type file struct {
	Name string
	Body string
}

var files = []file{
	{"file1", "This is file 1"},
	{"file2", "This is file 2"},
	{"file3", "This is file 3"},
	{"file4", "This is file 4"},
	{"file5", "This is file 5"},
	{"file6", "This is file 6"},
	{"file7", "This is file 7"},
	{"file8", "This is file 8"},
	{"file9", "This is file 9"},
}

func TestConcurrentZipWriter_ZipWrite(t *testing.T) {
	buf := new(bytes.Buffer)

	czw := NewConcurrentZipWriter(buf)

	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func(f file, wg *sync.WaitGroup) {
			defer wg.Done()

			err := czw.ZipWrite(f.Name, []byte(f.Body))
			require.NoError(t, err)
		}(f, &wg)
	}

	wg.Wait()

	err := czw.Close()
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)

	require.Len(t, reader.File, len(files))

	for _, f := range reader.File {
		rc, err := f.Open()
		require.NoError(t, err)

		bz, err := io.ReadAll(rc)
		require.NoError(t, err)

		res := file{
			Name: f.Name,
			Body: string(bz),
		}

		require.Contains(t, files, res)
	}
}
