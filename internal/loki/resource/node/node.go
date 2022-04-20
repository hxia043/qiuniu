package node

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Node struct {
	logDir string
}

func (n *Node) Log() error {
	collector := collector.NewCollector()

	nodeUrl := fmt.Sprintf(config.NodeUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(nodeUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, n.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewNode(dir string) *Node {
	logDir := path.Join(dir, "node")
	os.MkdirAll(logDir, os.ModePerm)

	return &Node{
		logDir: logDir,
	}
}
