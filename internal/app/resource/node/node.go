package node

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var nodeUrlPattern string = "%s/api/v1/nodes"

type Node struct {
	host   string
	token  string
	logDir string
}

func (n *Node) Log() error {
	fmt.Println("Info: collect node log start...")

	collector := collector.NewCollector()

	nodeUrl := fmt.Sprintf(nodeUrlPattern, n.host)
	resp, err := collector.CollectLog(n.token, nodeUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, n.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect node log finished.")

	return nil
}

func NewNode(host, token, dir string) *Node {
	return &Node{
		host:   host,
		token:  token,
		logDir: path.Join(dir, "node"),
	}
}
