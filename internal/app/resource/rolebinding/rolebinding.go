package rolebinding

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
)

var rolebindingUrlPattern string = "https://%s:%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/rolebindings"

type RoleBinding struct {
	namespace string
	logDir    string
}

func (r *RoleBinding) Log() error {
	collector := collector.NewCollector()

	roleBindingUrl := fmt.Sprintf(rolebindingUrlPattern, request.Request.Host, request.Request.Port, r.namespace)
	resp, err := collector.CollectLog(roleBindingUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, r.logDir); err != nil {
		return err
	}

	return nil
}

func NewRoleBinding(dir string) *RoleBinding {
	return &RoleBinding{
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "role_binding"),
	}
}
