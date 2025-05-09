package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sunnyyssh/designing-software-cw2/gateway/internal/router"
	"gopkg.in/yaml.v3"
)

const configPath = "/etc/gateway/config.yaml"

func main() {
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	router := router.New(config)

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}

func readConfig(path string) (*router.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rawConfig struct {
		Locations map[string]struct{
			URL string `yaml:"url"`
		} `yaml:"locations"`
	}

	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, err
	}

	config := new(router.Config)
	for prefix, location := range rawConfig.Locations {
		config.Locations = append(config.Locations, router.Location{
			Prefix: prefix,
			URL:    location.URL,
		})
	}

	return config, nil
}
