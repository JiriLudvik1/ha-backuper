package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HomeAssistantPath  string
	ServiceAccountPath string
	BucketName         string
	LocationIdentifier string
	WebhookEnabled     bool
	WebhookUrl         *string
}

const configFileName = "config.json"

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s not found", configFileName)
	}

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	configFile.Close()

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}
	fmt.Println("Config loaded")
	return config, nil
}

func validateConfig(config *Config) error {
	var missingFields []string
	if config.HomeAssistantPath == "" {
		missingFields = append(missingFields, "Home Assistant path")
	}
	if config.ServiceAccountPath == "" {
		missingFields = append(missingFields, "Service account path")
	}
	if config.BucketName == "" {
		missingFields = append(missingFields, "Bucket name")
	}
	if config.LocationIdentifier == "" {
		missingFields = append(missingFields, "Location identifier")
	}
	if config.WebhookEnabled && config.WebhookUrl == nil {
		missingFields = append(missingFields, "Webhook URL")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing fields: %s", missingFields)
	}

	return nil
}
