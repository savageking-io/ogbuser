package main

import (
	"github.com/savageking-io/ogbrest/restlib"
	steam "github.com/savageking-io/ogbsteam/client"
	"github.com/savageking-io/ogbuser/db"
	"github.com/savageking-io/ogbuser/kafka"
	"github.com/savageking-io/ogbuser/token"
)

var (
	AppVersion     = "Undefined"
	ConfigFilepath = "user-config.yaml"
	LogLevel       = ""
	AppConfig      ServiceConfig
)

type ServiceConfig struct {
	LogLevel    string                         `yaml:"log_level"`
	Rest        restlib.RestInterServiceConfig `yaml:"rest"`
	Rpc         RpcConfig                      `yaml:"rpc"`
	Postgres    db.PostgresConfig              `yaml:"postgres"`
	Crypto      CryptoConfig                   `yaml:"crypto"`
	Kafka       kafka.Config                   `yaml:"kafka"`
	SteamClient steam.Config                   `yaml:"steam_client"`
}

type RpcConfig struct {
	Hostname string `yaml:"hostname"`
	Port     uint16 `yaml:"port"`
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
