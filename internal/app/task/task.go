package task

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/resource"
	"github/hxia043/qiuniu/internal/pkg/clean"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/helmclient"
	"github/hxia043/qiuniu/internal/pkg/k8splugin"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/request"
	"github/hxia043/qiuniu/internal/pkg/zip"
	"os"
	"runtime"
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
	fmt.Printf("The log collector for Kubernetes application\n\n")

	fmt.Println("Available Commands:")
	fmt.Println("  help			help for qiuniu")
	fmt.Println("  version 		print the qiuniu version information")
	fmt.Println("  env			qiuniu client env information")
	fmt.Println("  zip			compress the log")
	fmt.Println("  clean			clean the log")
	fmt.Println("  helm			collect the helm release log for application")
	fmt.Println("  log			collect the kubernetes application log")

	fmt.Printf("\nOptions:\n")
	fmt.Println("  log")
	fmt.Println("    -h, --host			kubernetes cluster hostname or ip")
	fmt.Println("    -p, --port			kubernetes cluster port")
	fmt.Println("    -t, --token			kubernetes cluster token")
	fmt.Println("    -n, --namespace		kubernetes cluster namespace")
	fmt.Println("    -w, --workspace		the workspace for qiuniu")
	fmt.Println("  helm")
	fmt.Println("    -k, --kubeconfig		local kubeconfig path of kubernetes cluster")
	fmt.Println("    -n, --namespace		kubernetes cluster namespace")
	fmt.Println("    -w, --workspace		the workspace for qiuniu")
	fmt.Println("  zip")
	fmt.Println("    -d, --dir			the dir of compress the log")
	fmt.Println("  clean")
	fmt.Println("    -w, --workspace		the workspace for qiuniu")
	fmt.Println("    -i, --interval		the time interval between log collect time and current time")

	return nil
}

func Version() error {
	_, err := fmt.Println("qiuniu version:", config.Version, config.Type, " go version:", runtime.Version())
	return err
}

func Env() error {
	for key, value := range config.Env {
		if _, err := fmt.Printf("%s=%s\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

func Zip() error {
	err := zip.Zip(config.Config.ZipDir)
	return err
}

func Clean() error {
	err := clean.Clean(config.Config.CleanDir, config.Config.Interval)
	return err
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

func parseWorkspace(workspace string) string {
	dir := ""
	if workspace != "" {
		dir = path.Join(workspace, "qiuniu")
	} else {
		defaultWorkspace := os.Getenv("HOME")
		dir = path.Join(defaultWorkspace, "qiuniu")
	}

	return dir
}

func Helm() error {
	workDir := parseWorkspace(config.Config.Workspace)
	helmDir := path.Join(workDir+"/helm", time.Now().Format(time.RFC3339))
	kubeconfig, err := k8splugin.TransferKubeconfigToBase64(config.Config.Kubeconfig)
	if err != nil {
		return err
	}

	options := helmclient.NewOptions(kubeconfig)
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
		if err := getReleaseValues(client, release.Name, helmDir); err != nil {
			return err
		}

		if err := getReleaseManifest(release, helmDir); err != nil {
			return err
		}

		if err := getReleaseInfo(release, helmDir); err != nil {
			return err
		}
	}

	return nil
}

func createWorkDir(dir string) (string, string, error) {
	describeId := time.Now().Format(time.RFC3339)
	describeDir := path.Join(dir, describeId)
	if err := os.MkdirAll(describeDir, os.ModePerm); err != nil {
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
			"namespace":    config.Config.Namespace,
			"collect_time": id,
		},
	}

	data, err := json.MarshalIndent(descriptionInfo, "", "    ")
	if err != nil {
		return err
	}

	f := file.NewFile(dir+"/qiuniu_description.json", data)
	if err = f.WriteFile(); err != nil {
		return err
	}

	return nil
}

func Log() error {
	request.InitLogRequest()

	dir := parseWorkspace(config.Config.Workspace)
	id, workdir, err := createWorkDir(dir)
	if err != nil {
		return err
	}

	createDescriptionFile(id, workdir)

	resources := resource.NewResources(workdir)
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
