package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	S3 struct {
		Region          string `yaml:"region"`
		FileBucket      string `yaml:"file_bucket"`
		ImageBucket     string `yaml:"image_bucket"`
		ImageURLPrefix  string `yaml:"image_url_prefix"`
		EndpointURL     string `yaml:"endpoint_url"`
		SecretAccessKey string
		AccessKeyId     string
	} `yaml:"s3"`

	PGConnString string

	StorageURL string `yaml:"storage_url"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.S3.AccessKeyId = os.Getenv("ACCESS_KEY_ID")
	cfg.S3.SecretAccessKey = os.Getenv("SECRET_ACCESS_KEY")

	cfg.PGConnString = os.Getenv("PG_CONN_STRING")

	return cfg, nil
}
