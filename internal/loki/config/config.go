package config

import "time"

var AppConfig *Config = new(Config)

var (
	Version string = "0.1"
	Type    string = "Debug" // Type: Debug/Release
)

type Config struct {
	Host       string
	Port       string
	Token      string
	Command    string
	IsVerify   bool
	Namespace  string
	Version    string
	Type       string
	WorkSpace  string
	ZipDir     string
	Interval   time.Duration
	KubeConfig string
}
