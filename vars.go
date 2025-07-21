package main

import (
	"github.com/savageking-io/ogbrest/restlib"
	"github.com/savageking-io/ogbuser/token"
)

var (
	AppVersion     = "Undefined"
	ConfigFilepath = "user-config.yaml"
	LogLevel       = ""
	AppConfig      ServiceConfig
)

type ServiceConfig struct {
	LogLevel string                         `yaml:"log_level"`
	Rest     restlib.RestInterServiceConfig `yaml:"rest"`
	Rpc      RpcConfig                      `yaml:"rpc"`
	Postgres PostgresConfig                 `yaml:"postgres"`
	Crypto   CryptoConfig                   `yaml:"crypto"`
}

type RpcConfig struct {
	Hostname string `yaml:"hostname"`
	Port     uint16 `yaml:"port"`
}

type PostgresConfig struct {
	Hostname        string `yaml:"hostname"`
	Port            uint16 `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	SslMode         bool   `yaml:"ssl_mode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type CryptoConfig struct {
	Argon ArgonConfig  `yaml:"argon"`
	JWT   token.Config `yaml:"jwt"`
}

type ArgonConfig struct {
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt_length"`
	KeyLength   uint32 `yaml:"key_length"`
}
