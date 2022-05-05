package container

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var containerUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/pods/%s/log"

type Container struct {
	Name      string
	Namespace string
	logDir    string
	PodName   string
}

func (c *Container) collectContainerLog(isPrevious bool) error {
	collector := collector.NewCollector()

	// common means common part url for previous-container and current-container
	commonContainerUrl := fmt.Sprintf(containerUrlPattern, request.Request.Host, request.Request.Port, c.Namespace, c.PodName)

	logDir := ""
	containerUrl := ""
	if isPrevious {
		logDir = path.Join(c.logDir+"/previous-container", c.Name)
		containerUrl = commonContainerUrl + "?previous=true&container=" + c.Name
	} else {
		logDir = path.Join(c.logDir+"/current-container", c.Name)
		containerUrl = commonContainerUrl + "?container=" + c.Name
	}

	resp, err := collector.CollectLog(containerUrl)
	if err != nil {
		return err
	}

	logPath := fmt.Sprintf("%s/%s.log", logDir, c.Name)
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

func NewContainer(podName, containerName, dir string) *Container {
	return &Container{
		Name:      containerName,
		Namespace: config.Config.Namespace,
		logDir:    path.Join(dir, podName),
		PodName:   podName,
	}
}
