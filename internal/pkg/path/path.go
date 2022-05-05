package path

import (
	"os"
	"path"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return s.IsDir()
}

func CreateDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func Join(dir, dirName string) string {
	return path.Join(dir, dirName)
}
