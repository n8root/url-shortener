package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

	serve(cfg)

}

func serve(cfg *config.Config) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, URL Shortener!")
	})

	addr := cfg.Server.Addr()

	http.ListenAndServe(addr, nil)
}
