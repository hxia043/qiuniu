package helmclient

import (
	"io"

	"helm.sh/helm/v3/pkg/release"
)

type Client interface {
	List() ([]*release.Release, error)
	Get(name string) (*release.Release, error)
	GetValues(name string) (map[string]interface{}, error)
	GetStatus(name string) (*release.Release, error)
	GetHistory(string) ([]*release.Release, error)
	Install(chart string, relSpec ReleaseSpec) (*release.Release, error)
	Uninstall(name string) (*release.UninstallReleaseResponse, error)
	Test(name string, out io.Writer) error
	CheckCluster() error
}
