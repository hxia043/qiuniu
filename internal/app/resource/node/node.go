package node

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var nodeUrlPattern string = "https://%s:%s/api/v1/nodes"

type Node struct {
	logDir string
}

func (n *Node) Log() error {
	collector := collector.NewCollector()

	nodeUrl := fmt.Sprintf(nodeUrlPattern, request.Request.Host, request.Request.Port)
	resp, err := collector.CollectLog(nodeUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, n.logDir); err != nil {
		return err
	}

	return nil
}

func NewNode(dir string) *Node {
	return &Node{
		logDir: path.Join(dir, "node"),
	}
}
