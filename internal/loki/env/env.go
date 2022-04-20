package env

var Env map[string]string = make(map[string]string)

const (
	ENV_HOST      string = "LOKI_HOST"
	ENV_PORT      string = "LOKI_PORT"
	ENV_TOKEN     string = "LOKI_TOKEN"
	ENV_NAMESPACE string = "LOKI_NAMESPACE"
	ENV_WORKSPACE string = "LOKI_WORKSPACE"
)
