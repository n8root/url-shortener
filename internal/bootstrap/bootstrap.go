package bootstrap

import (
	"log"
	"url-shortener/internal/config"
)

var Application struct {
	Config *config.Config
}

var cfg *config.Config

func Boot(cfgPath string) {
}

func loadConfig(cfgPath string) {
	var err error

	cfg, err = config.NewConfig(cfgPath)

	if err != nil {
		log.Fatalf("failed initialization config %v", err)
	}
}
