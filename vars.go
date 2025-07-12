package main

var (
	AppVersion     = "Undefined"
	ConfigFilepath = "user-config.yaml"
	LogLevel       = "info"
	AppConfig      Config
)

type ServerConfig struct {
	Hostname string `yaml:"hostname"`
	Port     uint16 `yaml:"port"`
	Token    string `yaml:"token"`
}

type EndpointConfig struct {
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
}

type RestConfig struct {
	Root      string           `yaml:"root"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Rest   RestConfig   `yaml:"rest"`
}
