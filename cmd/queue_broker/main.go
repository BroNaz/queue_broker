package main

import (
	"flag"
	"log"

	"github.com/BroNaz/queue_broker/internal/app/server"
	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.toml", "path to toml config")
}

func main() {
	flag.Parse()

	config := server.NewConfiger()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(config)
	if err := s.Start(); err != nil {
		log.Fatal()
	}
}
