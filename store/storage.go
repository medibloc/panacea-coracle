package store

type Storage interface {
	// Data Pool Model
	UploadFileWithSgx(path string, name string, sgxSecretKey, additional, data []byte) error

	// Data Deal Model
	UploadFile(path string, name string, data []byte) error
	MakeDownloadURL(path string, name string) string
	MakeRandomFilename() string
}
