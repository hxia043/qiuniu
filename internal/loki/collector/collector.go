package collector

import (
	"encoding/json"
	"fmt"
	"loki/internal/pkg/file"
	"loki/internal/pkg/request"
	"os"
	"path"
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

func (c *Collector) GenerateContainerLog(resp []byte, log string) error {
	f := file.NewFile(log, resp)
	err := f.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func (c *Collector) GenerateLog(resp []byte, dir string) error {
	ep := make(map[string]interface{})
	err := json.Unmarshal(resp, &ep)
	if err != nil {
		return err
	}

	_, ok := ep["items"]
	if !ok {
		return nil
	}

	items := ep["items"].([]interface{})
	for _, item := range items {
		itemElem := item.(map[string]interface{})
		metadata := itemElem["metadata"].(map[string]interface{})

		logDir := path.Join(dir, metadata["name"].(string))
		os.MkdirAll(logDir, os.ModePerm)
		logFile := fmt.Sprintf("%s/%s.json", logDir, metadata["name"].(string))

		data, err := json.MarshalIndent(itemElem, "", "    ")
		if err != nil {
			return err
		}

		f := file.NewFile(logFile, data)
		err = f.WriteFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewCollector() *Collector {
	return &Collector{}
}
