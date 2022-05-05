package collector

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

type Collector struct {
	Dir string
}

func (c *Collector) CollectLog(url string) ([]byte, error) {
	resp, err := request.Handler(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Collector) GenerateContainerLog(resp []byte, path string) error {
	f := file.NewFile(path, resp)
	if err := f.WriteFile(); err != nil {
		return err
	}

	return nil
}

func (c *Collector) GenerateLog(resp []byte, dir string) error {
	data := make(map[string]interface{})
	if err := json.Unmarshal(resp, &data); err != nil {
		return err
	}

	_, ok := data["items"]
	if !ok {
		return nil
	}

	items := data["items"].([]interface{})
	for _, item := range items {
		itemElem := item.(map[string]interface{})
		metadata := itemElem["metadata"].(map[string]interface{})

		logName := metadata["name"].(string)
		logDir := path.Join(dir, logName)
		logFile := fmt.Sprintf("%s/%s.json", logDir, logName)

		text, err := json.MarshalIndent(itemElem, "", "    ")
		if err != nil {
			return err
		}

		f := file.NewFile(logFile, text)
		if err = f.WriteFile(); err != nil {
			return err
		}
	}

	return nil
}

func NewCollector() *Collector {
	return &Collector{}
}
