package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type BackuperConfig struct {
	HomeAssistantPath   string
	ServiceAccountPath  string
	GcloudProject       string
	BucketName          string
	LocationIdentifier  string
	FirestoreCollection string
	WebhookEnabled      bool
	WebhookUrl          *string
}

const configFileName = "config.json"

func LoadConfig() (*BackuperConfig, error) {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s not found", configFileName)
	}

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	config := &BackuperConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}
	fmt.Println("BackuperConfig loaded")
	return config, nil
}

func validateConfig(config *BackuperConfig) error {
	var missingFields []string
	if config.HomeAssistantPath == "" {
		missingFields = append(missingFields, "Home Assistant path")
	}
	if config.ServiceAccountPath == "" {
		missingFields = append(missingFields, "Service account path")
	}
	if config.GcloudProject == "" {
		missingFields = append(missingFields, "Gcloud project")
	}
	if config.BucketName == "" {
		missingFields = append(missingFields, "Bucket name")
	}
	if config.LocationIdentifier == "" {
		missingFields = append(missingFields, "Location identifier")
	}
	if config.FirestoreCollection == "" {
		missingFields = append(missingFields, "Firebase collection")
	}
	if config.WebhookEnabled && config.WebhookUrl == nil {
		missingFields = append(missingFields, "Webhook URL")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing fields: %s", missingFields)
	}

	return nil
}
