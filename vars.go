package main

import "github.com/savageking-io/ogbrest/restlib"

var (
	AppVersion     = "Undefined"
	ConfigFilepath = "user-config.yaml"
	LogLevel       = "info"
	AppConfig      ServiceConfig
)

type ServiceConfig struct {
	Rest restlib.RestInterServiceConfig `yaml:"rest"`
	Rpc  RpcConfig                      `yaml:"rpc"`
}

type RpcConfig struct {
	Hostname string `yaml:"hostname"`
	Port     uint16 `yaml:"port"`
}
