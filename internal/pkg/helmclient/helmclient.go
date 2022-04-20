package helmclient

import (
	"encoding/base64"
	"fmt"
	"io"
	"loki/internal/loki/resource/config"
	"loki/internal/pkg/logging"
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
	KubeConfig           string
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

func NewOptions(kubeConfig string) Options {
	return Options{
		KubeConfig: kubeConfig,
		Log:        logging.NewLogger().Infof,
		Namespace:  config.ResourceConfig.Namespace,
	}
}

// New returns a new Helm client with the provided options
func NewHelmClient(options *Options) (Client, error) {
	if options == nil {
		return nil, fmt.Errorf("options argument is nil")
	}

	if options.KubeConfig == "" {
		return nil, fmt.Errorf("kubeconfig missing")
	}

	kubeConfig, err := base64.StdEncoding.DecodeString(options.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("decode kubeconfig failed: %w", err)
	}

	clientGetter, err := NewRESTClientGetter(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("new rest client getter failed: %w", err)
	}

	ns, _, err := clientGetter.ToRawKubeConfigLoader().Namespace()
	fmt.Println(ns)
	if err != nil {
		return nil, fmt.Errorf("get namespace from KubeConfig fail: %w", err)
	}

	// Helm does not expose Namespace property, it has to be configured in this
	// environment way.
	ns = options.Namespace

	if err := os.Setenv("HELM_NAMESPACE", ns); err != nil {
		return nil, fmt.Errorf("setting namespace fail, %w", err)
	}

	settings := cli.New()
	setEnvSettings(options, settings)

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(
		clientGetter,
		settings.Namespace(),
		os.Getenv("HELM_DRIVER"),
		options.Log,
	)
	if err != nil {
		return nil, err
	}

	return &client{
		Settings:     settings,
		ActionConfig: actionConfig,
	}, nil
}

// set environment variables for helm.
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

// Install a release from a chart
func (c *client) Install(chart string, relSpec ReleaseSpec) (*release.Release, error) {
	return nil, nil
}

// Test for origin k8s cluster
func (c *client) CheckCluster() error {
	return c.ActionConfig.KubeClient.IsReachable()
}

// Get is the action for checking a given release's information
func (c *client) Get(name string) (*release.Release, error) {
	get := action.NewGet(c.ActionConfig)
	return get.Run(name)
}

func (c *client) GetValues(name string) (map[string]interface{}, error) {
	getValues := action.NewGetValues(c.ActionConfig)
	return getValues.Run(name)
}

func (c *client) GetStatus(name string) (*release.Release, error) {
	getStatus := action.NewStatus(c.ActionConfig)
	return getStatus.Run(name)
}

func (c *client) GetHistory(name string) ([]*release.Release, error) {
	getHistory := action.NewHistory(c.ActionConfig)
	return getHistory.Run(name)
}

// List is to get the release name in specific namespace
// same as "helm list -n [namespace]"
func (c *client) List() ([]*release.Release, error) {
	list := action.NewList(c.ActionConfig)
	return list.Run()
}

// Test performs "helm test" for given release
func (c *client) Test(name string, out io.Writer) error {
	return nil
}

// Uninstall delete a helm release
func (c *client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	uninstall := action.NewUninstall(c.ActionConfig)
	result, err := uninstall.Run(name)
	if err != nil {
		return nil, err
	}
	return result, nil
}
