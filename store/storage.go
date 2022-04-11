package store

type Storage interface {
	// Data Pool Model
	UploadFileWithSgx(path string, name string, data []byte) error

	// Data Deal Model
	UploadFile(path string, name string, data []byte) error
	MakeDownloadURL(path string, name string) string

	// Common
	MakeRandomFilename() string
}
