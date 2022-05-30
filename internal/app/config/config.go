package config

import "time"

var Config *config = new(config)
var Env map[string]string = make(map[string]string)

const (
	ENV_NAMESPACE    string = "QIUNIU_NAMESPACE"
	ENV_WORKSPACE    string = "QIUNIU_WORKSPACE"
	ENV_KUBECONFIG   string = "QIUNIU_KUBECONFIG"
	ENV_SERVICE_IP   string = "QIUNIU_SERVICE_IP"
	ENV_SERVICE_PORT string = "QIUNIU_SERVICE_PORT"
)

var (
	Version string = "1.0"

	// Type: Draft or Release
	Type string = "Release"

	DefaultServiceIp   string = "0.0.0.0"
	DefaultServicePort string = "9189"

	DescriptionFile string = ""
)

type config struct {
	Command     string
	IsVerify    bool
	Namespace   string
	Workspace   string
	ZipDir      string
	CleanDir    string
	Interval    time.Duration
	Kubeconfig  string
	ServiceIp   string
	ServicePort string
}
