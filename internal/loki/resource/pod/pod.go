package pod

import (
	"encoding/json"
	"fmt"
	"loki/internal/loki/collector"
	container "loki/internal/loki/resource/Container"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Pod struct {
	Namespace string
	logDir    string
}

func handleContainerLog(resp []byte, dir string) error {
	pod := make(map[string]interface{})
	err := json.Unmarshal(resp, &pod)
	if err != nil {
		return err
	}

	_, ok := pod["items"]
	if !ok {
		return nil
	}

	errChan := make(chan error)
	itemFinished := make(chan int)
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
					err := c.Log()
					errChan <- err
				}
			}

			containers := spec["containers"].([]interface{})
			for _, containerInterface := range containers {
				containerElem := containerInterface.(map[string]interface{})
				containerName := containerElem["name"].(string)
				c := container.NewContainer(podName, containerName, dir)
				err := c.Log()
				errChan <- err
			}

			itemFinished <- i

		}(i, item)
	}

	for {
		select {
		case <-itemFinished:
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

	podUrl := fmt.Sprintf(config.PodUrlPattern, request.Request.Host, request.Request.Port, p.Namespace)
	resp, err := collector.CollectLog(podUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, p.logDir)
	if err != nil {
		return err
	}

	handleContainerLog(resp, p.logDir)

	return nil
}

func NewPod(dir string) *Pod {
	logDir := path.Join(dir, "pod")
	os.MkdirAll(logDir, os.ModePerm)

	return &Pod{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
