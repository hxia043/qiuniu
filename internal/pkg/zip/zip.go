package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"loki/internal/pkg/path"
	"os"
	"path/filepath"
	"strings"
)

var Dir string = ""

func Zip(dir string) error {
	if !path.IsDir(dir) {
		err := fmt.Errorf("loki: the workdir %s is not a directory :(", dir)
		return err
	}

	zipFileName := dir + ".zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer func() { _ = zipFile.Close() }()

	zipWriter := zip.NewWriter(zipFile)
	defer func() { _ = zipWriter.Close() }()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(dir)+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()
		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return err
	}

	fmt.Printf("loki: compress content of dir %s into zip file %s\n", dir, zipFileName)
	return nil
}
