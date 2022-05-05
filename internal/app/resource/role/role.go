package role

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var roleUrlPattern string = "https://%s:%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/roles"

type Role struct {
	namespace string
	logDir    string
}

func (r *Role) Log() error {
	collector := collector.NewCollector()

	roleUrl := fmt.Sprintf(roleUrlPattern, request.Request.Host, request.Request.Port, r.namespace)
	resp, err := collector.CollectLog(roleUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, r.logDir); err != nil {
		return err
	}

	return nil
}

func NewRole(dir string) *Role {
	return &Role{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "role"),
	}
}
