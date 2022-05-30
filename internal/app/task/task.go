package task

import (
	"encoding/json"
	"fmt"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/resource"
	"github/hxia043/qiuniu/internal/app/service"
	"github/hxia043/qiuniu/internal/pkg/clean"
	"github/hxia043/qiuniu/internal/pkg/file"
	"github/hxia043/qiuniu/internal/pkg/k8splugin"
	"github/hxia043/qiuniu/internal/pkg/path"
	"github/hxia043/qiuniu/internal/pkg/zip"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	HELP    string = "help"
	LOG     string = "log"
	VERSION string = "version"
	ENV     string = "env"
	ZIP     string = "zip"
	CLEAN   string = "clean"
	SERVICE string = "service"

	KIND string = "logs"
)

type descriptionInfo struct {
	Kind  string            `json:"Kind,omitempty"`
	Items []descriptionItem `json:"Items,omitempty"`
}

type descriptionItem struct {
	CollectTime string `json:"collect_time,omitempty"`
	Host        string `json:"host,omitempty"`
	Name        string `json:"name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
}

func Help() error {
	fmt.Printf("The log collector for Kubernetes application\n\n")

	fmt.Println("Available Commands:")
	fmt.Println("  help			print help information")
	fmt.Println("  version 		print the version information")
	fmt.Println("  env			print host env information")
	fmt.Println("  log			collect the kubernetes application log")
	fmt.Println("  zip			compress the log")
	fmt.Println("  clean			clean the log")
	fmt.Println("  service               provide the log collect service by restful api")

	fmt.Printf("\nOptions:\n")
	fmt.Println("  log")
	fmt.Println("    -n, --namespace		kubernetes cluster namespace")
	fmt.Println("    -w, --workspace		the workspace of qiuniu")
	fmt.Println("    -k, --kubeconfig		local kubeconfig path of kubernetes cluster")
	fmt.Println("  zip")
	fmt.Println("    -d, --dir			the dir of compress the log")
	fmt.Println("  clean")
	fmt.Println("    -w, --workspace		the workspace for qiuniu")
	fmt.Println("    -i, --interval		the time interval between log collect time and current time, unit(h)")
	fmt.Println("  service")
	fmt.Println("    -si, --service-ip           listening ip of qiuniu")
	fmt.Println("    -sp, --service-port         listening port of qiuniu")
	fmt.Println()

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
	dir := parseDirectory(config.Config.ZipDir)

	_, err := zip.Zip(dir)
	return err
}

func Clean() error {
	dir := parseDirectory(config.Config.CleanDir)

	err := clean.Clean(dir, config.Config.Interval)
	return err
}

func parseDirectory(dir string) string {
	if dir != "" {
		if strings.HasPrefix(dir, "~") {
			dir = strings.Replace(dir, "~", os.Getenv("HOME"), 1)
		}

		return dir
	} else {
		return path.Join(os.Getenv("HOME"), "qiuniu")
	}
}

func parseWorkspace(workspace string) string {
	if workspace != "" {
		if strings.HasPrefix(workspace, "~") {
			workspace = strings.Replace(workspace, "~", os.Getenv("HOME"), 1)
		}

		return path.Join(workspace, "qiuniu")
	} else {
		return path.Join(os.Getenv("HOME"), "qiuniu")
	}
}

func Service() error {
	s := service.New()
	if err := s.Run(); err != nil {
		return err
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

func updateDescriptionFile(id, host, workspace string) error {
	item := descriptionItem{
		Name:        id,
		Namespace:   config.Config.Namespace,
		Host:        host,
		CollectTime: id,
	}

	text, err := json.MarshalIndent(item, "", "    ")
	if err != nil {
		return err
	}

	descriptionInfo := new(descriptionInfo)
	f := file.NewFile(path.Join(workspace, "qiuniu_description.json"), text)
	if !f.Exist() {
		descriptionInfo.Kind = KIND
		descriptionInfo.Items = append(descriptionInfo.Items, item)
	} else {
		descriptionInfoBytes, err := ioutil.ReadFile(f.Path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(descriptionInfoBytes, descriptionInfo); err != nil {
			return err
		}

		descriptionInfo.Items = append(descriptionInfo.Items, item)
	}

	f.Text, err = json.MarshalIndent(descriptionInfo, "", "    ")
	if err != nil {
		return err
	}

	if err := f.WriteFile(); err != nil {
		return err
	}

	return nil
}

func Log() error {
	fmt.Println("Info: collect log start...")

	kubeconfig, err := k8splugin.TransferKubeconfigFromPath(config.Config.Kubeconfig)
	if err != nil {
		return err
	}

	namespace, host, token, err := k8splugin.ParseKubeconfig(kubeconfig)
	if err != nil {
		return err
	}

	if config.Config.Namespace == "" {
		config.Config.Namespace = namespace
	}

	workspace := parseWorkspace(config.Config.Workspace)
	id, workdir, err := createWorkDir(workspace)
	if err != nil {
		return err
	}

	updateDescriptionFile(id, host, workspace)

	resources := resource.NewResources(host, token, workdir, kubeconfig)
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
