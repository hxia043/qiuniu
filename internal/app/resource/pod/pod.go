package pod

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/resource/container"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var podUrlPattern string = "%s/api/v1/namespaces/%s/pods/"

type Pod struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (p *Pod) handleContainerLog(resp []byte) error {
	fmt.Println("Info: collect container log start...")

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
					c := container.NewContainer(podName, initContainerName, p.logDir, p.host, p.token)
					errChan <- c.Log()
				}
			}

			containers := spec["containers"].([]interface{})
			for _, containerInterface := range containers {
				containerElem := containerInterface.(map[string]interface{})
				containerName := containerElem["name"].(string)
				c := container.NewContainer(podName, containerName, p.logDir, p.host, p.token)
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
				fmt.Println("Info: collect container log finished.")

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
	fmt.Println("Info: collect pod log start...")

	collector := collector.NewCollector()

	podUrl := fmt.Sprintf(podUrlPattern, p.host, p.namespace)
	resp, err := collector.CollectLog(p.token, podUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, p.logDir); err != nil {
		return err
	}

	if err = p.handleContainerLog(resp); err != nil {
		return err
	}

	fmt.Println("Info: collect pod log finished.")

	return nil
}

func NewPod(host, token, dir string) *Pod {
	return &Pod{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "pod"),
	}
}
