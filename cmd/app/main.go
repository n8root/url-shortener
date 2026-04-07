package main

import (
	"flag"
	"fmt"
	"log"
	"url-shortener/internal/config"
)

func main() {
	//Чтобы запускать билд с аргументами --config=...
	cfgPath := flag.String("config", "configs/local.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)

	if err != nil {
		log.Fatalf("error get config %v", err)
	}

	fmt.Printf("%v", cfg)
}
