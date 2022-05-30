package role

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var roleUrlPattern string = "%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/roles"

type Role struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (r *Role) Log() error {
	fmt.Println("Info: collect role log start...")

	collector := collector.NewCollector()

	roleUrl := fmt.Sprintf(roleUrlPattern, r.host, r.namespace)
	resp, err := collector.CollectLog(r.token, roleUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, r.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect role log finished.")

	return nil
}

func NewRole(host, token, dir string) *Role {
	return &Role{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "role"),
	}
}
