package parser

import (
	"fmt"
	"loki/internal/loki/config"
	"loki/internal/loki/env"
	"loki/internal/loki/options"
	"loki/internal/loki/task"
	"os"
	"time"
)

type Parser struct{}

func (p *Parser) Parse() error {
	parseDefaultConfig()
	parseConfigFromEnv()
	err := parseConfigFromCommandLine()
	if err != nil {
		return err
	}

	return nil
}

func parseDefaultConfig() {
	config.AppConfig.Version = config.Version
	config.AppConfig.Type = config.Type
}

func initEnv() {
	env.Env[env.ENV_HOST] = ""
	env.Env[env.ENV_PORT] = ""
	env.Env[env.ENV_TOKEN] = ""
	env.Env[env.ENV_NAMESPACE] = ""
	env.Env[env.ENV_WORKSPACE] = ""
}

func parseConfigFromEnv() {
	initEnv()

	ip := os.Getenv(env.ENV_HOST)
	if ip != "" {
		config.AppConfig.Host = ip
		env.Env[env.ENV_HOST] = ip
	}

	port := os.Getenv(env.ENV_PORT)
	if port != "" {
		config.AppConfig.Port = port
		env.Env[env.ENV_PORT] = port
	}

	token := os.Getenv(env.ENV_TOKEN)
	if token != "" {
		config.AppConfig.Token = token
		env.Env[env.ENV_TOKEN] = token
	}

	namespace := os.Getenv(env.ENV_NAMESPACE)
	if namespace != "" {
		config.AppConfig.Namespace = namespace
		env.Env[env.ENV_NAMESPACE] = namespace
	}

	workspace := os.Getenv(env.ENV_WORKSPACE)
	if workspace != "" {
		config.AppConfig.WorkSpace = workspace
		env.Env[env.ENV_WORKSPACE] = workspace
	}
}

// if there is the same config in env and command line
// the command line will override the config from env
func parseConfigFromCommandLine() error {
	if len(os.Args) <= 1 {
		config.AppConfig.Command = task.HELP
		return nil
	}

	// define os.Args[1] as the command flag
	// define os.Args[2:] as the options flag
	command := parseCommandConfig(os.Args[1])
	if err := parseOptionsConfig(command, os.Args[2:]); err != nil {
		return err
	}

	return nil
}

func parseCommandConfig(command string) string {
	switch command {
	case task.HELP:
		config.AppConfig.Command = task.HELP
	case task.LOG:
		config.AppConfig.Command = task.LOG
	case task.VERSION:
		config.AppConfig.Command = task.VERSION
	case task.ENV:
		config.AppConfig.Command = task.ENV
	case task.ZIP:
		config.AppConfig.Command = task.ZIP
	case task.CLEAN:
		config.AppConfig.Command = task.CLEAN
	case task.HELM:
		config.AppConfig.Command = task.HELM
	default:
		config.AppConfig.Command = task.HELP
	}

	return config.AppConfig.Command
}

func parseZipOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SZIPDIR, options.LZIPDIR:
			i += 1
			config.AppConfig.ZipDir = args[i]
		default:
			err := fmt.Errorf("loki: unexpected options config for %s :(", command)
			return err
		}
	}

	return nil
}

func parseIntervalOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SINTERVAL, options.LINTERVAL:
			i += 1
			interval, err := time.ParseDuration(args[i] + "s")
			if err != nil {
				return err
			}

			config.AppConfig.Interval = interval
		case options.SWORKSPACE, options.LWORKSPACE:
			i += 1
			config.AppConfig.WorkSpace = args[i]
		default:
			err := fmt.Errorf("loki: unexpected options config for %s :(", command)
			return err
		}
	}

	return nil
}

func checkOtionsConfigAvaiable(command string, args []string) error {
	if command != task.LOG && command != task.ZIP && command != task.CLEAN && command != task.HELM {
		if len(args) > 0 {
			err := fmt.Errorf("loki: wrong options config for %s :(", command)
			return err
		}
	} else {
		if len(args)%2 != 0 {
			err := fmt.Errorf("loki: wrong options config for %s :(", command)
			return err
		}
	}

	return nil
}

func parseLogOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SHOST, options.LHOST:
			i += 1
			config.AppConfig.Host = args[i]
		case options.SPORT, options.LPORT:
			i = i + 1
			config.AppConfig.Port = args[i]
		case options.STOKEN, options.LTOKEN:
			i = i + 1
			config.AppConfig.Token = args[i]
		case options.SNAMESPACE, options.LNAMESPACE:
			i = i + 1
			config.AppConfig.Namespace = args[i]
		case options.SWORKSPACE, options.LWORKSPACE:
			i = i + 1
			config.AppConfig.WorkSpace = args[i]
		default:
			err := fmt.Errorf("loki: unexpected options config for %s :(", command)
			return err
		}
	}

	return nil
}

func parseHelmOptions(command string, args []string) error {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case options.SKUBECONFIG, options.LKUBECONFIG:
			i += 1
			config.AppConfig.KubeConfig = args[i]
		case options.SWORKSPACE, options.LWORKSPACE:
			i += 1
			config.AppConfig.WorkSpace = args[i]
		case options.SNAMESPACE, options.LNAMESPACE:
			i += 1
			config.AppConfig.Namespace = args[i]
		default:
			err := fmt.Errorf("loki: unexpected options config for %s :(", command)
			return err
		}
	}

	return nil
}

func parseOptionsConfig(command string, args []string) error {
	if err := checkOtionsConfigAvaiable(command, args); err != nil {
		return err
	}

	if command == task.ZIP {
		if err := parseZipOptions(command, args); err != nil {
			return err
		}
	}

	if command == task.CLEAN {
		if err := parseIntervalOptions(command, args); err != nil {
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

	return nil
}

func NewParser() *Parser {
	return &Parser{}
}
