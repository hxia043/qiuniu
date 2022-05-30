package collector

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
	"io"
	"net/http"
	"time"
)

type Collector struct {
	Dir string
}

func (c *Collector) CollectLog(token, url string) ([]byte, error) {
	req, err := request.NewRequest(request.GET_REQUEST, token, url)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: request.IsSkipVerifyDefault,
			},
		},
		Timeout: time.Duration(request.Timeout) * time.Second,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
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
