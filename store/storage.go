package store

type Storage interface {
	UploadFile(path string, round string, name string, data []byte) error
	MakeDownloadURL(path string, round string, name string) string
	MakeRandomFilename() string
}
