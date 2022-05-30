package container

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var containerUrlPattern string = "%s/api/v1/namespaces/%s/pods/%s/log"

type Container struct {
	name      string
	namespace string
	logDir    string
	podName   string
	host      string
	token     string
}

func (c *Container) collectContainerLog(isPrevious bool) error {
	collector := collector.NewCollector()

	// common means common part url for previous-container and current-container
	commonContainerUrl := fmt.Sprintf(containerUrlPattern, c.host, c.namespace, c.podName)

	logDir := ""
	containerUrl := ""
	if isPrevious {
		logDir = path.Join(c.logDir+"/previous-container", c.name)
		containerUrl = commonContainerUrl + "?previous=true&container=" + c.name
	} else {
		logDir = path.Join(c.logDir+"/current-container", c.name)
		containerUrl = commonContainerUrl + "?container=" + c.name
	}

	resp, err := collector.CollectLog(c.token, containerUrl)
	if err != nil {
		return err
	}

	logPath := fmt.Sprintf("%s/%s.log", logDir, c.name)
	if err = collector.GenerateContainerLog(resp, logPath); err != nil {
		return err
	}

	return nil
}

func (c *Container) Log() error {
	isPrevious := true
	if err := c.collectContainerLog(isPrevious); err != nil {
		return err
	}

	isPrevious = false
	if err := c.collectContainerLog(isPrevious); err != nil {
		return err
	}

	return nil
}

func NewContainer(podName, containerName, dir, host, token string) *Container {
	return &Container{
		name:      containerName,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, podName),
		podName:   podName,
		host:      host,
		token:     token,
	}
}
