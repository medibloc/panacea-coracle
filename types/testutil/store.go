package testutil

import (
	"errors"
	"github.com/google/uuid"
	"github.com/medibloc/panacea-oracle/store"
	"strings"
)

var _ store.Storage = (*MockStore)(nil)

type MockStore struct {
	storeMap map[string][]byte
}

func NewMockStore() MockStore {
	return MockStore{
		make(map[string][]byte),
	}
}

func (m MockStore) UploadFile(path, name string, data []byte) error {
	m.storeMap[combinesKey(path, name)] = data
	return nil
}

func (m MockStore) MakeDownloadURL(path, name string) string {
	return strings.Join([]string{path, name}, "/")
}

func (m MockStore) MakeRandomFilename() string {
	return uuid.New().String()
}

func (m MockStore) DownloadFile(path, name string) ([]byte, error) {
	data, ok := m.storeMap[combinesKey(path, name)]
	if !ok {
		return nil, errors.New("not found file")
	}
	return data, nil
}
