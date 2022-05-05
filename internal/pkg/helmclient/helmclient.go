package helmclient

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/pkg/logging"
	"io"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

var (
	defaultCachePath            = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".cache/helm")
	defaultRepositoryConfigPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".config/helm")
)

// client defines the struct of a helm client
type client struct {
	Settings       *cli.EnvSettings
	ActionConfig   *action.Configuration
	ReleaseTesting ReleaseTester
}

type ReleaseTester interface {
	Run(name string) (*release.Release, error)
	GetPodLogs(out io.Writer, rel *release.Release) error
}

// Options defines the options for helm client's configuration
type Options struct {
	RepositoryConfigPath string
	RepositoryCachePath  string
	Debug                bool
	Kubeconfig           string
	Namespace            string
	Log                  action.DebugLog
}

// ReleaseSpec specify properties of release when installed from a chart
type ReleaseSpec struct {
	ReleaseName     string
	ValuesYaml      string
	Wait            bool
	CreateNamespace bool
}

func NewOptions(kubeconfig string) *Options {
	return &Options{
		Kubeconfig: kubeconfig,
		Log:        logging.NewLogger().Infof,
		Namespace:  config.Config.Namespace,
	}
}

// New returns a new Helm client with the provided options
func NewHelmClient(options *Options) (Client, error) {
	if options == nil {
		return nil, errors.New("error: helm options argument is nil")
	}

	if options.Kubeconfig == "" {
		return nil, errors.New("error: kubeconfig missing")
	}

	kubeconfig, err := base64.StdEncoding.DecodeString(options.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error: decode kubeconfig failed: %w", err)
	}

	clientGetter, err := NewRESTClientGetter(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("new rest client getter failed: %w", err)
	}

	ns, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, fmt.Errorf("get namespace from KubeConfig fail: %w", err)
	}

	// Helm does not expose Namespace property, it has to be configured in this environment way
	if options.Namespace != "" {
		ns = options.Namespace
	}

	if err := os.Setenv("HELM_NAMESPACE", ns); err != nil {
		return nil, fmt.Errorf("setting namespace fail, %w", err)
	}

	settings := cli.New()
	setEnvSettings(options, settings)

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(clientGetter, settings.Namespace(), os.Getenv("HELM_DRIVER"), options.Log)
	if err != nil {
		return nil, err
	}

	return &client{
		Settings:     settings,
		ActionConfig: actionConfig,
	}, nil
}

// set environment variables for helm
func setEnvSettings(options *Options, settings *cli.EnvSettings) {
	if options.RepositoryConfigPath == "" {
		options.RepositoryConfigPath = defaultRepositoryConfigPath
	}

	if options.RepositoryCachePath == "" {
		options.RepositoryCachePath = defaultCachePath
	}

	settings.RepositoryCache = options.RepositoryCachePath
	settings.RepositoryConfig = options.RepositoryConfigPath
	settings.Debug = options.Debug
}

// Test for origin k8s cluster
func (c *client) CheckCluster() error {
	return c.ActionConfig.KubeClient.IsReachable()
}

func (c *client) GetValues(name string) (map[string]interface{}, error) {
	getValues := action.NewGetValues(c.ActionConfig)
	return getValues.Run(name)
}

// List is to get the release name in specific namespace
// same as "helm list -n [namespace]"
func (c *client) List() ([]*release.Release, error) {
	list := action.NewList(c.ActionConfig)
	return list.Run()
}
