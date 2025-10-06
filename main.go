package main

import (
	"github.com/savageking-io/ogbuser/token"
	"os"
	"time"

	ogb "github.com/savageking-io/ogbcommon"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ogbuser"
	app.Version = AppVersion
	app.Description = "User management service for online games"
	app.Usage = "User Microservice of OnlineGameBase ecosystem"

	app.Authors = []cli.Author{
		{
			Name:  "savageking.io",
			Email: "i@savageking.io",
		},
		{
			Name:  "Mike Savochkin (crioto)",
			Email: "mike@crioto.com",
		},
	}

	app.Copyright = "2025 (c) savageking.io. All Rights Reserved"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Start user service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "config",
					Usage:       "Configuration filepath",
					Value:       ConfigFilepath,
					Destination: &ConfigFilepath,
				},
				cli.StringFlag{
					Name:        "log",
					Usage:       "Specify logging level",
					Value:       "",
					Destination: &LogLevel,
				},
			},
			Action: Serve,
		},
	}

	_ = app.Run(os.Args)
}

func Serve(c *cli.Context) error {
	err := ogb.ReadYAMLConfig(ConfigFilepath, &AppConfig)
	if err != nil {
		log.Errorf("Failed to read configuration file: %v", err)
		return err
	}

	if LogLevel == "" && AppConfig.LogLevel != "" {
		LogLevel = AppConfig.LogLevel
	}
	if LogLevel == "" {
		LogLevel = "info"
	}
	err = ogb.SetLogLevel(LogLevel)
	if err != nil {
		log.Errorf("Failed to set logging level: %v", err)
		return err
	}

	log.Infof("Configuration loaded from %s", ConfigFilepath)

	token.SetConfig(&AppConfig.Crypto.JWT)

	service := NewService(&AppConfig)
	startedAt := time.Unix(0, 0)
	for {
		if time.Since(startedAt) < time.Duration(time.Millisecond*1000) {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if err := service.Init(); err != nil {
			log.Errorf("Failed to initialize service: %v", err)
			startedAt = time.Now()
		} else {
			break
		}
	}

	return service.Start()
}
