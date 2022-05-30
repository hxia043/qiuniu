package rolebinding

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/collector"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/path"
)

var rolebindingUrlPattern string = "%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/rolebindings"

type RoleBinding struct {
	host      string
	token     string
	namespace string
	logDir    string
}

func (r *RoleBinding) Log() error {
	fmt.Println("Info: collect role binding log start...")

	collector := collector.NewCollector()

	roleBindingUrl := fmt.Sprintf(rolebindingUrlPattern, r.host, r.namespace)
	resp, err := collector.CollectLog(r.token, roleBindingUrl)
	if err != nil {
		return err
	}

	if err = collector.GenerateLog(resp, r.logDir); err != nil {
		return err
	}

	fmt.Println("Info: collect role binding log finished.")

	return nil
}

func NewRoleBinding(host, token, dir string) *RoleBinding {
	return &RoleBinding{
		host:      host,
		token:     token,
		namespace: config.Config.Namespace,
		logDir:    path.Join(dir, "role_binding"),
	}
}
