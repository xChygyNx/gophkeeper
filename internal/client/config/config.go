package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

// ConfigClient - config of gophkeeper cliernt
type ConfigClient struct {
	GRPC       string       `env:"AddressGRPC"`
	DebugLevel logrus.Level `env:"DEBUG_LEVEL" envDefault:"debug"`
	FileFolder string       `env:"DATA_FOLDER"`
	FileSize   int          `env:"FILE_SIZE"`
}

// NewConfig - creates a new instance with the configuration for the client
func NewConfig() *ConfigClient {
	configClient := ConfigClient{
		GRPC:       "localhost:8080",
		FileFolder: "./data/client_keeper",
		FileSize:   4000000,
	}
	flag.StringVar(&configClient.GRPC, "g", configClient.GRPC, "Server address")
	flag.StringVar(&configClient.FileFolder, "f", configClient.FileFolder, "File Folder")
	flag.Parse()
	err := env.Parse(&configClient)
	if err != nil {
		log.Fatal(err)
	}

	return &configClient
}
