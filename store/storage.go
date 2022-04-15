package store

type Storage interface {
	UploadFile(path string, name string, data []byte) error
	MakeDownloadURL(path string, name string) string
	MakeRandomFilename() string
}
