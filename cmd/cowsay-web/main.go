package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vad1malekseev/cowsay-web/internal/apiserver"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", "configs/server.json", "Path to the server.json")
	flag.Parse()

	cfg, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(cfg)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	s.WaitWithGracefulShutdown()
}

func getConfig(configPath string) (*apiserver.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening the config file: %w", err)
	}

	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading the config file: %w", err)
	}

	cfg := apiserver.NewConfig()

	if err := json.Unmarshal(configContent, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return cfg, nil
}
