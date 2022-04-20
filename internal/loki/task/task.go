package task

import (
	"encoding/json"
	"fmt"
	"loki/internal/loki/config"
	"loki/internal/loki/env"
	"loki/internal/loki/resource"
	rconfig "loki/internal/loki/resource/config"
	"loki/internal/loki/version"
	"loki/internal/pkg/clean"
	"loki/internal/pkg/file"
	"loki/internal/pkg/helmclient"
	"loki/internal/pkg/k8splugin"
	"loki/internal/pkg/request"
	"loki/internal/pkg/zip"
	"os"
	"path"
	"time"

	"helm.sh/helm/v3/pkg/release"
)

const (
	HELP    string = "help"
	LOG     string = "log"
	VERSION string = "version"
	ENV     string = "env"
	ZIP     string = "zip"
	CLEAN   string = "clean"
	HELM    string = "helm"

	KIND string = "logs"
)

func Help() error {
	fmt.Println("hi, I'm loki :)")

	return nil
}

func Version() error {
	_, err := fmt.Println("loki version:", version.Version.LokiVersion, version.Version.LokiType, " go version:", version.Version.GoVersion)

	return err
}

func Env() error {
	for key, value := range env.Env {
		_, err := fmt.Printf("%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func Zip() error {
	err := zip.Zip(zip.Dir)
	if err != nil {
		return err
	}

	return nil
}

func Clean() error {
	err := clean.Clean(clean.Workdir, clean.Interval)
	if err != nil {
		return err
	}

	return nil
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

	if err := file.WriteFile(); err != nil {
		return err
	}

	return nil
}

func getReleaseInfo(release *release.Release, dir string) error {
	path := fmt.Sprintf("%s/%s/release_info.txt", dir, release.Name)

	info, err := json.Marshal(release.Info)
	if err != nil {
		return err
	}

	file := file.NewFile(path, info)
	if err := file.WriteFile(); err != nil {
		return err
	}

	return nil
}

func Helm() error {
	workDir := path.Join(config.AppConfig.WorkSpace, "loki", time.Now().Format(time.RFC3339))
	kubeConfig, err := k8splugin.TransferKubeConfigToBase64(rconfig.ResourceConfig.KubeConfig)
	if err != nil {
		return err
	}

	options := helmclient.NewOptions(kubeConfig)
	client, err := helmclient.NewHelmClient(&options)
	if err != nil {
		return err
	}

	releases, err := client.List()
	if err != nil {
		return err
	}

	for _, release := range releases {
		if err := getReleaseValues(client, release.Name, workDir); err != nil {
			return err
		}

		if err := getReleaseManifest(release, workDir); err != nil {
			return err
		}

		if err := getReleaseInfo(release, workDir); err != nil {
			return err
		}
	}

	return nil
}

func createWorkDir() (string, string, error) {
	describeId := time.Now().Format(time.RFC3339)
	describeDir := path.Join(rconfig.ResourceConfig.Workdir, describeId)
	err := os.MkdirAll(describeDir, os.ModePerm)
	if err != nil {
		return "", "", err
	}

	return describeId, describeDir, nil
}

func createDescriptionFile(id, dir string) error {
	descriptionInfo := make(map[string]interface{})
	descriptionInfo["Kind"] = KIND
	descriptionInfo["Items"] = []map[string]string{
		{
			"host":         request.Request.Host,
			"port":         request.Request.Port,
			"namespace":    rconfig.ResourceConfig.Namespace,
			"collect_time": id,
		},
	}

	data, err := json.MarshalIndent(descriptionInfo, "", "    ")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", dir, "loki_description.json")
	f := file.NewFile(path, data)
	err = f.WriteFile()
	if err != nil {
		return err
	}

	return nil
}

func Log() error {
	id, dir, err := createWorkDir()
	if err != nil {
		return err
	}

	createDescriptionFile(id, dir)

	resources := resource.NewResources(dir)
	done := make(chan error)
	resourceCount := 0

	for _, res := range resources {
		go func(r resource.Resource) {
			err := r.Log()
			done <- err
		}(res)
	}

	for err := range done {
		if err != nil {
			return err
		}

		resourceCount += 1
		if resourceCount == len(resources) {
			break
		}
	}

	return nil
}
