package parser

import (
	"errors"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/options"
	"github/hxia043/qiuniu/internal/app/task"
	"os"
	"time"
)

type Parser struct{}

func (p *Parser) Parse() error {
	parseConfigFromEnv()
	err := parseConfigFromCommandLine()

	return err
}

func parseConfigFromEnv() {
	if host := os.Getenv(config.ENV_HOST); host != "" {
		config.Config.Host, config.Env[config.ENV_HOST] = host, host
	}

	if port := os.Getenv(config.ENV_PORT); port != "" {
		config.Config.Port, config.Env[config.ENV_PORT] = port, port
	}

	if token := os.Getenv(config.ENV_TOKEN); token != "" {
		config.Config.Token, config.Env[config.ENV_TOKEN] = token, token
	}

	if namespace := os.Getenv(config.ENV_NAMESPACE); namespace != "" {
		config.Config.Namespace, config.Env[config.ENV_NAMESPACE] = namespace, namespace
	}

	if workspace := os.Getenv(config.ENV_WORKSPACE); workspace != "" {
		config.Config.Workspace, config.Env[config.ENV_WORKSPACE] = workspace, workspace
	}

	if serviceIp := os.Getenv(config.ENV_SERVICE_IP); serviceIp != "" {
		config.Config.ServiceIp, config.Env[config.ENV_SERVICE_IP] = serviceIp, serviceIp
	} else {
		config.Config.ServiceIp = config.DefaultServiceIp
	}

	if servicePort := os.Getenv(config.ENV_SERVICE_PORT); servicePort != "" {
		config.Config.ServicePort, config.Env[config.ENV_SERVICE_PORT] = servicePort, servicePort
	} else {
		config.Config.ServicePort = config.DefaultServicePort
	}
}

// if there is the same config in env and command line
// the command line will override the config from env
func parseConfigFromCommandLine() error {
	if len(os.Args) <= 1 {
		task.Help()
		return errors.New("error: unexpected command options")
	}

	// define os.Args[1] as the command flag
	// define os.Args[2:] as the options flag
	command, err := parseCommandConfig(os.Args[1])
	if err != nil {
		return err
	}

	if err := parseOptionsConfig(command, os.Args[2:]); err != nil {
		return err
	}

	return nil
}

func parseCommandConfig(command string) (string, error) {
	var err error = nil

	switch command {
	case task.HELP:
		config.Config.Command = task.HELP
	case task.LOG:
		config.Config.Command = task.LOG
	case task.VERSION:
		config.Config.Command = task.VERSION
	case task.ENV:
		config.Config.Command = task.ENV
	case task.ZIP:
		config.Config.Command = task.ZIP
	case task.CLEAN:
		config.Config.Command = task.CLEAN
	case task.HELM:
		config.Config.Command = task.HELM
	case task.SERVICE:
		config.Config.Command = task.SERVICE
	default:
		task.Help()
		err = errors.New("error: unexpected command options " + command)
	}

	return config.Config.Command, err
}

func parseZipOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SZIPDIR, options.LZIPDIR:
			i += 1
			config.Config.ZipDir = args[i]
		default:
			return errors.New("error: unexpected options config for " + command)
		}
	}

	return nil
}

func parseCleanOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SINTERVAL, options.LINTERVAL:
			i += 1
			interval, err := time.ParseDuration(args[i] + "h")
			if err != nil {
				return err
			}

			config.Config.Interval = interval
		case options.SWORKSPACE, options.LWORKSPACE:
			i += 1
			config.Config.Workspace = args[i]
		default:
			return errors.New("error: unexpected options config for " + command)
		}
	}

	return nil
}

func checkOptionsConfigAvaiable(command string, args []string) error {
	if command != task.LOG && command != task.ZIP && command != task.CLEAN && command != task.HELM && command != task.SERVICE {
		if len(args) > 0 {
			return errors.New("error: wrong options config for " + command)
		}
	} else {
		if len(args)%2 != 0 {
			return errors.New("error: wrong options config for " + command)
		}
	}

	return nil
}

func parseLogOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SHOST, options.LHOST:
			i += 1
			config.Config.Host = args[i]
		case options.SPORT, options.LPORT:
			i = i + 1
			config.Config.Port = args[i]
		case options.STOKEN, options.LTOKEN:
			i = i + 1
			config.Config.Token = args[i]
		case options.SNAMESPACE, options.LNAMESPACE:
			i = i + 1
			config.Config.Namespace = args[i]
		case options.SWORKSPACE, options.LWORKSPACE:
			i = i + 1
			config.Config.Workspace = args[i]
		default:
			return errors.New("error: unexpected options config for " + command)
		}
	}

	return nil
}

func parseHelmOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SKUBECONFIG, options.LKUBECONFIG:
			i += 1
			config.Config.Kubeconfig = args[i]
		case options.SWORKSPACE, options.LWORKSPACE:
			i += 1
			config.Config.Workspace = args[i]
		case options.SNAMESPACE, options.LNAMESPACE:
			i += 1
			config.Config.Namespace = args[i]
		default:
			return errors.New("error: unexpected options config for " + command)
		}
	}

	return nil
}

func parseServiceOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SSERVICE_IP, options.LSERVICE_IP:
			i += 1
			config.Config.ServiceIp = args[i]
		case options.SSERVICE_PORT, options.LSERVICE_PORT:
			i += 1
			config.Config.ServicePort = args[i]
		default:
			return errors.New("error: unexpected options config for " + command)
		}
	}

	return nil
}

// @ToDo: The parse options will be instead of the Cobra in the future
func parseOptionsConfig(command string, args []string) error {
	if err := checkOptionsConfigAvaiable(command, args); err != nil {
		return err
	}

	if command == task.ZIP {
		if err := parseZipOptions(command, args); err != nil {
			return err
		}
	}

	if command == task.CLEAN {
		if err := parseCleanOptions(command, args); err != nil {
			return err
		}
	}

	if command == task.LOG {
		if err := parseLogOptions(command, args); err != nil {
			return err
		}
	}

	if command == task.HELM {
		if err := parseHelmOptions(command, args); err != nil {
			return err
		}
	}

	if command == task.SERVICE {
		if err := parseServiceOptions(command, args); err != nil {
			return err
		}
	}

	return nil
}

func NewParser() *Parser {
	return &Parser{}
}
