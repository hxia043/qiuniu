package role

import (
	"fmt"
	"loki/internal/loki/collector"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/request"
	"os"
	"path"
)

type Role struct {
	Namespace string
	logDir    string
}

func (r *Role) Log() error {
	collector := collector.NewCollector()

	roleUrl := fmt.Sprintf(config.RoleUrlPattern, request.Request.Host, request.Request.Port, r.Namespace)
	resp, err := collector.CollectLog(roleUrl)
	if err != nil {
		return err
	}

	err = collector.GenerateLog(resp, r.logDir)
	if err != nil {
		return err
	}

	return nil
}

func NewRole(dir string) *Role {
	logDir := path.Join(dir, "role")
	os.MkdirAll(logDir, os.ModePerm)

	return &Role{
		Namespace: config.ResourceConfig.Namespace,
		logDir:    logDir,
	}
}
