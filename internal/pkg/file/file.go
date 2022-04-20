package file

import (
	"loki/internal/pkg/path"
	"os"
	"path/filepath"
)

type File struct {
	Path string
	Text []byte
}

func (f *File) WriteFile() error {
	dir := filepath.Dir(f.Path)
	if !path.IsDir(dir) {
		if err := path.CreateDir(dir); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(f.Path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	err = os.WriteFile(file.Name(), f.Text, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Exist() bool {
	_, err := os.Stat(f.Path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func NewFile(path string, text []byte) *File {
	return &File{
		Path: path,
		Text: text,
	}
}
