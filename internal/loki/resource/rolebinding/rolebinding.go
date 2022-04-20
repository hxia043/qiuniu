package rolebinding

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type RoleBinding struct {
	Namespace string
	logDir    string
}

func (r *RoleBinding) Log() error {
	collector := collector.NewCollector()

	roleBindingUrl := fmt.Sprintf(config.RolebindingUrlPattern, request.Request.Host, request.Request.Port, r.Namespace)
	resp, err := collector.CollectLog(roleBindingUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, r.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewRoleBinding(dir string) *RoleBinding {
	logDir := path.Join(dir, "role_binding")
	os.MkdirAll(logDir, os.ModePerm)

	return &RoleBinding{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
