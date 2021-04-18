package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vad1malekseev/cowsay-web/internal/apiserver"
	"io/ioutil"
	"log"
	"os"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "c", "configs/server.json", "Path to the server.json")
	flag.Parse()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(cfg)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	s.WaitWithGracefulShutdown()
}

func getConfig() (*apiserver.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening the config file: %v", err)
	}

	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading the config file: %v", err)
	}

	cfg := apiserver.NewConfig()

	if err := json.Unmarshal(configContent, cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return cfg, nil
}
