package clean

import (
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	Interval time.Duration = 0
	Workdir  string        = ""
)

func Clean(dir string, interval time.Duration) error {
	timeNow := time.Now()

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	dirNames := make([]string, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirNames = append(dirNames, fileInfo.Name())
		}
	}

	for _, dirName := range dirNames {
		logTime, err := time.Parse(time.RFC3339, dirName)
		if err != nil {
			return err
		}

		if timeNow.Sub(logTime).Seconds() > interval.Seconds() {
			err := os.RemoveAll(path.Join(dir, dirName))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
