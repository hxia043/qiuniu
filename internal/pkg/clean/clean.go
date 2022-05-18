package clean

import (
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/path"
	"io/ioutil"
	"os"
	"time"
)

func Clean(dir string, interval time.Duration) error {
	if !path.IsDir(dir) {
		return fmt.Errorf("error: the clean dir %s is not a directory", dir)
	}

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

		if timeNow.Sub(logTime).Hours() > interval.Hours() {
			if err := os.RemoveAll(path.Join(dir, dirName)); err != nil {
				return err
			}
		}
	}

	return nil
}
