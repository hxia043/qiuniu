package config

import "time"

var Config *config = new(config)
var Env map[string]string = make(map[string]string)

const (
	ENV_HOST      string = "QIUNIU_HOST"
	ENV_PORT      string = "QIUNIU_PORT"
	ENV_TOKEN     string = "QIUNIU_TOKEN"
	ENV_NAMESPACE string = "QIUNIU_NAMESPACE"
	ENV_WORKSPACE string = "QIUNIU_WORKSPACE"
)

var (
	Version string = "1.0"

	// Type: Draft or Release
	Type string = "Release"
)

type config struct {
	Host       string
	Port       string
	Token      string
	Command    string
	IsVerify   bool
	Namespace  string
	Workspace  string
	ZipDir     string
	CleanDir   string
	Interval   time.Duration
	Kubeconfig string
}

func init() {
	Env[ENV_HOST] = ""
	Env[ENV_PORT] = ""
	Env[ENV_TOKEN] = ""
	Env[ENV_NAMESPACE] = ""
	Env[ENV_WORKSPACE] = ""
}
