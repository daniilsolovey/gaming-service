package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/daniilsolovey/gaming-task/internal/config"
	"github.com/daniilsolovey/gaming-task/internal/database"
	"github.com/daniilsolovey/gaming-task/internal/handler"
	"github.com/daniilsolovey/gaming-task/internal/requester"
	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `gaming-service

Handle http requests from client, send requests to game platform

Usage:
  gaming-service [options]

Options:
  -c --config <path>                Read specified config file. [default: config.yaml]
  --debug                           Enable debug messages.
  -v --version                      Print version.
  -h --help                         Show this help.
`

func main() {
	args, err := docopt.ParseArgs(
		usage,
		nil,
		"gaming-service "+version,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("version", version),
		"gaming-service started",
	)

	if args["--debug"].(bool) {
		log.SetLevel(log.LevelDebug)
	}

	log.Infof(nil, "loading configuration file: %q", args["--config"].(string))

	config, err := config.Load(args["--config"].(string))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("database", config.Database.Name),
		"connecting to the database",
	)

	database := database.NewDatabase(
		config.Database.Name, config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password,
	)
	defer database.Close()

	err = database.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	client, err := createClient(config)
	if err != nil {
		log.Fatal(err)
	}

	requester := requester.NewRequester(config, client)

	newHandler := handler.NewHandler(database, config, requester)
	newHandler.StartServer(config)
}

func createClient(config *config.Config) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(config.SSLPathPem, config.SSLPathKey)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to read ssl by path:% s",
			config.SSLPathPem,
		)

	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},

		Timeout: 1 * time.Minute,
	}

	return client, nil
}
