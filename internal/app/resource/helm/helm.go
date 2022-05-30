package helm

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/helmclient"
	"path"

	"helm.sh/helm/v3/pkg/release"
)

type Helm struct {
	kubeconfig []byte
	logDir     string
}

func getReleaseValues(client helmclient.Client, name string, dir string) error {
	values, err := client.GetValues(name)
	if err != nil {
		return err
	}

	valuesLog, err := json.MarshalIndent(values, "", "    ")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s/values.json", dir, name)
	file := file.NewFile(path, valuesLog)
	if err := file.WriteFile(); err != nil {
		return err
	}

	return nil
}

func getReleaseManifest(release *release.Release, dir string) error {
	path := fmt.Sprintf("%s/%s/manifest.txt", dir, release.Name)

	file := file.NewFile(path, []byte(release.Manifest))
	return file.WriteFile()
}

func getReleaseInfo(release *release.Release, dir string) error {
	path := fmt.Sprintf("%s/%s/release_info.txt", dir, release.Name)

	info, err := json.MarshalIndent(release.Info, "", "    ")
	if err != nil {
		return err
	}

	file := file.NewFile(path, info)
	return file.WriteFile()
}

func (h *Helm) Log() error {
	fmt.Println("Info: collect helm log start...")

	options := helmclient.NewOptions(h.kubeconfig)
	client, err := helmclient.NewHelmClient(options)
	if err != nil {
		return err
	}

	if err := client.CheckCluster(); err != nil {
		return err
	}

	releases, err := client.List()
	if err != nil {
		return err
	}

	for _, release := range releases {
		if err := getReleaseValues(client, release.Name, h.logDir); err != nil {
			return err
		}

		if err := getReleaseManifest(release, h.logDir); err != nil {
			return err
		}

		if err := getReleaseInfo(release, h.logDir); err != nil {
			return err
		}
	}

	fmt.Println("Info: collect helm log finished.")

	return nil
}

func NewHelm(kubeconfig []byte, dir string) *Helm {
	return &Helm{
		kubeconfig: kubeconfig,
		logDir:     path.Join(dir, "helm"),
	}
}
