package container

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Container struct {
	Name      string
	Namespace string
	logDir    string
	PodName   string
}

func (c *Container) collectContainerLog(isPrevious bool) error {
	collector := collector.NewCollector()

	// common means common part url for previous-container and current-container
	commonContainerUrl := fmt.Sprintf(config.ContainerUrlPattern, request.Request.Host, request.Request.Port, c.Namespace, c.PodName)

	containerUrl := ""
	logDir := ""
	if isPrevious {
		logDir = path.Join(c.logDir, "previous-container", c.Name)
		containerUrl = commonContainerUrl + "?previous=true&container=" + c.Name
	} else {
		logDir = path.Join(c.logDir, "current-container", c.Name)
		containerUrl = commonContainerUrl + "?container=" + c.Name
	}

	resp, err := collector.CollectLog(containerUrl)
	if err != nil {
		return err
	}

	os.MkdirAll(logDir, os.ModePerm)
	logFile := fmt.Sprintf("%s/%s.log", logDir, c.Name)
	err = collector.GenerateContainerLog(resp, logFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) Log() error {
	isPrevious := true
	err := c.collectContainerLog(isPrevious)
	if err != nil {
		return err
	}

	isPrevious = false
	err = c.collectContainerLog(isPrevious)
	if err != nil {
		return err
	}

	return nil

}

func NewContainer(podName, containerName, dir string) *Container {
	return &Container{
		Name:      containerName,
		Namespace: config.ResourceConfig.Namespace,
		logDir:    path.Join(dir, podName),
		PodName:   podName,
	}
}
