package main

import (
	"fmt"
	"net"
	"os"

	ogb "github.com/savageking-io/ogbcommon"
	pb "github.com/savageking-io/ogbuser/proto"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()
	app.Name = "user-service"
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
			Usage: "Start REST",
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
					Value:       LogLevel,
					Destination: &LogLevel,
				},
			},
			Action: Serve,
		},
	}

	_ = app.Run(os.Args)
}

func Serve(c *cli.Context) error {
	err := ogb.SetLogLevel(LogLevel)
	if err != nil {
		log.Errorf("Failed to set logging level: %v", err)
		return err
	}

	err = ogb.ReadYAMLConfig(ConfigFilepath, &AppConfig)
	if err != nil {
		log.Errorf("Failed to read configuration file: %v", err)
		return err
	}

	log.Infof("Configuration loaded from %s", ConfigFilepath)

	grpcServer := grpc.NewServer()
	server := &Server{
		token: AppConfig.Server.Token,
	}
	pb.RegisterUserServiceServer(grpcServer, server)

	addr := fmt.Sprintf("%s:%d", AppConfig.Server.Hostname, AppConfig.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("Failed to listen on %s: %v", addr, err)
		return err
	}

	log.Infof("Starting gRPC server on %s", addr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Errorf("Failed to serve: %v", err)
		return err
	}

	return nil
}
