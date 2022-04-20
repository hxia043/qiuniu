package config

var ResourceConfig *config = new(config)

var (
	PodUrlPattern            string = "https://%s:%s/api/v1/namespaces/%s/pods/"
	StorageclassUrlPattern   string = "https://%s:%s/apis/storage.k8s.io/v1/storageclasses"
	ConfigmapUrlPattern      string = "https://%s:%s/api/v1/namespaces/%s/configmaps"
	DeploymentUrlPattern     string = "https://%s:%s/apis/apps/v1/namespaces/%s/deployments/"
	EndpointUrlPattern       string = "https://%s:%s/api/v1/namespaces/%s/endpoints"
	EventUrlPattern          string = "https://%s:%s/api/v1/namespaces/%s/events/"
	ImagestreamUrlPattern    string = "https://%s:%s/apis/image.openshift.io/v1/namespaces/%s/imagestreams"
	JobUrlPattern            string = "https://%s:%s/apis/batch/v1/namespaces/%s/jobs"
	NodeUrlPattern           string = "https://%s:%s/api/v1/nodes"
	PvUrlPattern             string = "https://%s:%s/api/v1/persistentvolumes/"
	PvcUrlPattern            string = "https://%s:%s/api/v1/namespaces/%s/persistentvolumeclaims"
	RoleUrlPattern           string = "https://%s:%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/roles"
	RolebindingUrlPattern    string = "https://%s:%s/apis/rbac.authorization.k8s.io/v1/namespaces/%s/rolebindings"
	ServiceaccountUrlPattern string = "https://%s:%s/api/v1/namespaces/%s/serviceaccounts"
	StatefulsetUrlPattern    string = "https://%s:%s/apis/apps/v1/namespaces/%s/statefulsets/"
	ServiceUrlPattern        string = "https://%s:%s/api/v1/namespaces/%s/services"
	ContainerUrlPattern      string = "https://%s:%s/api/v1/namespaces/%s/pods/%s/log"
)

type config struct {
	Namespace  string
	Workdir    string
	KubeConfig string
}
