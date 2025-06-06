package infrastructure

import (
	"os"

	"github.com/pkg/errors"
)

type FileInfrastructure interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

type FileInfrastructureImpl struct{}

func NewFileInfrastructure() *FileInfrastructureImpl {
	return &FileInfrastructureImpl{}
}

const filePermission = 0o666

func (f *FileInfrastructureImpl) ReadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, errors.Wrap(err, "file infrastructure: failed to open file")
	}

	//nolint:gosec
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "file infrastructure: failed to read file")
	}

	return data, nil
}

func (f *FileInfrastructureImpl) WriteFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, filePermission); err != nil {
		return errors.Wrap(err, "file infrastructure: failed to write file")
	}

	return nil
}
