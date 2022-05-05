package pod

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/resource/container"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var podUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/pods/"

type Pod struct {
	namespace string
	logDir    string
}

func handleContainerLog(resp []byte, dir string) error {
	pod := make(map[string]interface{})
	if err := json.Unmarshal(resp, &pod); err != nil {
		return err
	}

	_, ok := pod["items"]
	if !ok {
		return nil
	}

	errChan := make(chan error)
	itemChan := make(chan int)
	items := pod["items"].([]interface{})
	itemCount := 0
	for i, item := range items {
		go func(i int, item interface{}) {
			itemElem := item.(map[string]interface{})
			metadata := itemElem["metadata"].(map[string]interface{})
			podName := metadata["name"].(string)
			spec := itemElem["spec"].(map[string]interface{})

			if _, ok := spec["initContainers"]; ok {
				initContainers := spec["initContainers"].([]interface{})
				for _, initContainer := range initContainers {
					initContainerElem := initContainer.(map[string]interface{})
					initContainerName := initContainerElem["name"].(string)
					c := container.NewContainer(podName, initContainerName, dir)
					errChan <- c.Log()
				}
			}

			containers := spec["containers"].([]interface{})
			for _, containerInterface := range containers {
				containerElem := containerInterface.(map[string]interface{})
				containerName := containerElem["name"].(string)
				c := container.NewContainer(podName, containerName, dir)
				errChan <- c.Log()
			}

			itemChan <- i

		}(i, item)
	}

	for {
		select {
		case <-itemChan:
			itemCount += 1
			if itemCount == len(items) {
				return nil
			}
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
}

func (p *Pod) Log() error {
	collector := collector.NewCollector()

	podUrl := fmt.Sprintf(podUrlPattern, request.Request.Host, request.Request.Port, p.namespace)
	resp, err := collector.CollectLog(podUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, p.logDir); err != nil {
		return err
	}

	if err = handleContainerLog(resp, p.logDir); err != nil {
		return err
	}

	return nil
}

func NewPod(dir string) *Pod {
	return &Pod{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "pod"),
	}
}
