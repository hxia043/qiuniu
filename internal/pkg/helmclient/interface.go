package helmclient

import (
	"helm.sh/helm/v3/pkg/release"
)

type Client interface {
	List() ([]*release.Release, error)
	GetValues(name string) (map[string]interface{}, error)
	CheckCluster() error
}
